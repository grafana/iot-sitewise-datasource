//go:generate mockery --name SitewiseAPIClient

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

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type DescribeAssetPropertyAPIClient interface {
	DescribeAssetProperty(ctx context.Context, params *iotsitewise.DescribeAssetPropertyInput, optFns ...func(*iotsitewise.Options)) (*iotsitewise.DescribeAssetPropertyOutput, error)
}
type DescribeTimeseriesAPIClient interface {
	DescribeTimeSeries(ctx context.Context, params *iotsitewise.DescribeTimeSeriesInput, optFns ...func(*iotsitewise.Options)) (*iotsitewise.DescribeTimeSeriesOutput, error)
}
type GetAssetPropertyValueAPIClient interface {
	GetAssetPropertyValue(ctx context.Context, params *iotsitewise.GetAssetPropertyValueInput, optFns ...func(*iotsitewise.Options)) (*iotsitewise.GetAssetPropertyValueOutput, error)
}

type SitewiseAPIClient interface {
	iotsitewise.BatchGetAssetPropertyAggregatesAPIClient
	iotsitewise.BatchGetAssetPropertyValueAPIClient
	iotsitewise.BatchGetAssetPropertyValueHistoryAPIClient
	iotsitewise.DescribeAssetAPIClient
	iotsitewise.DescribeAssetModelAPIClient
	iotsitewise.ExecuteQueryAPIClient
	iotsitewise.GetAssetPropertyAggregatesAPIClient
	iotsitewise.GetAssetPropertyValueHistoryAPIClient
	iotsitewise.GetInterpolatedAssetPropertyValuesAPIClient
	iotsitewise.ListAssetsAPIClient
	iotsitewise.ListAssetModelsAPIClient
	iotsitewise.ListAssetPropertiesAPIClient
	iotsitewise.ListAssociatedAssetsAPIClient
	iotsitewise.ListTimeSeriesAPIClient

	DescribeAssetPropertyAPIClient
	DescribeTimeseriesAPIClient
	GetAssetPropertyValueAPIClient

	BatchGetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, error)
	GetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyValueHistoryOutput, error)
	GetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyAggregatesOutput, error)
	BatchGetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyAggregatesOutput, error)
	GetInterpolatedAssetPropertyValuesPageAggregation(ctx context.Context, req *iotsitewise.GetInterpolatedAssetPropertyValuesInput, maxPages int, maxResults int) (*iotsitewise.GetInterpolatedAssetPropertyValuesOutput, error)
}

type SitewiseClient struct {
	*iotsitewise.Client
}

// NewSitewiseClientForRegion is mainly for testing in this case
// TODO: move this into one of the test files
func NewSitewiseClientForRegion(region string) SitewiseAPIClient {
	cfg, _ := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(region))
	return &SitewiseClient{Client: iotsitewise.NewFromConfig(cfg)}
}

func (c *SitewiseClient) GetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyValueHistoryOutput, error) {
	var (
		numPages  = 0
		values    []iotsitewisetypes.AssetPropertyValue
		nextToken *string
	)

	pager := iotsitewise.NewGetAssetPropertyValueHistoryPaginator(c.Client, req)
	for pager.HasMorePages() && numPages < maxPages && len(values) <= maxResults {
		numPages += 1
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		nextToken = page.NextToken
		values = append(values, page.AssetPropertyValueHistory...)
	}

	return &iotsitewise.GetAssetPropertyValueHistoryOutput{
		AssetPropertyValueHistory: values,
		NextToken:                 nextToken,
	}, nil
}

