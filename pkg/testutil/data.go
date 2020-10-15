package testutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer/fdata"
)

var UnmarshallFileContents = func(filename string, val interface{}) error {

	cwd, _ := os.Getwd()
	fmt.Println(cwd)

	b, err := ioutil.ReadFile("../testdata/" + filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, val); err != nil {
		return err
	}
	return nil
}

var GetAssetProp = func(t *testing.T, filename string) iotsitewise.DescribeAssetPropertyOutput {
	property := iotsitewise.DescribeAssetPropertyOutput{}
	err := UnmarshallFileContents(filename, &property)
	if err != nil {
		t.Fatal(err)
	}
	return property
}

var GetPropVals = func(t *testing.T, filename string) fdata.AssetPropertyValue {
	propVals := fdata.AssetPropertyValue{}
	err := UnmarshallFileContents(filename, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var GetPropHistoryVals = func(t *testing.T, filename string) fdata.AssetPropertyValueHistory {
	propVals := fdata.AssetPropertyValueHistory{}
	err := UnmarshallFileContents(filename, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}

var GetAssetPropAggregates = func(t *testing.T, filename string) fdata.AssetPropertyAggregates {
	propVals := fdata.AssetPropertyAggregates{}
	err := UnmarshallFileContents(filename, &propVals)
	if err != nil {
		t.Fatal(err)
	}
	return propVals
}
