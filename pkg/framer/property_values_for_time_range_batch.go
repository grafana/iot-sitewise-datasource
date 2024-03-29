package framer

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
	"github.com/pkg/errors"
)

type AssetPropertyValuesForTimeRangeBatch struct {
	History    *AssetPropertyValueHistoryBatch
	Aggregates *AssetPropertyAggregatesBatch
}

func (a *AssetPropertyValuesForTimeRangeBatch) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {

	if a.History != nil && a.Aggregates != nil {
		return nil, errors.New("unexpected state: AssetPropertyValuesForTimeRangeBatch should only have 'history' OR 'aggregate' response")
	}

	if a.History != nil {
		return a.History.Frames(ctx, resources)
	}

	if a.Aggregates != nil {
		return a.Aggregates.Frames(ctx, resources)
	}

	return nil, errors.New("no response found for AssetPropertyValuesForTimeRangeBatch")
}
