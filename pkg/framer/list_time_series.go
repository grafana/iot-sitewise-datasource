package framer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
	
)

type TimeSeries iotsitewise.ListTimeSeriesOutput

type timeSeriesSummaryFields struct {
	alias         		*data.Field
	assetId           *data.Field
	dataType          *data.Field
	dataTypeSpec      *data.Field
	propertyId  	 *data.Field
	timeSeriesArn  *data.Field
	timeSeriesId  *data.Field
	timeSeriesCreationDate *data.Field
	timeSeriesLastUpdateDate  *data.Field
}

func (f *timeSeriesSummaryFields) fields() data.Fields {
	return data.Fields{
		f.alias,    
		f.assetId,
		f.dataType,
		f.dataTypeSpec,
		f.propertyId,
		f.timeSeriesArn,
		f.timeSeriesId,
		f.timeSeriesCreationDate,
		f.timeSeriesLastUpdateDate,
	}
}

func newTimeSeriesSummaryFields(length int) *timeSeriesSummaryFields {
	return &timeSeriesSummaryFields{
		alias:         fields.AliasField(length),
		assetId:           fields.AssetIdField(length),
		dataType:          fields.DataTypeField(length),
		dataTypeSpec:      fields.DataTypeSpecField(length),
		propertyId:  fields.PropertyIdField(length),
		timeSeriesArn:  fields.TimeSeriesArnField(length),
		timeSeriesId:  fields.TimeSeriesIdField(length),
		timeSeriesCreationDate: fields.TimeSeriesCreationDateField(length),
		timeSeriesLastUpdateDate:   fields.TimeSeriesLastUpdateDateField(length),
	}
}

func (t TimeSeries) Frames(_ context.Context, _ resource.ResourceProvider) (data.Frames, error) {

	length := len(t.TimeSeriesSummaries)

	timeSeriesSummaryFields := newTimeSeriesSummaryFields(length)

	for i, timeSeries := range t.TimeSeriesSummaries {

	
		if (timeSeries.Alias != nil) {
			timeSeriesSummaryFields.alias.Set(i, *timeSeries.Alias)
		} 
		if (timeSeries.AssetId != nil) {
			timeSeriesSummaryFields.assetId.Set(i, *timeSeries.AssetId)
		}
		if (timeSeries.DataType != nil) {
			timeSeriesSummaryFields.dataType.Set(i, *timeSeries.DataType)
		}
		if (timeSeries.DataTypeSpec != nil) {
			timeSeriesSummaryFields.dataTypeSpec.Set(i, *timeSeries.DataTypeSpec)
		}
		if (timeSeries.PropertyId != nil) {
			timeSeriesSummaryFields.propertyId.Set(i, *timeSeries.PropertyId)
		}
		if (timeSeries.TimeSeriesArn != nil) {
			timeSeriesSummaryFields.timeSeriesArn.Set(i, *timeSeries.TimeSeriesArn)
		}
		if (timeSeries.TimeSeriesId != nil) {
			timeSeriesSummaryFields.timeSeriesId.Set(i, *timeSeries.TimeSeriesId)
		}
		if (timeSeries.TimeSeriesCreationDate != nil) {
			timeSeriesSummaryFields.timeSeriesCreationDate.Set(i, *timeSeries.TimeSeriesCreationDate)
		}
		if (timeSeries.TimeSeriesLastUpdateDate != nil) {
			timeSeriesSummaryFields.timeSeriesLastUpdateDate.Set(i, *timeSeries.TimeSeriesLastUpdateDate)
		}
	}

	frame := data.NewFrame("", timeSeriesSummaryFields.fields()...)

	frame.Meta = &data.FrameMeta{
		Custom: models.SitewiseCustomMeta{
			NextToken: aws.StringValue(t.NextToken),
		},
	}

	return data.Frames{frame}, nil
}
