package api

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/api/propvals"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
)

func GetAssetPropertyValuesForTimeRange(ctx context.Context, client client.SitewiseClient, query models.AssetPropertyValueQuery) (*framer.AssetPropertyValuesForTimeRange, error) {

	if query.Resolution == "AUTO" {
		resolution := propvals.Resolution(query.BaseQuery)

		// todo: remove propvals.ResolutionSecond condition once 1s aggregation is supported
		if propvals.ResolutionRaw == resolution || propvals.ResolutionSecond == resolution {
			history, err := BatchGetAssetPropertyValues(ctx, client, query)
			if err != nil {
				return nil, err
			}
			return &framer.AssetPropertyValuesForTimeRange{History: history}, nil
		}

	}

	aggregates, err := GetAssetPropertyAggregates(ctx, client, query)
	if err != nil {
		return nil, err
	}
	return &framer.AssetPropertyValuesForTimeRange{Aggregates: aggregates}, nil
}
