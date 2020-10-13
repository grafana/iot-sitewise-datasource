package sitewise

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	framerImpl "github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type apiFunc func(context.Context, client.Client, models.AssetPropertyValueQuery) (framer.FrameData, error)

type testScenario struct {
	name         string
	query        models.AssetPropertyValueQuery
	propVals     framer.FrameData
	property     iotsitewise.DescribeAssetPropertyOutput
	validationFn func(t *testing.T, frames data.Frames)
}

var getPropVals = func(t *testing.T, filename string) AssetPropertyValue {
	propVals := AssetPropertyValue{}
	err := unmarshallFileContents(filename, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var getPropHistoryVals = func(t *testing.T, filename string) AssetPropertyValueHistory {
	propVals := AssetPropertyValueHistory{}
	err := unmarshallFileContents(filename, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var getAssetPropAggregates = func(t *testing.T, filename string) AssetPropertyAggregates {
	propVals := AssetPropertyAggregates{}
	err := unmarshallFileContents(filename, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var getAssetProp = func(t *testing.T, filename string) iotsitewise.DescribeAssetPropertyOutput {
	property := iotsitewise.DescribeAssetPropertyOutput{}
	err := unmarshallFileContents(filename, &property)
	if err != nil {
		t.Fatal(err)
	}
	return property
}

// fieldAssert will verify the field created by the framer contains the expected information.
// As we add additional field config + tags, expand this struct.
type fieldAssert struct {
	fields       data.Fields
	idx          int
	expectedName string
	expectedType data.FieldType
}

func (fa fieldAssert) assert(t *testing.T) {
	field := fa.fields[fa.idx]
	assert.Equal(t, fa.expectedName, field.Name, "wrong name for field in Field[%d]. got %s expected %s", fa.idx, field.Name, fa.expectedName)
	assert.Equal(t, fa.expectedType, field.Type(), "wrong type for field in Field[%d]. got %s expected %s", fa.idx, field.Type(), fa.expectedType)
}

var assertFramesAndGetFields = func(t *testing.T, frames data.Frames) data.Fields {
	assert.Len(t, frames, 1)
	frame := frames[0]
	t.Log(frame.StringTable(-1, -1))

	return frame.Fields
}

func getScenarios(t *testing.T) []*testScenario {

	return []*testScenario{
		{
			name: "TestAssetPropertyValue",
			query: models.AssetPropertyValueQuery{
				QueryType:  models.QueryTypePropertyValue,
				AssetId:    testAssetId,
				PropertyId: testPropIdAvgWind,
			},
			propVals: getPropVals(t, "property-value.json"),
			property: getAssetProp(t, "describe-asset-property-avg-wind.json"),
			validationFn: func(t *testing.T, frames data.Frames) {

				fields := assertFramesAndGetFields(t, frames)

				fieldAssert{
					fields:       fields,
					idx:          0,
					expectedName: "time",
					expectedType: data.FieldTypeInt64,
				}.assert(t)

				fieldAssert{
					fields:       fields,
					idx:          1,
					expectedName: "Average Wind Speed",
					expectedType: data.FieldTypeNullableFloat64,
				}.assert(t)

			},
		},
		{
			name: "TestNullResponseAssetPropertyValues",
			query: models.AssetPropertyValueQuery{
				QueryType:  models.QueryTypePropertyValue,
				AssetId:    testAssetId,
				PropertyId: testPropIdAvgWind,
			},
			propVals: AssetPropertyValue{
				PropertyValue: &iotsitewise.AssetPropertyValue{
					Quality: aws.String("GOOD"),
					Timestamp: &iotsitewise.TimeInNanos{
						OffsetInNanos: aws.Int64(0),
						TimeInSeconds: aws.Int64(1602219000),
					},
					Value: &iotsitewise.Variant{
						BooleanValue: nil,
						DoubleValue:  nil,
						IntegerValue: nil,
						StringValue:  nil,
					},
				},
			},
			property: getAssetProp(t, "describe-asset-property-avg-wind.json"),
			validationFn: func(t *testing.T, frames data.Frames) {
				fields := assertFramesAndGetFields(t, frames)
				fieldAssert{
					fields:       fields,
					idx:          0,
					expectedName: "time",
					expectedType: data.FieldTypeInt64,
				}.assert(t)

				fieldAssert{
					fields:       fields,
					idx:          1,
					expectedName: "Average Wind Speed",
					expectedType: data.FieldTypeNullableFloat64,
				}.assert(t)

			},
		},
		{
			name: "TestAssetPropertyHistoryValues",
			query: models.AssetPropertyValueQuery{
				QueryType:  models.QueryTypePropertyValueHistory,
				AssetId:    testAssetId,
				PropertyId: testPropIdAvgWind,
			},
			propVals: getPropHistoryVals(t, "property-history-values.json"),
			property: getAssetProp(t, "describe-asset-property-avg-wind.json"),
			validationFn: func(t *testing.T, frames data.Frames) {

				fields := assertFramesAndGetFields(t, frames)

				fieldAssert{
					fields:       fields,
					idx:          0,
					expectedName: "time",
					expectedType: data.FieldTypeInt64,
				}.assert(t)

				fieldAssert{
					fields:       fields,
					idx:          1,
					expectedName: "Average Wind Speed",
					expectedType: data.FieldTypeNullableFloat64,
				}.assert(t)

			},
		},
		{
			name: "TestAssetPropertyHistoryAggregates",
			query: models.AssetPropertyValueQuery{
				AssetId:        testAssetId,
				PropertyId:     testPropIdRawWin,
				AggregateTypes: []string{models.AggregateMax, models.AggregateMin, models.AggregateAvg},
				Resolution:     "1m",
				QueryType:      models.QueryTypePropertyAggregate,
			},
			propVals: getAssetPropAggregates(t, "property-aggregate-values.json"),
			property: getAssetProp(t, "describe-asset-property-raw-wind.json"),
			validationFn: func(t *testing.T, frames data.Frames) {

				fields := assertFramesAndGetFields(t, frames)

				// time, avg, min, max
				assert.Len(t, fields, 4, "expected [time, avg, min, max]")

				fieldAssert{
					fields:       fields,
					idx:          0,
					expectedName: "time",
					expectedType: data.FieldTypeInt64,
				}.assert(t)

				fieldAssert{
					fields:       fields,
					idx:          1,
					expectedName: "avg",
					expectedType: data.FieldTypeNullableFloat64,
				}.assert(t)

				fieldAssert{
					fields:       fields,
					idx:          2,
					expectedName: "min",
					expectedType: data.FieldTypeNullableFloat64,
				}.assert(t)

				fieldAssert{
					fields:       fields,
					idx:          3,
					expectedName: "max",
					expectedType: data.FieldTypeNullableFloat64,
				}.assert(t)

			},
		},
	}
}

func TestFrameData(t *testing.T) {

	for _, v := range getScenarios(t) {

		t.Run(v.name, func(t *testing.T) {

			var ctx = context.Background()

			sw := &mocks.Client{}
			sw.On("DescribeAssetPropertyWithContext", mock.Anything, mock.Anything).Return(&v.property, nil)

			rp := resource.NewSitewiseResourceProvider(sw)
			mp := framerImpl.NewPropertyValueMetaProvider(rp, v.query)
			fr := framerImpl.PropertyValueQueryFramer{
				FrameData:    v.propVals,
				MetaProvider: mp,
				Request:      v.query,
			}

			dataFrames, err := fr.Frames(ctx)

			if err != nil {
				t.Fatal(err)
			}

			v.validationFn(t, dataFrames)

		})

	}

}
