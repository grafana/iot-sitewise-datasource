package api

import (
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
)

func FilterAssetSummariesByArns(assetSummaries []*iotsitewise.AssetSummary, allowArns map[string]bool) []*iotsitewise.AssetSummary {
	filteredAssetSummaries := []*iotsitewise.AssetSummary{}
	for _, asset := range assetSummaries {
		if allowArns[*asset.Arn] {
			filteredAssetSummaries = append(filteredAssetSummaries, asset)
		}
	}
	return filteredAssetSummaries
}

func FilterAssociatedAssetSummariesByArns(assetSummaries []*iotsitewise.AssociatedAssetsSummary, allowArns map[string]bool) []*iotsitewise.AssociatedAssetsSummary {
	filteredAssetSummaries := []*iotsitewise.AssociatedAssetsSummary{}
	for _, asset := range assetSummaries {
		if allowArns[*asset.Arn] {
			filteredAssetSummaries = append(filteredAssetSummaries, asset)
		}
	}
	return filteredAssetSummaries
}

func FilterResourcesByTags(resources *resourcegroupstaggingapi.GetResourcesOutput, includedTagPatterns []map[string][]string) map[string]bool {
	// filter the given resources by matching allowTags
	allowArns := map[string]bool{}

	// none allowed
	if len(includedTagPatterns) == 0 {
		return allowArns
	}

	// for each resource
	for _, resource := range resources.ResourceTagMappingList {
		if resource.Tags != nil {
			// match every tag policy
			for _, allowTags := range includedTagPatterns {
				matchedTags := true
				// match every tag
				for allowTagKey, allowTagValues := range allowTags {
					matchedTagFound := false
					// look for the matching resource tag
					for _, tag := range resource.Tags {
						if tag.Key != nil && tag.Value != nil {
							if *tag.Key == allowTagKey {
								// if the tag value is in the allow list, return the resource
								for _, allowTagValue := range allowTagValues {
									if *tag.Value == allowTagValue {
										matchedTagFound = true
									}
								}
							}
						}
					}
					if !matchedTagFound {
						// mismatched found; terminates execution;
						matchedTags = false
					}
				}

				if matchedTags {
					allowArns[*resource.ResourceARN] = true
				}
			}
		}
	}

	return allowArns
}
