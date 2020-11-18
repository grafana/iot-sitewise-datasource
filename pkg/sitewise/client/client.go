//go:generate mockery --name SitewiseClient

package client

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/aws/aws-sdk-go/service/iotsitewise/iotsitewiseiface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
)

type SitewiseClient interface {
	iotsitewiseiface.IoTSiteWiseAPI
	GetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyValueHistoryOutput, error)
	GetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyAggregatesOutput, error)
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

func GetClient(region string, settings awsds.AWSDatasourceSettings, provider awsds.AmazonSessionProvider) (client SitewiseClient, err error) {

	sess, err := provider(region, settings)
	if err != nil {
		return nil, err
	}

	swcfg := &aws.Config{}

	if settings.Endpoint != "" {
		swcfg.Endpoint = aws.String(settings.Endpoint)
	}

	return &sitewiseClient{iotsitewise.New(sess, swcfg)}, nil
}
