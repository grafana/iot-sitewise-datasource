package sitewise

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

const (
	SKIPALL = false
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

	if SKIPALL {
		t.Skip("Integration Test")
	}

	m := make(map[string]testDataFunc)

	m["property-history-values.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()

		// hard coded values from my account
		query := models.AssetPropertyValueQuery{}
		query.AssetId = testdata.TestAssetId
		query.PropertyId = testdata.TestPropIdAvgWind
		query.TimeRange = backend.TimeRange{
			From: time.Now().Add(time.Hour * -3), // return 3 hours of data. 60*3/5 = 36 points
			To:   time.Now(),
		}

		resp, err := GetAssetPropertyValues(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}

		return resp
	}
	m["property-value.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()

		query := models.AssetPropertyValueQuery{}
		query.AssetId = testdata.TestAssetId
		query.PropertyId = testdata.TestPropIdAvgWind

		resp, err := GetAssetPropertyValue(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}

		return resp
	}
	m["property-aggregate-values.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()

		query := models.AssetPropertyValueQuery{}
		query.Resolution = "1m"
		query.AggregateTypes = []string{
			models.AggregateCount, models.AggregateAvg, models.AggregateMin, models.AggregateMax, models.AggregateSum, models.AggregateStdDev,
		}
		query.AssetId = testdata.TestAssetId
		query.PropertyId = testdata.TestPropIdRawWin
		query.TimeRange = backend.TimeRange{
			From: time.Now().Add(time.Hour * -3), // return 3 hours of data. 60*3/5 = 36 points
			To:   time.Now(),
		}

		resp, err := GetAssetPropertyAggregates(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}

		return resp
	}
	m["describe-asset.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetQuery{}
		query.AssetId = testdata.TestAssetId

		resp, err := DescribeAsset(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}
	m["describe-asset-top-level.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetQuery{}
		query.AssetId = testdata.TestTopLevelAssetId

		resp, err := DescribeAsset(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}
	m["describe-asset-property-avg-wind.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetPropertyQuery{}
		query.AssetId = testdata.TestAssetId
		query.PropertyId = testdata.TestPropIdAvgWind
		resp, err := GetAssetPropertyDescription(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["describe-asset-property-raw-wind.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetPropertyQuery{}
		query.AssetId = testdata.TestAssetId
		query.PropertyId = testdata.TestPropIdRawWin
		resp, err := GetAssetPropertyDescription(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-asset-models.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		resp, err := ListAssetModels(ctx, client, models.ListAssetModelsQuery{})
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-assets.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.ListAssetsQuery{}
		query.ModelId = testdata.TestAssetModelId
		query.Filter = "ALL"
		resp, err := ListAssets(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-assets-top-level.json"] = func(t *testing.T, client client.Client) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.ListAssetsQuery{}
		resp, err := ListAssets(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-associated-assets.json"] = func(t *testing.T, client client.Client) interface{} {
		//t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.ListAssociatedAssetsQuery{}
		query.AssetId = testdata.TestTopLevelAssetId
		query.HierarchyId = testdata.TestTopLevelAssetHierarchyId
		query.TraversalDirection = "CHILD"
		resp, err := ListAssociatedAssets(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	sesh := session.Must(session.NewSession())
	sw := iotsitewise.New(sesh, aws.NewConfig().WithRegion("us-east-1"))

	for k, v := range m {
		writeTestData(t, k, v, sw)
	}
}

func writeTestData(t *testing.T, filename string, tf testDataFunc, client client.Client) {

	t.Run(filename, func(t *testing.T) {
		resp := tf(t, client)

		js, err := json.MarshalIndent(resp, "", "    ")

		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Create("../testdata/" + filename)

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
	})

}
