package framer

import (
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

type AssetPropertyValueBatch struct {
	Responses       []*iotsitewise.BatchGetAssetPropertyValueOutput
	AnomalyAssetIds []string
	SitewiseClient  client.SitewiseClient
}

func (p AssetPropertyValueBatch) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	successEntriesLength := 0
	for _, r := range p.Responses {
		successEntriesLength += len(r.SuccessEntries)
	}
	frames := make(data.Frames, 0, successEntriesLength)

	properties, err := resources.Properties(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range p.Responses {
		for _, e := range r.SuccessEntries {
			property := properties[*e.EntryId]
			if util.IsAssetProperty(property) && !isPropertyDataTypeDefined(*property.AssetProperty.DataType) && e.AssetPropertyValue != nil {
				property.AssetProperty.DataType = aws.String(getPropertyVariantValueType(e.AssetPropertyValue.Value))
			}

			var frame *data.Frame
			if property.AssetId != nil && slices.Contains(p.AnomalyAssetIds, *property.AssetId) {
				frame, err = p.frameL4ePropertyValue(ctx, property, e.AssetPropertyValue)
				if err != nil {
					return nil, err
				}
			} else {
				frame = p.framePropertyValue(property, e.AssetPropertyValue)
			}
			frame.Meta = &data.FrameMeta{
				Custom: models.SitewiseCustomMeta{
					NextToken: aws.StringValue(r.NextToken),
					EntryId:   *e.EntryId,
				},
			}
			frames = append(frames, frame)
		}

		for _, e := range r.ErrorEntries {
			property := properties[*e.EntryId]
			frame := data.NewFrame(*property.AssetName)
			if e.ErrorMessage != nil {
				frame.Meta = &data.FrameMeta{
					Notices: []data.Notice{{Severity: data.NoticeSeverityError, Text: *e.ErrorMessage}},
				}
			}
			frames = append(frames, frame)
		}
	}

	return frames, nil
}

func (AssetPropertyValueBatch) framePropertyValue(property *iotsitewise.DescribeAssetPropertyOutput, assetPropertyValue *iotsitewise.AssetPropertyValue) *data.Frame {
	timeField := fields.TimeField(0)
	valueField := fields.PropertyValueField(property, 0)
	qualityField := fields.QualityField(0)

	frame := data.NewFrame(*property.AssetName, timeField, valueField, qualityField)

	if assetPropertyValue != nil && getPropertyVariantValue(assetPropertyValue.Value) != nil {
		timeField.Append(getTime(assetPropertyValue.Timestamp))
		valueField.Append(getPropertyVariantValue(assetPropertyValue.Value))
		qualityField.Append(*assetPropertyValue.Quality)
	}
	return frame
}

func (p AssetPropertyValueBatch) frameL4ePropertyValue(ctx context.Context, property *iotsitewise.DescribeAssetPropertyOutput, assetPropertyValue *iotsitewise.AssetPropertyValue) (*data.Frame, error) {
	dataFields := []*data.Field{}

	timeField := fields.TimeField(0)
	dataFields = append(dataFields, timeField)

	qualityField := fields.QualityField(0)
	dataFields = append(dataFields, qualityField)

	anomalyScoreField := fields.AnomalyScoreField(0)
	dataFields = append(dataFields, anomalyScoreField)

	predictionReasonField := fields.PredictionReasonField(0)
	dataFields = append(dataFields, predictionReasonField)

	if assetPropertyValue == nil {
		frame := data.NewFrame(*property.AssetName, dataFields...)
		return frame, nil
	}

	var l4eAnomalyResult models.L4eAnomalyResult
	err := json.Unmarshal([]byte(*assetPropertyValue.Value.StringValue), &l4eAnomalyResult)
	if err != nil {
		return nil, err
	}

	timeField.Append(getTime(assetPropertyValue.Timestamp))
	qualityField.Append(*assetPropertyValue.Quality)
	anomalyScoreField.Append(l4eAnomalyResult.AnomalyScore)
	predictionReasonField.Append(l4eAnomalyResult.PredictionReason)

	for _, diagnostics := range l4eAnomalyResult.Diagnostics {
		propertyId := strings.Split(diagnostics.Name, "\\")[0]

		req := &iotsitewise.DescribeAssetPropertyInput{
			AssetId:    property.AssetId,
			PropertyId: aws.String(propertyId),
		}
		resp, err := p.SitewiseClient.DescribeAssetPropertyWithContext(ctx, req)
		if err != nil {
			return nil, err
		}

		diagnosticsField := fields.DiagnosticField(0, *resp.AssetProperty.Name)
		diagnosticsField.Append(diagnostics.Value)
		dataFields = append(dataFields, diagnosticsField)
	}

	frame := data.NewFrame(*property.AssetName, dataFields...)

	return frame, nil
}
