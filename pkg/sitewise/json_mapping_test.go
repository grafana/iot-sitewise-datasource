package sitewise

import (
	"context"
	"math"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

// testResourceProvider implements the full ResourceProvider interface minimally for tests.
type testResourceProvider struct {
	assetOut      *iotsitewise.DescribeAssetOutput
	assetsOut     map[string]*iotsitewise.DescribeAssetOutput
	propertyOut   *iotsitewise.DescribeAssetPropertyOutput
	propertiesOut map[string]*iotsitewise.DescribeAssetPropertyOutput
	modelOut      *iotsitewise.DescribeAssetModelOutput
	err           error
}

func (t *testResourceProvider) Asset(ctx context.Context) (*iotsitewise.DescribeAssetOutput, error) {
	return t.assetOut, t.err
}

func (t *testResourceProvider) Assets(ctx context.Context) (map[string]*iotsitewise.DescribeAssetOutput, error) {
	return t.assetsOut, t.err
}

func (t *testResourceProvider) Property(ctx context.Context) (*iotsitewise.DescribeAssetPropertyOutput, error) {
	return t.propertyOut, t.err
}

func (t *testResourceProvider) Properties(ctx context.Context) (map[string]*iotsitewise.DescribeAssetPropertyOutput, error) {
	return t.propertiesOut, t.err
}

func (t *testResourceProvider) AssetModel(ctx context.Context) (*iotsitewise.DescribeAssetModelOutput, error) {
	return t.modelOut, t.err
}

// helpers for tests
func getFieldByName(frame *data.Frame, name string) *data.Field {
	for _, f := range frame.Fields {
		if f != nil && f.Name == name {
			return f
		}
	}
	return nil
}

func floatEquals(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestBuildPropertyNameMap_EmptyAssetID(t *testing.T) {
	// When assetID is empty, BuildPropertyNameMap returns empty map and doesn't call AssetModel.
	out := BuildPropertyNameMap(context.Background(), nil, "")
	if len(out) != 0 {
		t.Fatalf("expected empty map for empty assetID, got %v", out)
	}
}

func TestBuildPropertyNameMap_WithModel(t *testing.T) {
	provider := &testResourceProvider{
		modelOut: &iotsitewise.DescribeAssetModelOutput{
			AssetModelProperties: []iotsitewisetypes.AssetModelProperty{
				{Id: aws.String("id1"), Name: aws.String("Pressure")},
				{Id: aws.String("id2"), Name: aws.String("Temperature")},
			},
		},
		err: nil,
	}

	out := BuildPropertyNameMap(context.Background(), provider, "asset-1")
	if len(out) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(out))
	}
	if out["id1"] != "Pressure" || out["id2"] != "Temperature" {
		t.Fatalf("unexpected mapping result: %v", out)
	}
}

func TestParseJSONFields_SimpleNumericAndDiagnostic(t *testing.T) {
	ctx := context.Background()

	js := `{
  "timestamp":"2025-11-13T06:42:21.549955Z",
  "prediction":0,
  "prediction_reason":"NO_ANOMALY_DETECTED",
  "anomaly_score":0.4312,
  "some_value": 12.34,
  "flag": true,
  "diagnostics": [
   {"name":"root\\prop-1","value":45.8,"anomaly_score":0.1},
   {"name":"root\\prop-2","value":54.2}
  ]
 }`

	stringField := data.NewField("value", nil, []string{js})
	frame := data.NewFrame("TestFrame", stringField)
	frames := data.Frames{frame}

	// pass nil provider with empty assetID — BuildPropertyNameMap returns empty map early
	outFrames := ParseJSONFields(ctx, frames, (resource.ResourceProvider)(nil), "")
	if len(outFrames) != 1 {
		t.Fatalf("expected 1 output frame, got %d", len(outFrames))
	}
	of := outFrames[0]

	// timestamp should be skipped (original string field preserved but timestamp not added as a numeric field)
	if getFieldByName(of, "timestamp") != nil || getFieldByName(of, "Timestamp") != nil {
		t.Fatalf("timestamp should not be present as a parsed numeric/string field")
	}

	// check anomaly_score field exists and value set
	asField := getFieldByName(of, "anomaly_score")
	if asField == nil {
		t.Fatalf("anomaly_score field missing, fields: %v", func() []string {
			names := []string{}
			for _, f := range of.Fields {
				if f != nil {
					names = append(names, f.Name)
				}
			}
			return names
		}())
	}
	if v := asField.At(0).(float64); !floatEquals(v, 0.4312, 1e-9) {
		t.Fatalf("unexpected anomaly_score: %v", v)
	}

	// check simple fields some_value and flag
	valField := getFieldByName(of, "some_value")
	if valField == nil {
		t.Fatalf("expected some_value field")
	}
	if got := valField.At(0).(float64); !floatEquals(got, 12.34, 1e-9) {
		t.Fatalf("unexpected some_value: %v", got)
	}

	flagField := getFieldByName(of, "flag")
	if flagField == nil {
		t.Fatalf("expected flag field")
	}
	if got := flagField.At(0).(bool); !got {
		t.Fatalf("expected flag true")
	}

	// contributions should be normalized to percentages
	c1 := getFieldByName(of, "contrib_prop-1")
	c2 := getFieldByName(of, "contrib_prop-2")
	if c1 == nil || c2 == nil {
		t.Fatalf("expected contrib fields present, fields: %v", func() []string {
			names := []string{}
			for _, f := range of.Fields {
				if f != nil {
					names = append(names, f.Name)
				}
			}
			return names
		}())
	}
	v1 := c1.At(0).(float64)
	v2 := c2.At(0).(float64)
	if !floatEquals(v1+v2, 100.0, 1e-6) {
		t.Fatalf("contribs not normalized: %v + %v != 100", v1, v2)
	}

	// diag anomaly field from first diag
	diagField := getFieldByName(of, "diag_anomaly_prop-1")
	if diagField == nil {
		t.Fatalf("diag anomaly field missing")
	}
	if got := diagField.At(0).(float64); !floatEquals(got, 0.1, 1e-9) {
		t.Fatalf("unexpected diag anomaly value: %v", got)
	}
}

