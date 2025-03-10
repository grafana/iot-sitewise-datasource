package api

import (
	"context"
	"math"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

const (
	l4eAnomalyResultPropertyName = "AWS/L4E_ANOMALY_RESULT"
)

var (
	MaxSitewiseResults = aws.Int64(250)
)

func getNextToken(query models.BaseQuery) *string {
	// NextTokens should only be set for time series queries with batched entries
	if len(query.NextTokens) > 0 && len(query.AssetPropertyEntries) > 0 {
		var entryId string
		// We can look up the nextToken for any property in the given query
		// because it has been batched to account for API max entry limits
		// Every property in the batch has the same nextToken
		entry := query.AssetPropertyEntries[0]
		if entry.AssetId != "" && entry.PropertyId != "" {
			entryId = *util.GetEntryIdFromAssetProperty(entry.AssetId, entry.PropertyId)
		} else {
			entryId = *util.GetEntryIdFromPropertyAlias(entry.PropertyAlias)
		}
		// If there are any issues looking up the nextToken it should error and bubble up
		nextToken := query.NextTokens[entryId]
		return aws.String(nextToken)
	} else {
		if query.NextToken == "" {
			return nil
		}
		return aws.String(query.NextToken)
	}
}

func getAssetIdAndPropertyId(query models.AssetPropertyValueQuery, client client.SitewiseClient, ctx context.Context) (models.AssetPropertyValueQuery, error) {
	result := query
	result.AssetPropertyEntries = []models.AssetPropertyEntry{}
	// There should only be a list of property aliases OR lists for assetIds and propertyIds
	// Look up the assetId and propertyId for a property alias
	if len(query.PropertyAliases) > 0 {
		for _, propertyAlias := range query.PropertyAliases {
			resp, err := client.DescribeTimeSeriesWithContext(ctx, &iotsitewise.DescribeTimeSeriesInput{
				Alias: aws.String(propertyAlias),
			})
			if err != nil {
				return models.AssetPropertyValueQuery{}, err
			}
			if resp.AssetId != nil && resp.PropertyId != nil {
				result.AssetPropertyEntries = append(result.AssetPropertyEntries, models.AssetPropertyEntry{
					AssetId:       *resp.AssetId,
					PropertyId:    *resp.PropertyId,
					PropertyAlias: propertyAlias,
				})
			} else {
				// For disassociated streams with just a propertyAlias
				result.AssetPropertyEntries = append(result.AssetPropertyEntries, models.AssetPropertyEntry{
					PropertyAlias: propertyAlias,
				})
			}
		}
	} else {
		for _, assetId := range query.AssetIds {
			for _, propertyId := range query.PropertyIds {
				result.AssetPropertyEntries = append(result.AssetPropertyEntries, models.AssetPropertyEntry{
					AssetId:    assetId,
					PropertyId: propertyId,
				})
			}
		}
	}
	return result, nil
}

func getFirstAssetId(query models.BaseQuery) *string {
	if len(query.AssetIds) == 0 {
		return nil
	}
	return aws.String(query.AssetIds[0])
}

func getFirstPropertyId(query models.BaseQuery) *string {
	if len(query.PropertyIds) == 0 {
		return nil
	}
	return aws.String(query.PropertyIds[0])
}

func getFirstPropertyAlias(query models.BaseQuery) *string {
	if len(query.PropertyAliases) == 0 {
		return nil
	}
	return aws.String(query.PropertyAliases[0])
}

func filterAnomalyAssetIds(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) ([]string, error) {
	anomalyAssetIds := []string{}

	switch {
	case len(query.PropertyAliases) > 0:
		return nil, nil

	default:
		for _, assetId := range query.AssetIds {
			for _, propertyId := range query.PropertyIds {
				var id *string
				if assetId != "" {
					id = aws.String(assetId)
				}

				req := &iotsitewise.DescribeAssetPropertyInput{
					AssetId:    id,
					PropertyId: aws.String(propertyId),
				}

				resp, err := client.DescribeAssetPropertyWithContext(ctx, req)
				if err != nil {
					return nil, err
				}

				if resp.CompositeModel != nil && *resp.CompositeModel.AssetProperty.Name == l4eAnomalyResultPropertyName {
					anomalyAssetIds = append(anomalyAssetIds, assetId)
				}
			}
		}
	}

	return anomalyAssetIds, nil
}

// Batch queries with a consistent next token for each batch
func batchQueriesWithNextToken(query models.AssetPropertyValueQuery) []models.AssetPropertyValueQuery {
	batchQueries := []models.AssetPropertyValueQuery{}

	nextTokenGroups := map[string][]models.AssetPropertyEntry{}
	for _, entry := range query.AssetPropertyEntries {
		entryId := util.GetEntryIdFromAssetPropertyEntry(entry)
		// The NextTokens field is set in src/getNextQueries.ts and maps a batched query entryId to a nextToken
		nextToken := query.NextTokens[*entryId]
		// Do not continue to process batches that are complete
		if nextToken != "" {
			// Build map of a nextToken to a list of entries to continue querying with that nextToken
			if _, exists := nextTokenGroups[nextToken]; exists {
				nextTokenGroups[nextToken] = append(nextTokenGroups[nextToken], entry)
			} else {
				entries := []models.AssetPropertyEntry{entry}
				nextTokenGroups[nextToken] = entries
			}
		}
	}
	for _, v := range nextTokenGroups {
		q := query
		q.AssetPropertyEntries = v
		batchQueries = append(batchQueries, q)
	}

	return batchQueries
}

func batchQueriesInitial(query models.AssetPropertyValueQuery, maxBatchSize int) []models.AssetPropertyValueQuery {
	batchQueries := []models.AssetPropertyValueQuery{}

	numEntries := len(query.AssetPropertyEntries)
	numBatches := int(math.Ceil(float64(numEntries) / float64(maxBatchSize)))
	var index = 0
	for i := 0; i < numBatches; i++ {
		q := query
		// Iterate in batches each with a size at most maxBatchSize
		batchEndIndex := index + maxBatchSize
		if batchEndIndex <= numEntries {
			q.AssetPropertyEntries = query.AssetPropertyEntries[index:batchEndIndex]
			index += maxBatchSize
		} else {
			q.AssetPropertyEntries = query.AssetPropertyEntries[index:]
		}
		batchQueries = append(batchQueries, q)
	}

	return batchQueries
}

func batchQueries(query models.AssetPropertyValueQuery, maxBatchSize int) []models.AssetPropertyValueQuery {
	// If the API entry limit is not exceeded no need to batch further
	if len(query.AssetPropertyEntries) <= maxBatchSize {
		return []models.AssetPropertyValueQuery{query}
	}

	if len(query.NextTokens) > 0 {
		return batchQueriesWithNextToken(query)
	} else {
		return batchQueriesInitial(query, maxBatchSize)
	}
}
