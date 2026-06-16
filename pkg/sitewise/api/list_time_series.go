package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func ListTimeSeries(ctx context.Context, sw client.SitewiseAPIClient, query models.ListTimeSeriesQuery) (*framer.TimeSeries, error) {

	var (
		aliasPrefix *string
	)
	timeSeriesType := query.TimeSeriesType
	assetId := util.GetAssetId(query.BaseQuery)

	if timeSeriesType == "ALL" {
		timeSeriesType = ""
	}

	if query.AliasPrefix != "" {
		aliasPrefix = aws.String(query.AliasPrefix)
	}

	if query.TimeSeriesType == iotsitewisetypes.ListTimeSeriesTypeDisassociated {
		// cannot filter on assetId for disassociated data
		assetId = nil
	}

	if query.TimeSeriesType == iotsitewisetypes.ListTimeSeriesTypeAssociated {
		// cannot filter by alias prefix on associated data
		aliasPrefix = nil
	}

	resp, err := sw.ListTimeSeries(ctx, &iotsitewise.ListTimeSeriesInput{
		AssetId:        assetId,
		TimeSeriesType: timeSeriesType,
		AliasPrefix:    aliasPrefix,
		MaxResults:     aws.Int32(250),
		NextToken:      getNextToken(query.BaseQuery),
	})

	if err != nil {
		return nil, err
	}

	return &framer.TimeSeries{
		TimeSeriesSummaries: resp.TimeSeriesSummaries,
		NextToken:           resp.NextToken,
	}, nil
}
