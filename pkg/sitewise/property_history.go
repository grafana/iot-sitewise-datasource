package sitewise

import (
	"context"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fdata"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

// GetAssetPropertyValueHistory requires either PropertyAlias OR (AssetID and PropertyID) to be set.
// The front end component should ensure that both cannot be sent at the same time by the user.
// If an invalid combo of assetId/propertyId/propertyAlias are sent to the API, an exception will be returned.
// The Framer consumer should bubble up that error to the user.
func historyQueryToInput(query models.AssetPropertyValueQuery) *iotsitewise.GetAssetPropertyValueHistoryInput {

	var (
		propertyId *string
		assetId    *string
		nextToken  *string
		qualities  []*string
	)

	if query.AssetId != "" {
		assetId = aws.String(query.AssetId)
	}

	if query.PropertyId != "" {
		propertyId = aws.String(query.PropertyId)
	}

	if query.NextToken != "" {
		nextToken = aws.String(query.NextToken)
	}

	if len(query.Qualities) > 0 {
		qualities = aws.StringSlice(query.Qualities)
	}

	from, to := util.TimeRangeToUnix(query.TimeRange)

	return &iotsitewise.GetAssetPropertyValueHistoryInput{
		AssetId:    assetId,
		StartDate:  from,
		EndDate:    to,
		MaxResults: aws.Int64(250), // should this even be configurable? 250 == max
		NextToken:  nextToken,
		PropertyId: propertyId,
		Qualities:  qualities,
	}
}

func GetAssetPropertyValues(ctx context.Context, client client.Client, query models.AssetPropertyValueQuery) (*fdata.AssetPropertyValueHistory, error) {

	awsReq := historyQueryToInput(query)

	// NOTE: there is a paginated API if we want to push pagination requests down to the server
	// See: https://docs.aws.amazon.com/sdk-for-go/api/service/iotsitewise/#IoTSiteWise.GetAssetPropertyValueHistoryPagesWithContext
	resp, err := client.GetAssetPropertyValueHistoryWithContext(ctx, awsReq)

	if err != nil {
		return nil, err
	}

	return &fdata.AssetPropertyValueHistory{
		AssetPropertyValueHistory: resp.AssetPropertyValueHistory,
		NextToken:                 resp.NextToken,
	}, nil
}
