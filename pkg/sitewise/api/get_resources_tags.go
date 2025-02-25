// const { GetResourcesCommand } = require("@aws-sdk/client-resource-groups-tagging-api");

// exports.getResources = async ({rgtClient, arns}) => {
//   const response = {
//     ResourceTagMappingList: [],
//   };

//   // Get resources by 100 per request
//   for (let i = 0; i < arns.length; i += 100) {
//     const resources = await getResourcesMaxNum({rgtClient, arns: arns.slice(i, i + 100)});

//     response.ResourceTagMappingList = response.ResourceTagMappingList.concat(...resources.ResourceTagMappingList);
//   }

//   return response;
// };

// const getResourcesMaxNum = async ({rgtClient, arns}) => {
//   const response = {
//     ResourceTagMappingList: [],
//   };

//   // paginate through the next token
//   let nextToken;
//   do {
//     const command = new GetResourcesCommand({
//       ResourceARNList: arns,
//       NextToken: nextToken,
//     });

//     const resources = await rgtClient.send(command);

//     response.ResourceTagMappingList = response.ResourceTagMappingList.concat(...resources.ResourceTagMappingList);

//     nextToken = resources.NextToken;
//   } while (nextToken);

//   return response;
// };

// Translate the js above into go lang
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
	resources, err := rgtClient.GetResourcesWithContext(ctx, &resourcegroupstaggingapi.GetResourcesInput{
		ResourceARNList: arns,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %v", err)
	}

	return resources, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
