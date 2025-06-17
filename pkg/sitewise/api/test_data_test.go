// nolint
package api

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/testdata"
)

const (
	SKIPALL = true
)

type testDataFunc func(t *testing.T, client client.SitewiseAPIClient) interface{}

// How to run tests:
//
// Use shell environment variables. Ex:
// export AWS_ACCESS_KEY_ID="<key id>"
// export AWS_SECRET_ACCESS_KEY="<secret key>"
// export AWS_SESSION_TOKEN="<session token>"
func TestGenerateTestData(t *testing.T) {

	if SKIPALL {
		t.Skip("Integration Test")
	}

	m := make(map[string]testDataFunc)

	m["property-history-values.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()

		// hard coded values from my account
		query := models.AssetPropertyValueQuery{}
		query.AssetIds = []string{testdata.DemoTurbineAsset1}
		query.PropertyIds = []string{testdata.TurbinePropAvgWindSpeed}
		query.TimeRange = backend.TimeRange{
			From: time.Now().Add(time.Hour * -3), // return 3 hours of data. 60*3/5 = 36 points
			To:   time.Now(),
		}
		query.MaxPageAggregations = 1

		_, resp, err := BatchGetAssetPropertyValues(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}

		return resp
	}
	m["property-history-values-with-alias.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()

		// hard coded values from my account
		query := models.AssetPropertyValueQuery{}
		query.AssetIds = []string{"709e2e02-7b28-4b41-b669-3fb501a11853"}
		//query.PropertyId = "5ff66b29-5b79-427a-978b-29f8dfc2757a"
		query.PropertyAliases = []string{"/amazon/renton/1/rpm"}
		query.TimeRange = backend.TimeRange{
			From: time.Now().Add(time.Hour * -3), // return 3 hours of data. 60*3/5 = 36 points
			To:   time.Now(),
		}
		query.MaxPageAggregations = 1

		_, resp, err := BatchGetAssetPropertyValues(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}

		return resp
	}
	m["property-value.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()

		query := models.AssetPropertyValueQuery{}
		query.AssetIds = []string{testdata.DemoTurbineAsset1}
		query.PropertyIds = []string{testdata.TurbinePropAvgWindSpeed}

		_, resp, err := BatchGetAssetPropertyValue(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}

		return resp
	}
	m["property-aggregate-values.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()

		query := models.AssetPropertyValueQuery{}
		query.Resolution = "1m"
		query.AggregateTypes = []types.AggregateType{
			models.AggregateCount, models.AggregateAvg, models.AggregateMin, models.AggregateMax, models.AggregateSum, models.AggregateStdDev,
		}
		query.AssetIds = []string{testdata.DemoTurbineAsset1}
		query.PropertyIds = []string{testdata.TurbinePropWindSpeed}
		query.TimeRange = backend.TimeRange{
			From: time.Now().Add(time.Hour * -3), // return 3 hours of data. 60*3/5 = 36 points
			To:   time.Now(),
		}
		query.MaxPageAggregations = 1

		_, resp, err := GetAssetPropertyAggregates(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}

		return resp
	}
	m["describe-asset.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetQuery{}
		query.AssetIds = []string{testdata.DemoTurbineAsset1}

		resp, err := DescribeAsset(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}
	m["describe-asset-top-level.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetQuery{}
		query.AssetIds = []string{testdata.DemoWindFarmAssetId}

		resp, err := DescribeAsset(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}
	m["describe-asset-property-avg-wind.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetPropertyQuery{}
		query.AssetIds = []string{testdata.DemoTurbineAsset1}
		query.PropertyIds = []string{testdata.TurbinePropAvgWindSpeed}
		resp, err := GetAssetPropertyDescription(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["describe-asset-property-raw-wind.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.DescribeAssetPropertyQuery{}
		query.AssetIds = []string{testdata.DemoTurbineAsset1}
		query.PropertyIds = []string{testdata.TurbinePropWindSpeed}
		resp, err := GetAssetPropertyDescription(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-asset-models.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		resp, err := ListAssetModels(ctx, client, models.ListAssetModelsQuery{})
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-assets.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.ListAssetsQuery{}
		query.ModelId = testdata.DemoTurbineAssetModelId
		query.Filter = "ALL"
		resp, err := ListAssets(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-assets-top-level.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.ListAssetsQuery{}
		resp, err := ListAssets(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-associated-assets.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.ListAssociatedAssetsQuery{}
		query.AssetIds = []string{testdata.DemoWindFarmAssetId}
		query.HierarchyId = testdata.TurbineAssetModelHierarchyId
		resp, err := ListAssociatedAssets(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["list-associated-assets-parent.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this
		ctx := context.Background()
		query := models.ListAssociatedAssetsQuery{}
		query.AssetIds = []string{testdata.DemoTurbineAsset1}
		resp, err := ListAssociatedAssets(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	m["describe-asset-model.json"] = func(t *testing.T, client client.SitewiseAPIClient) interface{} {
		t.Skip("Integration Test") // comment line to run this

		ctx := context.Background()
		query := models.DescribeAssetModelQuery{}
		query.AssetModelId = testdata.DemoTurbineAssetModelId
		resp, err := DescribeAssetModel(ctx, client, query)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}

	sw := client.NewSitewiseClientForRegion("us-east-1")

	for k, v := range m {
		writeTestData(t, k, v, sw)
	}
}

func writeTestData(t *testing.T, filename string, tf testDataFunc, client client.SitewiseAPIClient) {

	t.Run(filename, func(t *testing.T) {
		resp := tf(t, client)

		js, err := json.MarshalIndent(resp, "", "    ")

		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Create("../../testdata/" + filename)

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
