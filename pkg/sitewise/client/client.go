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
	GetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyValueHistoryOutput, error)
	GetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyAggregatesOutput, error)
	BatchGetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, error)
	BatchGetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyAggregatesOutput, error)
	GetInterpolatedAssetPropertyValuesPageAggregation(ctx context.Context, req *iotsitewise.GetInterpolatedAssetPropertyValuesInput, maxPages int, maxResults int) (*iotsitewise.GetInterpolatedAssetPropertyValuesOutput, error)
}

type ListAssetPropertiesClient interface {
	ListAssetPropertiesWithContext(aws.Context, *iotsitewise.ListAssetPropertiesInput, ...request.Option) (*iotsitewise.ListAssetPropertiesOutput, error)
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

func (c *sitewiseClient) GetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyValueHistoryOutput, error) {
	var (
		numPages  = 0
		values    []*iotsitewise.AssetPropertyValue
		nextToken *string
	)

	err := c.GetAssetPropertyValueHistoryPagesWithContext(ctx, req, func(output *iotsitewise.GetAssetPropertyValueHistoryOutput, b bool) bool {
		numPages++
		values = append(values, output.AssetPropertyValueHistory...)
		nextToken = output.NextToken
		return numPages < maxPages && len(values) <= maxResults
	})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.GetAssetPropertyValueHistoryOutput{
		AssetPropertyValueHistory: values,
		NextToken:                 nextToken,
	}, nil
}

func (c *sitewiseClient) BatchGetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, error) {
	var (
		count     = 0
		numPages  = 0
		success   []*iotsitewise.BatchGetAssetPropertyValueHistorySuccessEntry
		skipped   []*iotsitewise.BatchGetAssetPropertyValueHistorySkippedEntry
		errors    []*iotsitewise.BatchGetAssetPropertyValueHistoryErrorEntry
		nextToken *string
	)

	err := c.BatchGetAssetPropertyValueHistoryPagesWithContext(ctx, req, func(output *iotsitewise.BatchGetAssetPropertyValueHistoryOutput, b bool) bool {
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
		nextToken = output.NextToken
		return numPages < maxPages && count <= maxResults
	})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: success,
		SkippedEntries: skipped,
		ErrorEntries:   errors,
		NextToken:      nextToken,
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

func (c *sitewiseClient) GetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyAggregatesOutput, error) {
	var (
		numPages  = 0
		values    []*iotsitewise.AggregatedValue
		nextToken *string
	)

	err := c.GetAssetPropertyAggregatesPagesWithContext(ctx, req, func(output *iotsitewise.GetAssetPropertyAggregatesOutput, b bool) bool {
		numPages++
		values = append(values, output.AggregatedValues...)
		nextToken = output.NextToken
		return numPages < maxPages && len(values) <= maxResults
	})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.GetAssetPropertyAggregatesOutput{
		AggregatedValues: values,
		NextToken:        nextToken,
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

type AmazonSessionProvider func(c awsds.GetSessionConfig, as awsds.AuthSettings) (*session.Session, error)

func GetClient(region string, settings models.AWSSiteWiseDataSourceSetting, provider AmazonSessionProvider, authSettings *awsds.AuthSettings) (client SitewiseClient, err error) {
	awsSettings := settings.ToAWSDatasourceSettings()
	awsSettings.Region = region
	sess, err := provider(awsds.GetSessionConfig{Settings: awsSettings}, *authSettings)
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
