package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func BatchGetAssetPropertyValuesForTimeRange(ctx context.Context, client client.SitewiseClient,
	query models.AssetPropertyValueQuery) (models.AssetPropertyValueQuery, *framer.AssetPropertyValuesForTimeRangeBatch, error) {

	if query.Resolution == "AUTO" {
		resolution := propvals.Resolution(query.BaseQuery)

		// todo: remove propvals.ResolutionSecond condition once 1s aggregation is supported
		if propvals.ResolutionRaw == resolution || propvals.ResolutionSecond == resolution {
			modifiedQuery, history, err := BatchGetAssetPropertyValues(ctx, client, query)
			if err != nil {
				return modifiedQuery, nil, err
			}
			return modifiedQuery, &framer.AssetPropertyValuesForTimeRangeBatch{History: history}, nil
		}

	}

	modifiedQuery, aggregates, err := BatchGetAssetPropertyAggregates(ctx, client, query)
	if err != nil {
		return modifiedQuery, nil, err
	}
	return modifiedQuery, &framer.AssetPropertyValuesForTimeRangeBatch{Aggregates: aggregates}, nil
}