func TestParseJSONFields_WithAssetNameMapping(t *testing.T) {
	ctx := context.Background()

	modelOut := &iotsitewise.DescribeAssetModelOutput{
		AssetModelProperties: []iotsitewisetypes.AssetModelProperty{
			{Id: aws.String("prop-1"), Name: aws.String("pressure")},
			{Id: aws.String("prop-2"), Name: aws.String("temperature")},
		},
	}
	provider := &testResourceProvider{
		modelOut: modelOut,
		err:      nil,
	}

	js := `{
  "timestamp":"2025-11-13T06:42:21.549955Z",
  "prediction":1,
  "prediction_reason":"ANOMALY_DETECTED",
  "anomaly_score":0.97,
  "diagnostics":[
   {"name":"x\\prop-1","value":10,"anomaly_score":0.11},
   {"name":"x\\prop-2","value":90,"anomaly_score":0.89}
  ]
 }`

	frame := data.NewFrame("MappedFrame", data.NewField("value", nil, []string{js}))
	in := data.Frames{frame}

	outFrames := ParseJSONFields(ctx, in, provider, "asset-123")
	if len(outFrames) != 1 {
		t.Fatalf("expected 1 output frame, got %d", len(outFrames))
	}
	of := outFrames[0]

	// with mapping, contrib_pressure and contrib_temperature should be present
	cPressure := getFieldByName(of, "contrib_pressure")
	cTemp := getFieldByName(of, "contrib_temperature")
	if cPressure == nil || cTemp == nil {
		t.Fatalf("expected contrib_pressure & contrib_temperature present, got: %v", func() []string {
			names := []string{}
			for _, f := range of.Fields {
				if f != nil {
					names = append(names, f.Name)
				}
			}
			return names
		}())
	}
	pv := cPressure.At(0).(float64)
	tv := cTemp.At(0).(float64)
	if !floatEquals(pv+tv, 100.0, 1e-6) {
		t.Fatalf("mapped contribs not normalized: %v + %v != 100", pv, tv)
	}

	// diag anomaly fields should be created with mapped names
	da1 := getFieldByName(of, "diag_anomaly_pressure")
	da2 := getFieldByName(of, "diag_anomaly_temperature")
	if da1 == nil || da2 == nil {
		t.Fatalf("expected diag anomaly fields for mapped names")
	}
	if !floatEquals(da1.At(0).(float64), 0.11, 1e-9) || !floatEquals(da2.At(0).(float64), 0.89, 1e-9) {
		t.Fatalf("unexpected diag anomaly values: %v, %v", da1.At(0), da2.At(0))
	}
}

func TestParseJSONFields_CorruptedJSONAndMissingFields(t *testing.T) {
	ctx := context.Background()
	corrupt := `{"invalid":`
	missingReq := `{"some":"value"}`

	frame := data.NewFrame("BadFrame", data.NewField("value", nil, []string{corrupt, missingReq}))
	in := data.Frames{frame}
	provider := &testResourceProvider{
		modelOut: &iotsitewise.DescribeAssetModelOutput{},
		err:      nil,
	}

	outFrames := ParseJSONFields(ctx, in, provider, "")
	if len(outFrames) != 1 {
		t.Fatalf("expected 1 output frame despite bad rows, got %d", len(outFrames))
	}
	of := outFrames[0]

	// should still produce a frame with at least the original string field present
	if len(of.Fields) == 0 {
		t.Fatalf("expected at least one field in output frame, got none")
	}

	// should not have created anomaly_score/contrib fields for the broken rows
	if getFieldByName(of, "anomaly_score") != nil {
		t.Fatalf("did not expect anomaly_score field for corrupted rows")
	}
	if getFieldByName(of, "contrib_foo") != nil {
		t.Fatalf("did not expect contrib_foo field")
	}
	if getFieldByName(of, "timestamp") != nil {
		t.Fatalf("did not expect timestamp field")
	}
}
