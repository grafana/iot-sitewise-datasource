package api

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
)

const (
	MaxResourcesPerRequest = 100
)

func GetResourcesTags(ctx context.Context, rgtClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI, arns []*string) ([]*resourcegroupstaggingapi.ResourceTagMapping, error) {
	var resources []*resourcegroupstaggingapi.ResourceTagMapping

	// Get resources by MaxResourcesPerRequest per request
	for i := 0; i < len(arns); i += MaxResourcesPerRequest {
		res, err := getResources(ctx, rgtClient, arns[i:min(i+MaxResourcesPerRequest, len(arns))])
		if err != nil {
			log.Printf("error getting resources: %v", err)
			return nil, err
		}
		resources = append(resources, res.ResourceTagMappingList...)
	}

	return resources, nil
}

func getResources(ctx context.Context, rgtClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI, arns []*string) (*resourcegroupstaggingapi.GetResourcesOutput, error) {
	var resourceTagMappingList []*resourcegroupstaggingapi.ResourceTagMapping
	var paginationToken *string

	for {
		resources, err := rgtClient.GetResourcesWithContext(ctx, &resourcegroupstaggingapi.GetResourcesInput{
			ResourceARNList: arns,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting resources: %v", err)
		}

		paginationToken = resources.PaginationToken

		if paginationToken == nil {
			break
		}
	}

	return &resourcegroupstaggingapi.GetResourcesOutput{
		ResourceTagMappingList: resourceTagMappingList,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
