package testdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/framer"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

var SerializeStruct = func(t *testing.T, val interface{}) []byte {
	vbytes, err := json.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}
	return vbytes
}

var UnmarshallFileContents = func(path string, val interface{}) error {
	cwd, _ := os.Getwd()
	fmt.Println(cwd)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, val); err != nil {
		return err
	}
	return nil
}

var GetIotSitewiseAssetProp = func(t *testing.T, path string) iotsitewise.DescribeAssetPropertyOutput {
	property := iotsitewise.DescribeAssetPropertyOutput{}
	err := UnmarshallFileContents(path, &property)
	if err != nil {
		t.Fatal(err)
	}
	return property
}

var GetPropVals = func(t *testing.T, path string) framer.AssetPropertyValue {
	propVals := framer.AssetPropertyValue{}
	err := UnmarshallFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var GetPropHistoryVals = func(t *testing.T, path string) framer.AssetPropertyValueHistory {
	propVals := framer.AssetPropertyValueHistory{}
	err := UnmarshallFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var GetIoTSitewisePropHistoryVals = func(t *testing.T, path string) iotsitewise.GetAssetPropertyValueHistoryOutput {
	propVals := iotsitewise.GetAssetPropertyValueHistoryOutput{}
	err := UnmarshallFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var GetIoTSitewisePropVal = func(t *testing.T, path string) iotsitewise.GetAssetPropertyValueOutput {
	propVal := iotsitewise.GetAssetPropertyValueOutput{}
	err := UnmarshallFileContents(path, &propVal)
	if err != nil {
		t.Fatal(err)
	}
	return propVal
}

var GetIoTSitewisePropAggregateVals = func(t *testing.T, path string) iotsitewise.GetAssetPropertyAggregatesOutput {
	propAggs := iotsitewise.GetAssetPropertyAggregatesOutput{}
	err := UnmarshallFileContents(path, &propAggs)
	if err != nil {
		t.Fatal(err)
	}
	return propAggs
}

var GetAssetPropAggregates = func(t *testing.T, path string) framer.AssetPropertyAggregates {
	propVals := framer.AssetPropertyAggregates{}
	err := UnmarshallFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var GetIoTSitewiseAssetModels = func(t *testing.T, path string) iotsitewise.ListAssetModelsOutput {
	assetModels := iotsitewise.ListAssetModelsOutput{}
	err := UnmarshallFileContents(path, &assetModels)
	if err != nil {
		t.Fatal(err)
	}
	return assetModels
}

var GetIoTSitewiseAssets = func(t *testing.T, path string) iotsitewise.ListAssetsOutput {
	assets := iotsitewise.ListAssetsOutput{}
	err := UnmarshallFileContents(path, &assets)
	if err != nil {
		t.Fatal(err)
	}
	return assets
}

var GetIoTSitewiseAssetDescription = func(t *testing.T, path string) iotsitewise.DescribeAssetOutput {
	asset := iotsitewise.DescribeAssetOutput{}
	err := UnmarshallFileContents(path, &asset)
	if err != nil {
		t.Fatal(err)
	}
	return asset
}
