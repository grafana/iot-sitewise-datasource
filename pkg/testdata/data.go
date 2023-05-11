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

func SerializeStruct(t *testing.T, val interface{}) []byte {
	vbytes, err := json.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}
	return vbytes
}

func UnmarshalFileContents(path string, val interface{}) error {
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

func GetIotSitewiseAssetProp(t *testing.T, path string) iotsitewise.DescribeAssetPropertyOutput {
	property := iotsitewise.DescribeAssetPropertyOutput{}
	err := UnmarshalFileContents(path, &property)
	if err != nil {
		t.Fatal(err)
	}
	return property
}

func GetPropVals(t *testing.T, path string) framer.AssetPropertyValue {
	propVals := framer.AssetPropertyValue{}
	err := UnmarshalFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

func GetPropHistoryVals(t *testing.T, path string) framer.AssetPropertyValueHistory {
	propVals := framer.AssetPropertyValueHistory{}
	err := UnmarshalFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

func GetIoTSitewisePropHistoryVals(t *testing.T, path string) iotsitewise.BatchGetAssetPropertyValueHistoryOutput {
	propVals := iotsitewise.BatchGetAssetPropertyValueHistoryOutput{}
	err := UnmarshalFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

func GetIoTSitewisePropVal(t *testing.T, path string) iotsitewise.BatchGetAssetPropertyValueOutput {
	propVal := iotsitewise.BatchGetAssetPropertyValueOutput{}
	err := UnmarshalFileContents(path, &propVal)
	if err != nil {
		t.Fatal(err)
	}
	return propVal
}

func GetIoTSitewisePropAggregateVals(t *testing.T, path string) iotsitewise.BatchGetAssetPropertyAggregatesOutput {
	propAggs := iotsitewise.BatchGetAssetPropertyAggregatesOutput{}
	err := UnmarshalFileContents(path, &propAggs)
	if err != nil {
		t.Fatal(err)
	}
	return propAggs
}

func GetAssetPropAggregates(t *testing.T, path string) framer.AssetPropertyAggregates {
	propVals := framer.AssetPropertyAggregates{}
	err := UnmarshalFileContents(path, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

func GetIoTSitewiseAssetModels(t *testing.T, path string) iotsitewise.ListAssetModelsOutput {
	assetModels := iotsitewise.ListAssetModelsOutput{}
	err := UnmarshalFileContents(path, &assetModels)
	if err != nil {
		t.Fatal(err)
	}
	return assetModels
}

func GetIoTSitewiseAssets(t *testing.T, path string) iotsitewise.ListAssetsOutput {
	assets := iotsitewise.ListAssetsOutput{}
	err := UnmarshalFileContents(path, &assets)
	if err != nil {
		t.Fatal(err)
	}
	return assets
}

func GetIoTSitewiseAssetDescription(t *testing.T, path string) iotsitewise.DescribeAssetOutput {
	asset := iotsitewise.DescribeAssetOutput{}
	err := UnmarshalFileContents(path, &asset)
	if err != nil {
		t.Fatal(err)
	}
	return asset
}

func GetIoTSitewiseAssociatedAssets(t *testing.T, path string) iotsitewise.ListAssociatedAssetsOutput {
	assets := iotsitewise.ListAssociatedAssetsOutput{}
	err := UnmarshalFileContents(path, &assets)
	if err != nil {
		t.Fatal(err)
	}
	return assets
}

func GetIoTSitewiseAssetModelDescription(t *testing.T, path string) iotsitewise.DescribeAssetModelOutput {
	assets := iotsitewise.DescribeAssetModelOutput{}
	err := UnmarshalFileContents(path, &assets)
	if err != nil {
		t.Fatal(err)
	}
	return assets
}

func GetIoTSitewiseTimeSeries(t *testing.T, path string) iotsitewise.DescribeTimeSeriesOutput {
	timeSeries := iotsitewise.DescribeTimeSeriesOutput{}
	err := UnmarshalFileContents(path, &timeSeries)
	if err != nil {
		t.Fatal(err)
	}
	return timeSeries
}
