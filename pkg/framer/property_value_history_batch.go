package framer

import (
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fields"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

type AssetPropertyValueHistoryBatch struct {
	Responses       []*iotsitewise.BatchGetAssetPropertyValueHistoryOutput
	Query           models.AssetPropertyValueQuery
	AnomalyAssetIds []string
	SitewiseClient  client.SitewiseClient
}

func (p AssetPropertyValueHistoryBatch) Frames(ctx context.Context, resources resource.ResourceProvider) (data.Frames, error) {
	successEntriesLength := 0
	for _, r := range p.Responses {
		successEntriesLength += len(r.SuccessEntries)
	}
	frames := make(data.Frames, 0, successEntriesLength)

	properties, err := resources.Properties(ctx)
	if err != nil {
		return frames, err
	}

	for _, r := range p.Responses {
		for _, s := range r.SuccessEntries {
			frame, err := p.Frame(ctx, properties[*s.EntryId], s.AssetPropertyValueHistory)
			frame.Meta = &data.FrameMeta{
				Custom: models.SitewiseCustomMeta{
					NextToken:  aws.StringValue(r.NextToken),
					EntryId:    *s.EntryId,
					Resolution: models.PropertyQueryResolutionRaw,
				},
			}
			if err != nil {
				return nil, err
			}
			if frame != nil {
				frames = append(frames, frame)
			}
		}

		for _, e := range r.ErrorEntries {
			property := properties[*e.EntryId]
			frame := data.NewFrame(getFrameName(property))
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

func (p AssetPropertyValueHistoryBatch) Frame(ctx context.Context, property *iotsitewise.DescribeAssetPropertyOutput, h []*iotsitewise.AssetPropertyValue) (*data.Frame, error) {
	length := len(h)
	// TODO: make this work with the API instead of ad-hoc dataType inference
	// https://github.com/grafana/iot-sitewise-datasource/issues/98#issuecomment-892947756
	if util.IsAssetProperty(property) && *property.AssetProperty.DataType == *aws.String("?") {
		if length != 0 {
			property.AssetProperty.DataType = aws.String(getPropertyVariantValueType(h[0].Value))
		} else {
			property.AssetProperty.DataType = aws.String("")
		}
	}

	assetId := property.AssetId
	if assetId != nil && slices.Contains(p.AnomalyAssetIds, *assetId) {
		return p.frameL4ePropertyValues(ctx, property, h)
	} else {
		return p.framePropertyValues(property, h)
	}
}

// framePropertyValues creates a frame for a property value history.
func (p AssetPropertyValueHistoryBatch) framePropertyValues(property *iotsitewise.DescribeAssetPropertyOutput, h []*iotsitewise.AssetPropertyValue) (*data.Frame, error) {
	length := len(h)

	timeField := fields.TimeField(length)
	valueField := fields.PropertyValueFieldForQuery(p.Query, property, length)
	qualityField := fields.QualityField(length)
	frameName := ""
	if models.QueryTypePropertyAggregate == p.Query.QueryType {
		frameName = getFrameName(property)
	} else {
		frameName = *property.AssetName
	}
	frame := data.NewFrame(
		frameName,
		timeField,
		valueField,
		qualityField)

	for i, v := range h {
		if v.Value != nil && getPropertyVariantValue(v.Value) != nil {
			timeField.Set(i, getTime(v.Timestamp))
			valueField.Set(i, getPropertyVariantValue(v.Value))
			qualityField.Set(i, *v.Quality)
		}
	}

	return frame, nil
}

// frameL4ePropertyValues creates a frame for a property value history with L4E fields flatten.
func (p AssetPropertyValueHistoryBatch) frameL4ePropertyValues(ctx context.Context, property *iotsitewise.DescribeAssetPropertyOutput, h []*iotsitewise.AssetPropertyValue) (*data.Frame, error) {
	frameName := ""
	if models.QueryTypePropertyAggregate == p.Query.QueryType {
		frameName = getFrameName(property)
	} else {
		frameName = *property.AssetName
	}

	dataFields, err := p.parseL4eFields(ctx, property.AssetId, h)
	if err != nil {
		return nil, err
	}

	frame := data.NewFrame(
		frameName,
		dataFields...)

	return frame, nil
}

func (p AssetPropertyValueHistoryBatch) parseL4eFields(ctx context.Context, assetId *string, h []*iotsitewise.AssetPropertyValue) ([]*data.Field, error) {
	length := len(h)

	dataFields := []*data.Field{}

	timeField := fields.TimeField(length)
	dataFields = append(dataFields, timeField)

	qualityField := fields.QualityField(length)
	dataFields = append(dataFields, qualityField)

	anomalyScoreField := fields.AnomalyScoreField(length)
	dataFields = append(dataFields, anomalyScoreField)

	predictionReasonField := fields.PredictionReasonField(length)
	dataFields = append(dataFields, predictionReasonField)

	// Maps diagnostic property id to the corresponding data field
	diagnosticsMap := map[string]*data.Field{}

	for i, v := range h {
		var l4eAnomalyResult models.L4eAnomalyResult
		err := json.Unmarshal([]byte(*v.Value.StringValue), &l4eAnomalyResult)
		if err != nil {
			return nil, err
		}

		timeField.Set(i, getTime(v.Timestamp))
		qualityField.Set(i, *v.Quality)
		anomalyScoreField.Set(i, l4eAnomalyResult.AnomalyScore)
		predictionReasonField.Set(i, l4eAnomalyResult.PredictionReason)

		for _, diagnostics := range l4eAnomalyResult.Diagnostics {
			diagnosticsField, ok := diagnosticsMap[diagnostics.Name]
			if !ok {
				diagnosticsField = fields.DiagnosticField(length, diagnostics.Name)
				diagnosticsMap[diagnostics.Name] = diagnosticsField
				dataFields = append(dataFields, diagnosticsField)
			}
			diagnosticsField.Set(i, diagnostics.Value)
		}
	}

	// Rename diagnostic fields with human friendly names
	for _, diagnosticsField := range diagnosticsMap {
		propertyId := strings.Split(diagnosticsField.Name, "\\")[0]

		req := &iotsitewise.DescribeAssetPropertyInput{
			AssetId:    assetId,
			PropertyId: aws.String(propertyId),
		}
		resp, err := p.SitewiseClient.DescribeAssetPropertyWithContext(ctx, req)
		if err != nil {
			return nil, err
		}

		diagnosticsField.Name = *resp.AssetProperty.Name
	}

	return dataFields, nil
}
