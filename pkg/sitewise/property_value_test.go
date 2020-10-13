package sitewise

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/resource"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var unmarshallFileContents = func(filename string, val interface{}) error {

	b, err := ioutil.ReadFile("./testdata/" + filename)
	if err != nil {
		return err
	}

	//fmt.Println(string(b))
	if err := json.Unmarshal(b, val); err != nil {
		return err
	}
	//if err := jsonutil.UnmarshalJSON(val, bytes.NewReader(b)); err != nil {
	//	return err
	//}
	return nil
}

func mockPropertyValue(mockSw *mocks.Client) {

	resp := iotsitewise.GetAssetPropertyValueOutput{}
	err := unmarshallFileContents("property-value.json", &resp)

	mockSw.On(
		"GetAssetPropertyValueWithContext", mock.Anything, mock.Anything,
	).Return(&resp, err)

	fmt.Println(resp)
}

func mockDescribeProperty(mockSw *mocks.Client) {

	resp := iotsitewise.DescribeAssetPropertyOutput{}
	err := unmarshallFileContents("describe-asset-property-avg-wind.json", &resp)

	mockSw.On(
		"DescribeAssetPropertyWithContext", mock.Anything, mock.Anything,
	).Return(&resp, err)
}

func TestPropertyValue(t *testing.T) {

	var (
		ctx = context.Background()
	)

	sw := &mocks.Client{}
	mockPropertyValue(sw)
	mockDescribeProperty(sw)

	query := models.AssetPropertyValueQuery{
		QueryType:  models.QueryTypePropertyValue,
		AssetId:    testAssetId,
		PropertyId: testPropIdAvgWind,
	}

	fd, err := GetAssetPropertyValue(ctx, sw, query)
	if err != nil {
		t.Fatal(err)
	}

	rp := resource.NewSitewiseResourceProvider(sw)
	mp := framer.NewPropertyValueMetaProvider(rp, query)
	fr := framer.PropertyValueQueryFramer{
		FrameData:    fd,
		MetaProvider: mp,
		Request:      query,
	}
	dataFrames, err := fr.Frames(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, dataFrames, 1)

	frame := dataFrames[0]
	t.Log(frame.StringTable(-1, -1))

	if val, ok := frame.Fields[0].At(0).(int64); ok {
		t.Log("time: ", val)
		assert.True(t, ok)
	}

	if val, ok := frame.Fields[1].At(0).(float64); ok {
		t.Log("value: ", val)
		assert.True(t, ok)
	}
}

func generatePropertyValueTestData(t *testing.T, client client.Client) interface{} {

	var (
		ctx = context.Background()
	)

	query := models.AssetPropertyValueQuery{}
	query.AssetId = testAssetId
	query.PropertyId = testPropIdAvgWind

	resp, err := GetAssetPropertyValue(ctx, client, query)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}
