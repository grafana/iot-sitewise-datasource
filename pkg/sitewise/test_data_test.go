package sitewise

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

const (
	testAssetId       = "a9fe4e4a-e028-4be2-bd15-2f8dd0bee23b"
	testPropIdAvgWind = "1e1e256e-e32a-4666-8aeb-22b4131192eb"
	testPropIdRawWin  = "bfaa662d-0eb2-49d2-a24d-2dad5a75bfde"
)

type testDataFunc func(t *testing.T, client client.Client) interface{}

// How to run tests:
//
// Use shell environment variables. Ex:
// export AWS_ACCESS_KEY_ID="<key id>"
// export AWS_SECRET_ACCESS_KEY="<secret key>"
// export AWS_SESSION_TOKEN="<session token>"
//
func TestGenerateTestData(t *testing.T) {

	t.Skip("Integration Test") // comment line to run this

	m := make(map[string]testDataFunc)

	m["property-history-values.json"] = generatePropertyHistoryTestData
	m["property-value.json"] = generatePropertyValueTestData
	m["property-aggregate-values.json"] = generatePropertyAggregateTestData
	m["describe-asset.json"] = func(t *testing.T, client client.Client) interface{} {
		ctx := context.Background()
		resp, err := GetAssetDescription(ctx, client, models.DescribeAssetQuery{AssetId: testAssetId})
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}
	m["describe-asset-property-avg-wind.json"] = func(t *testing.T, client client.Client) interface{} {
		ctx := context.Background()
		resp, err := GetAssetPropertyDescription(ctx, client, models.DescribeAssetPropertyQuery{
			AssetId:    testAssetId,
			PropertyId: testPropIdAvgWind,
		})
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["describe-asset-property-raw-wind.json"] = func(t *testing.T, client client.Client) interface{} {
		ctx := context.Background()
		resp, err := GetAssetPropertyDescription(ctx, client, models.DescribeAssetPropertyQuery{
			AssetId:    testAssetId,
			PropertyId: testPropIdRawWin,
		})
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	sesh := session.Must(session.NewSession())
	sw := iotsitewise.New(sesh, aws.NewConfig().WithRegion("us-east-1"))

	for k, v := range m {
		writeTestData(k, v, sw, t)
	}
}

func writeTestData(filename string, tf testDataFunc, client client.Client, t *testing.T) {

	resp := tf(t, client)

	js, err := json.MarshalIndent(resp, "", "    ")
	//js, err := jsonutil.BuildJSON(resp)

	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create("./testdata/" + filename)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = f.Write(js)
	if err != nil {
		t.Fatal(err)
	}

}
