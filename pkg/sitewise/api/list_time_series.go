package api

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
	
)

func ListTimeSeries(ctx context.Context, client client.SitewiseClient, query models.ListTimeSeriesQuery) (*framer.TimeSeries, error) {

	var (
		timeSeriesType *string
		assetId *string = util.GetAssetId(query.BaseQuery)
		aliasPrefix *string
	)

	if query.TimeSeriesType != "" {
		// if user wants to see all timeseries data do not filter on type
		if query.TimeSeriesType == "ALL" { 
			timeSeriesType = nil
		} else {
			timeSeriesType = aws.String(query.TimeSeriesType)
		}	
	}

	if query.AliasPrefix != "" {
		aliasPrefix = aws.String(query.AliasPrefix)
	}

	if query.TimeSeriesType == "DISASSOCIATED" {
		// cannot filter on assetId for disassociated data
		assetId = nil
	}

	if query.TimeSeriesType == "ASSOCIATED" {
		// cannot filter by alias prefix on associated data
		aliasPrefix = nil
	}

	resp, err := client.ListTimeSeriesWithContext(ctx, &iotsitewise.ListTimeSeriesInput{
		AssetId: 		assetId,
		TimeSeriesType: timeSeriesType,
		AliasPrefix:  	aliasPrefix,
		MaxResults:   	aws.Int64(250),
		NextToken:    	getNextToken(query.BaseQuery),
	})


	if err != nil {
		return nil, err
	}

	return &framer.TimeSeries{
		TimeSeriesSummaries: resp.TimeSeriesSummaries,
		NextToken:      resp.NextToken,
	}, nil
}
