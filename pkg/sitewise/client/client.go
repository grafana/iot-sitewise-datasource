//go:generate mockery --name SitewiseClient

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/aws/aws-sdk-go/service/iotsitewise/iotsitewiseiface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type SitewiseClient interface {
	iotsitewiseiface.IoTSiteWiseAPI
	BatchGetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, error)
	BatchGetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyAggregatesOutput, error)
	GetInterpolatedAssetPropertyValuesPageAggregation(ctx context.Context, req *iotsitewise.GetInterpolatedAssetPropertyValuesInput, maxPages int, maxResults int) (*iotsitewise.GetInterpolatedAssetPropertyValuesOutput, error)
}

type sitewiseClient struct {
	iotsitewiseiface.IoTSiteWiseAPI
}

// NewSitewiseClient is mainly for testing in this case
func NewSitewiseClientForRegion(region string) SitewiseClient {
	sesh := session.Must(session.NewSession())
	sw := iotsitewise.New(sesh, aws.NewConfig().WithRegion(region))
	return &sitewiseClient{
		sw,
	}
}

type IndexedData struct {
	Index int
	Data  *iotsitewise.BatchGetAssetPropertyValueHistoryOutput
}

func fetchRawData(ctx context.Context, startTime int64, endTime int64, c iotsitewiseiface.IoTSiteWiseAPI, req *iotsitewise.BatchGetAssetPropertyValueHistoryInput) (*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, error) {
	var (
		count    = 0
		numPages = 0
		success  []*iotsitewise.BatchGetAssetPropertyValueHistorySuccessEntry
		skipped  []*iotsitewise.BatchGetAssetPropertyValueHistorySkippedEntry
		errors   []*iotsitewise.BatchGetAssetPropertyValueHistoryErrorEntry
	)

	entries := make([]*iotsitewise.BatchGetAssetPropertyValueHistoryEntry, 0)
	for _, entry := range req.Entries {
		entries = append(entries, &iotsitewise.BatchGetAssetPropertyValueHistoryEntry{
			StartDate:    aws.Time(time.Unix(startTime, 0)),
			EndDate:      aws.Time(time.Unix(endTime, 0)),
			EntryId:      entry.EntryId,
			AssetId:      entry.AssetId,
			PropertyId:   entry.PropertyId,
			TimeOrdering: entry.TimeOrdering,
			Qualities:    entry.Qualities,
		})
	}

	slicedReq := &iotsitewise.BatchGetAssetPropertyValueHistoryInput{
		Entries:    entries,
		MaxResults: req.MaxResults,
		NextToken:  req.NextToken,
	}

	err := c.BatchGetAssetPropertyValueHistoryPagesWithContext(ctx, slicedReq,
		func(output *iotsitewise.BatchGetAssetPropertyValueHistoryOutput, b bool) bool {
			numPages++
			if len(output.SuccessEntries) > 0 {
				count += len(output.SuccessEntries[0].AssetPropertyValueHistory)
			}
			if len(success) > 0 {
				for _, successEntry := range output.SuccessEntries {
					found := false
					for i, entry := range success {
						if *entry.EntryId == *successEntry.EntryId {
							success[i].AssetPropertyValueHistory = append(success[i].AssetPropertyValueHistory, successEntry.AssetPropertyValueHistory...)
							found = true
							break
						}
					}
					if !found {
						success = append(success, successEntry)
					}
				}
			} else {
				success = append(success, output.SuccessEntries...)
			}
			skipped = append(skipped, output.SkippedEntries...)
			errors = append(errors, output.ErrorEntries...)
			return true
		})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: success,
		SkippedEntries: skipped,
		ErrorEntries:   errors,
		NextToken:      nil,
	}, nil
}

func (c *sitewiseClient) BatchGetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, error) {
	var (
		success []*iotsitewise.BatchGetAssetPropertyValueHistorySuccessEntry
		skipped []*iotsitewise.BatchGetAssetPropertyValueHistorySkippedEntry
		errors  []*iotsitewise.BatchGetAssetPropertyValueHistoryErrorEntry
	)

	// Divide the time range into 10 chunks and fetch data in parallel
	startTime := req.Entries[0].StartDate.Unix()
	endTime := req.Entries[0].EndDate.Unix()
	intervalCount := 50
	totalDuration := endTime - startTime
	intervalDuration := totalDuration / int64(intervalCount)

	// If the time range is less than 10 seconds, we will fetch all data in one request
	if intervalDuration < 1 {
		intervalDuration = totalDuration
		intervalCount = 1
	}

	var wg sync.WaitGroup
	results := make([]*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, intervalCount)
	resultChannel := make(chan IndexedData, intervalCount)

	for i := 0; i < intervalCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			intervalStart := startTime + int64(i)*intervalDuration
			intervalEnd := intervalStart + intervalDuration
			if intervalEnd > endTime {
				intervalEnd = endTime
			}

			data, err := fetchRawData(ctx, intervalStart, intervalEnd, c, req)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			resultChannel <- IndexedData{Index: i, Data: data}
		}(i)
	}

	// wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	i := 0
	for res := range resultChannel {
		results[res.Index] = res.Data
		i++
	}

	// aggregate result back
	for _, result := range results {
		if len(success) > 0 {
			for _, successEntry := range result.SuccessEntries {
				found := false
				for i, entry := range success {
					if *entry.EntryId == *successEntry.EntryId {
						success[i].AssetPropertyValueHistory = append(success[i].AssetPropertyValueHistory, successEntry.AssetPropertyValueHistory...)
						found = true
						break
					}
				}
				if !found {
					success = append(success, successEntry)
				}
			}
		} else {
			success = append(success, result.SuccessEntries...)
		}
		skipped = append(skipped, result.SkippedEntries...)
		errors = append(errors, result.ErrorEntries...)
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: success,
		SkippedEntries: skipped,
		ErrorEntries:   errors,
		NextToken:      nil,
	}, nil
}

