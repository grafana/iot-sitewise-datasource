package framer

import (
	"context"
	"time"

	"github.com/grafana/iot-sitewise-datasource/pkg/util"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

// AssetPropertyValueMetadata handles creating Fields for all 3 data fetching APIS. A few TODOs/remarks:
// - Might be able to abstract the Fields() method out to another interface. That code path has the most branching.
// - Composing the Frame Meta is awkward with this abstraction, as it requires the response as well.
// -
type AssetPropertyValueMetadata struct {
	// asset may not be needed for property values
	property *iotsitewise.DescribeAssetPropertyOutput
	query    models.AssetPropertyValueQuery
}

// TODO: add Field Labels/Config
func (md AssetPropertyValueMetadata) Fields() ([]*data.Field, error) {

	fields := []*data.Field{
		data.NewField("time", nil, []time.Time{}),
	}

	qfields, err := md.getQueryTypeValueFields()

	if err != nil {
		return nil, err
	}

	fields = append(fields, qfields...)
	return fields, nil
}

func (md AssetPropertyValueMetadata) getQueryTypeValueFields() ([]*data.Field, error) {

	queryType := md.query.QueryType
	fieldName := *md.property.AssetProperty.Name

	switch queryType {
	case models.QueryTypePropertyAggregate:
		return getAggregationFields(md.query)
	case models.QueryTypePropertyValue:
		return []*data.Field{data.NewField(fieldName, nil, fieldTypeForPropertyValue(md.property))}, nil
	case models.QueryTypePropertyValueHistory:
		return []*data.Field{data.NewField(fieldName, nil, fieldTypeForPropertyValue(md.property))}, nil
	}

	return nil, nil
}

func fieldTypeForPropertyValue(property *iotsitewise.DescribeAssetPropertyOutput) interface{} {
	switch *property.AssetProperty.DataType {
	case "BOOLEAN":
		return []*bool{}
	case "DOUBLE":
		return []*float64{}
	case "INTEGER":
		return []*int64{}
	case "STRING":
		return []*string{}
	default:
		// todo: unsure what/if to default. Should never be any values outside these types
		return []*int64{}
	}
}

func getAggregationFields(query models.AssetPropertyValueQuery) ([]*data.Field, error) {
	var fields []*data.Field

	// convert the query aggregate params to a "set"
	aggregations := util.StringSliceToSet(query.AggregateTypes)

	for _, k := range models.AggregateOrder {
		// if the aggregate is in the "set", add to fields
		if _, found := aggregations[k]; found {
			if agg, ok := models.AggregateFields[k]; ok {
				fields = append(fields, data.NewField(agg.FieldName, nil, []*float64{}))
			}
		}

	}

	return fields, nil
}

func (md AssetPropertyValueMetadata) FrameName() string {
	return *md.property.AssetName
}

type propertyValueMetaProvider struct {
	resources resource.SitewiseResourceProvider
	query     models.AssetPropertyValueQuery
}

func NewPropertyValueMetaProvider(resource resource.SitewiseResourceProvider, query models.AssetPropertyValueQuery) *propertyValueMetaProvider {
	return &propertyValueMetaProvider{
		resources: resource,
		query:     query,
	}
}

func (qmp *propertyValueMetaProvider) Provide(ctx context.Context) (framer.Metadata, error) {

	assetId := qmp.query.AssetId
	propertyId := qmp.query.PropertyId

	property, err := qmp.resources.Property(ctx, assetId, propertyId)

	if err != nil {
		return nil, err
	}

	return &AssetPropertyValueMetadata{
		property: property,
		query:    qmp.query,
	}, nil
}
