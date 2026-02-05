package sitewise

import (
	"context"
	"math"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type testResourceProvider struct {
	assetOut *iotsitewise.DescribeAssetOutput
	modelOut *iotsitewise.DescribeAssetModelOutput
	err      error
}

func (t *testResourceProvider) Asset(ctx context.Context) (*iotsitewise.DescribeAssetOutput, error) {
	return t.assetOut, t.err
}

func (t *testResourceProvider) Assets(ctx context.Context) (map[string]*iotsitewise.DescribeAssetOutput, error) {
	return nil, t.err
}

func (t *testResourceProvider) Property(ctx context.Context) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	return nil, t.err
}

func (t *testResourceProvider) Properties(ctx context.Context) (map[string]*iotsitewise.DescribeAssetPropertyOutput, error) {
	return nil, t.err
}

func (t *testResourceProvider) AssetModel(ctx context.Context) (*iotsitewise.DescribeAssetModelOutput, error) {
	return t.modelOut, t.err
}

func getField(frame *data.Frame, name string) *data.Field {
	for _, f := range frame.Fields {
		if f != nil && f.Name == name {
			return f
		}
	}
	return nil
}

func floatEq(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

func TestBuildPropertyNameMap_EmptyAssetID(t *testing.T) {
	out, err := BuildPropertyNameMap(context.Background(), nil, "")
	if err == nil {
		t.Fatalf("expected error for empty assetID")
	}
	if out != nil {
		t.Fatalf("expected nil map")
	}
}

func TestBuildPropertyNameMap_ErrorFromModel(t *testing.T) {
	p := &testResourceProvider{err: context.Canceled}

	out, err := BuildPropertyNameMap(context.Background(), p, "asset-1")
	if err == nil {
		t.Fatalf("expected error")
	}
	if out != nil {
		t.Fatalf("expected nil map")
	}
}

func TestBuildPropertyNameMap_NoProperties(t *testing.T) {
	p := &testResourceProvider{
		modelOut: &iotsitewise.DescribeAssetModelOutput{},
	}

	out, err := BuildPropertyNameMap(context.Background(), p, "asset-1")
	if err == nil {
		t.Fatalf("expected error when no properties")
	}
	if out != nil {
		t.Fatalf("expected nil map")
	}
}

func TestBuildPropertyNameMap_Success(t *testing.T) {
	p := &testResourceProvider{
		modelOut: &iotsitewise.DescribeAssetModelOutput{
			AssetModelProperties: []iotsitewisetypes.AssetModelProperty{
				{Id: aws.String("p1"), Name: aws.String("Pressure")},
				{Id: aws.String("p2"), Name: aws.String("Temperature")},
			},
		},
	}

	out, err := BuildPropertyNameMap(context.Background(), p, "asset-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries")
	}
}

func TestRequiresJsonParsing(t *testing.T) {
	if !requiresJsonParsing(models.BaseQuery{QueryType: models.QueryTypePropertyValue}) {
		t.Fatalf("expected true")
	}
	if !requiresJsonParsing(models.BaseQuery{QueryType: models.QueryTypePropertyAggregate}) {
		t.Fatalf("expected true")
	}
	if !requiresJsonParsing(models.BaseQuery{QueryType: models.QueryTypePropertyValueHistory}) {
		t.Fatalf("expected true")
	}
	if requiresJsonParsing(models.BaseQuery{}) {
		t.Fatalf("expected false for unknown query type")
	}

}

func TestParseJSONFields_HappyPath_WithMappingAndAssetName(t *testing.T) {
	ctx := context.Background()

	p := &testResourceProvider{
		assetOut: &iotsitewise.DescribeAssetOutput{
			AssetName: aws.String("Pump1"),
		},
		modelOut: &iotsitewise.DescribeAssetModelOutput{
			AssetModelProperties: []iotsitewisetypes.AssetModelProperty{
				{Id: aws.String("prop-1"), Name: aws.String("pressure")},
				{Id: aws.String("prop-2"), Name: aws.String("temperature")},
			},
		},
	}

	js := `{
  "timestamp":1,
  "prediction":"ANOMALY",
  "prediction_reason":"threshold",
  "anomaly_score":0.8,
  "value":12.3,
  "diagnostics":[
   {"name":"x\\prop-1","value":20},
   {"name":"x\\prop-2","value":80}
  ]
 }`

	frame := data.NewFrame("f", data.NewField("raw", nil, []string{js}))
	out := ParseJSONFields(ctx, data.Frames{frame}, p, "asset-1")
	of := out[0]

	if getField(of, "prediction") == nil {
		t.Fatalf("prediction missing")
	}

	if !floatEq(getField(of, "anomaly_score").At(0).(float64), 0.8) {
		t.Fatalf("wrong anomaly_score")
	}

	c1 := getField(of, "contrib_Pump1_pressure")
	c2 := getField(of, "contrib_Pump1_temperature")
	if c1 == nil || c2 == nil {
		t.Fatalf("missing contrib fields")
	}
	if !floatEq(c1.At(0).(float64)+c2.At(0).(float64), 100) {
		t.Fatalf("contrib not normalized")
	}
}

func TestParseJSONFields_CorruptAndNonJSON(t *testing.T) {
	ctx := context.Background()

	frame := data.NewFrame(
		"bad",
		data.NewField("v", nil, []string{
			`{"invalid":`,
			"hello",
		}),
	)

	out := ParseJSONFields(ctx, data.Frames{frame}, &testResourceProvider{
		modelOut: &iotsitewise.DescribeAssetModelOutput{},
	}, "asset-1")

	if len(out[0].Fields) != 1 {
		t.Fatalf("expected only original field")
	}
}

func TestParseJSONFields_AssetNameFallback(t *testing.T) {
	ctx := context.Background()

	p := &testResourceProvider{
		assetOut: &iotsitewise.DescribeAssetOutput{},
		modelOut: &iotsitewise.DescribeAssetModelOutput{
			AssetModelProperties: []iotsitewisetypes.AssetModelProperty{
				{Id: aws.String("p1"), Name: aws.String("speed")},
			},
		},
	}

	js := `{
  "timestamp":1,
  "prediction":"OK",
  "prediction_reason":"none",
  "diagnostics":[{"name":"x\\p1","value":100}]
 }`

	frame := data.NewFrame("f", data.NewField("v", nil, []string{js}))
	out := ParseJSONFields(ctx, data.Frames{frame}, p, "asset-X")

	if getField(out[0], "contrib_asset-X_speed") == nil {
		t.Fatalf("expected fallback assetID field")
	}
}