func (c *sitewiseClient) GetInterpolatedAssetPropertyValuesPageAggregation(ctx context.Context, req *iotsitewise.GetInterpolatedAssetPropertyValuesInput, maxPages int, maxResults int) (*iotsitewise.GetInterpolatedAssetPropertyValuesOutput, error) {
	var (
		numPages  = 0
		values    []*iotsitewise.InterpolatedAssetPropertyValue
		nextToken *string
	)

	err := c.GetInterpolatedAssetPropertyValuesPagesWithContext(ctx, req, func(output *iotsitewise.GetInterpolatedAssetPropertyValuesOutput, b bool) bool {
		numPages++
		values = append(values, output.InterpolatedAssetPropertyValues...)
		nextToken = output.NextToken
		return numPages < maxPages
	})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.GetInterpolatedAssetPropertyValuesOutput{
		InterpolatedAssetPropertyValues: values,
		NextToken:                       nextToken,
	}, nil
}

func (c *sitewiseClient) BatchGetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyAggregatesOutput, error) {

	var (
		count     = 0
		numPages  = 0
		success   []*iotsitewise.BatchGetAssetPropertyAggregatesSuccessEntry
		skipped   []*iotsitewise.BatchGetAssetPropertyAggregatesSkippedEntry
		errors    []*iotsitewise.BatchGetAssetPropertyAggregatesErrorEntry
		nextToken *string
	)

	err := c.BatchGetAssetPropertyAggregatesPagesWithContext(ctx, req, func(output *iotsitewise.BatchGetAssetPropertyAggregatesOutput, b bool) bool {
		if len(output.SuccessEntries) > 0 {
			count += len(output.SuccessEntries[0].AggregatedValues)
		}
		if len(success) > 0 {
			for _, successEntry := range output.SuccessEntries {
				found := false
				for i, entry := range success {
					if *entry.EntryId == *successEntry.EntryId {
						success[i].AggregatedValues = append(success[i].AggregatedValues, successEntry.AggregatedValues...)
						found = true
						break
					}
				}
				if !found {
					success = append(success, successEntry)
				}
			}
		} else {
			success = append(success, output.SuccessEntries...)
		}
		skipped = append(skipped, output.SkippedEntries...)
		errors = append(errors, output.ErrorEntries...)
		nextToken = output.NextToken
		return numPages < maxPages && count <= maxResults
	})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.BatchGetAssetPropertyAggregatesOutput{
		SuccessEntries: success,
		SkippedEntries: skipped,
		ErrorEntries:   errors,
		NextToken:      nextToken,
	}, nil
}

type AmazonSessionProvider func(c awsds.SessionConfig) (*session.Session, error)

func GetClient(region string, settings models.AWSSiteWiseDataSourceSetting, provider AmazonSessionProvider) (client SitewiseClient, err error) {
	awsSettings := settings.ToAWSDatasourceSettings()
	awsSettings.Region = region
	sess, err := provider(awsds.SessionConfig{Settings: awsSettings})
	if err != nil {
		return nil, err
	}

	swcfg := &aws.Config{}
	if settings.Endpoint != "" {
		swcfg.Endpoint = aws.String(settings.Endpoint)
	}

	if settings.Region == models.EDGE_REGION {
		pool, _ := x509.SystemCertPool()
		if pool == nil {
			pool = x509.NewCertPool()
		}

		if settings.Cert == "" {
			return nil, errors.New("certificate cannot be null")
		}

		block, _ := pem.Decode([]byte(settings.Cert))
		if block == nil || block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			return nil, fmt.Errorf("decode certificate failed: %s", settings.Cert)
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		pool.AddCert(cert)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //Not actually skipping, check the cert in VerifyPeerCertificate
				RootCAs:            pool,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					// If this is the first handshake on a connection, process and
					// (optionally) verify the server's certificates.
					certs := make([]*x509.Certificate, len(rawCerts))
					for i, asn1Data := range rawCerts {
						cert, err := x509.ParseCertificate(asn1Data)
						if err != nil {
							return errors.New("tls: failed to parse certificate from server: " + err.Error())
						}
						certs[i] = cert
					}

					opts := x509.VerifyOptions{
						Roots:         pool,
						CurrentTime:   time.Now(),
						DNSName:       "", // <- skip hostname verification
						Intermediates: x509.NewCertPool(),
					}

					for i, cert := range certs {
						if i == 0 {
							continue
						}
						opts.Intermediates.AddCert(cert)
					}
					_, err := certs[0].Verify(opts)
					return err
				},
			},
		}
		httpClient := &http.Client{Transport: tr}

		swcfg = swcfg.WithHTTPClient(httpClient)
		swcfg = swcfg.WithDisableEndpointHostPrefix(true)

	}

	c := iotsitewise.New(sess, swcfg)

	c.Handlers.Send.PushFront(func(r *request.Request) {
		r.HTTPRequest.Header.Set("User-Agent", awsds.GetUserAgentString("grafana-iot-sitewise-datasource"))
	})
	return &sitewiseClient{c}, nil
}