func (c *SitewiseClient) BatchGetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyValueHistoryOutput, error) {
	var (
		count     = 0
		numPages  = 0
		success   []iotsitewisetypes.BatchGetAssetPropertyValueHistorySuccessEntry
		skipped   []iotsitewisetypes.BatchGetAssetPropertyValueHistorySkippedEntry
		errs      []iotsitewisetypes.BatchGetAssetPropertyValueHistoryErrorEntry
		nextToken *string
	)

	pager := iotsitewise.NewBatchGetAssetPropertyValueHistoryPaginator(c.Client, req)
	for pager.HasMorePages() && numPages < maxPages && count <= maxResults {
		numPages += 1
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if len(page.SuccessEntries) > 0 {
			count += len(page.SuccessEntries[0].AssetPropertyValueHistory)
		}
		if len(success) > 0 {
			for _, successEntry := range page.SuccessEntries {
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
			success = append(success, page.SuccessEntries...)
		}
		skipped = append(skipped, page.SkippedEntries...)
		errs = append(errs, page.ErrorEntries...)
		nextToken = page.NextToken
	}

	return &iotsitewise.BatchGetAssetPropertyValueHistoryOutput{
		SuccessEntries: success,
		SkippedEntries: skipped,
		ErrorEntries:   errs,
		NextToken:      nextToken,
	}, nil
}

func (c *SitewiseClient) GetInterpolatedAssetPropertyValuesPageAggregation(ctx context.Context, req *iotsitewise.GetInterpolatedAssetPropertyValuesInput, maxPages int, maxResults int) (*iotsitewise.GetInterpolatedAssetPropertyValuesOutput, error) {
	var (
		numPages  = 0
		values    []iotsitewisetypes.InterpolatedAssetPropertyValue
		nextToken *string
	)

	pager := iotsitewise.NewGetInterpolatedAssetPropertyValuesPaginator(c.Client, req)
	for pager.HasMorePages() && numPages < maxPages {
		numPages += 1
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		values = append(values, page.InterpolatedAssetPropertyValues...)
		nextToken = page.NextToken
	}

	return &iotsitewise.GetInterpolatedAssetPropertyValuesOutput{
		InterpolatedAssetPropertyValues: values,
		NextToken:                       nextToken,
	}, nil
}

func (c *SitewiseClient) GetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyAggregatesOutput, error) {
	var (
		numPages  = 0
		values    []iotsitewisetypes.AggregatedValue
		nextToken *string
	)

	pager := iotsitewise.NewGetAssetPropertyAggregatesPaginator(c.Client, req)
	for pager.HasMorePages() && numPages < maxPages && len(values) <= maxResults {
		numPages += 1
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		values = append(values, page.AggregatedValues...)
		nextToken = page.NextToken
	}

	return &iotsitewise.GetAssetPropertyAggregatesOutput{
		AggregatedValues: values,
		NextToken:        nextToken,
	}, nil
}

func (c *SitewiseClient) BatchGetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.BatchGetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.BatchGetAssetPropertyAggregatesOutput, error) {

	var (
		count     = 0
		numPages  = 0
		success   []iotsitewisetypes.BatchGetAssetPropertyAggregatesSuccessEntry
		skipped   []iotsitewisetypes.BatchGetAssetPropertyAggregatesSkippedEntry
		errs      []iotsitewisetypes.BatchGetAssetPropertyAggregatesErrorEntry
		nextToken *string
	)

	pager := iotsitewise.NewBatchGetAssetPropertyAggregatesPaginator(c.Client, req)
	for pager.HasMorePages() && numPages < maxPages && count <= maxResults {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		if len(page.SuccessEntries) > 0 {
			count += len(page.SuccessEntries[0].AggregatedValues)
		}
		if len(success) > 0 {
			for _, successEntry := range page.SuccessEntries {
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
			success = append(success, page.SuccessEntries...)
		}
		skipped = append(skipped, page.SkippedEntries...)
		errs = append(errs, page.ErrorEntries...)
		nextToken = page.NextToken
	}

	return &iotsitewise.BatchGetAssetPropertyAggregatesOutput{
		SuccessEntries: success,
		SkippedEntries: skipped,
		ErrorEntries:   errs,
		NextToken:      nextToken,
	}, nil
}

func GetAWSConfig(ctx context.Context, settings models.AWSSiteWiseDataSourceSetting) (cfg aws.Config, err error) {
	options := make([]func(*awsconfig.LoadOptions) error, 0)
	if settings.Endpoint != "" {
		options = append(options, awsconfig.WithBaseEndpoint(settings.Endpoint))
	}

	if settings.Region == models.EDGE_REGION {
		pool, _ := x509.SystemCertPool()
		if pool == nil {
			pool = x509.NewCertPool()
		}

		if settings.Cert == "" {
			return cfg, errors.New("certificate cannot be null")
		}

		block, _ := pem.Decode([]byte(settings.Cert))
		if block == nil || block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			return cfg, fmt.Errorf("decode certificate failed: %s", settings.Cert)
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return cfg, err
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

		options = append(options, awsconfig.WithHTTPClient(&http.Client{Transport: tr}))
		// TODO: figure out how to replace this with smithy's version
		// https://pkg.go.dev/github.com/aws/smithy-go/transport/http#DisableEndpointHostPrefix
		//swcfg = swcfg.WithDisableEndpointHostPrefix(true)

	}
	return awsconfig.LoadDefaultConfig(ctx, options...)
}
