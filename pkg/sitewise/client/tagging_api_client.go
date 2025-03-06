//go:generate mockery --name SitewiseClient

package client

import (
	"context"
	"math"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"golang.org/x/sync/errgroup"
)

const (
	MaxResourcesPerRequest = 100
)

type TaggingApiClient interface {
	resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	GetResourcesPage(ctx context.Context, arns []*string) (*resourcegroupstaggingapi.GetResourcesOutput, error)
}

type taggingApiClient struct {
	*resourcegroupstaggingapi.ResourceGroupsTaggingAPI
}

func (rgtClient *taggingApiClient) GetResourcesPage(ctx context.Context, arns []*string) (*resourcegroupstaggingapi.GetResourcesOutput, error) {
	requestNumber := int(math.Ceil(float64(len(arns)) / float64(MaxResourcesPerRequest)))
	resultChan := make(chan *resourcegroupstaggingapi.GetResourcesOutput, requestNumber)
	eg, ectx := errgroup.WithContext(ctx)

	// Get resources by MaxResourcesPerRequest per request
	for i := 0; i < len(arns); i += MaxResourcesPerRequest {
		eg.Go(func() error {
			res, err := rgtClient.getResources(ectx, arns[i:min(i+MaxResourcesPerRequest, len(arns))])
			if err != nil {
				return err
			}
			resultChan <- res

			return nil
		})
	}

	err := eg.Wait()
	close(resultChan)
	if err != nil {
		return nil, err
	}

	var resources []*resourcegroupstaggingapi.ResourceTagMapping

	// append the ResourceTagMappingList from each response
	for result := range resultChan {
		resources = append(resources, result.ResourceTagMappingList...)
	}

	output := resourcegroupstaggingapi.GetResourcesOutput{
		ResourceTagMappingList: resources,
	}

	return &output, nil
}

func (rgtClient *taggingApiClient) getResources(ctx context.Context, arns []*string) (*resourcegroupstaggingapi.GetResourcesOutput, error) {
	resources, err := rgtClient.GetResourcesWithContext(ctx, &resourcegroupstaggingapi.GetResourcesInput{
		ResourceARNList: arns,
	})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func GetTaggingApiClient(region string, settings models.AWSSiteWiseDataSourceSetting, provider AmazonSessionProvider, authSettings *awsds.AuthSettings) (client TaggingApiClient, err error) {
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

	c := resourcegroupstaggingapi.New(sess, swcfg)

	// TODO: user agent for metrics tracking purpose
	// c.Handlers.Send.PushFront(func(r *request.Request) {
	// 	r.HTTPRequest.Header.Set("User-Agent", awsds.GetUserAgentString("grafana-iot-sitewise-datasource"))
	// })
	return &taggingApiClient{c}, nil
}
