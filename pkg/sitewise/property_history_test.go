package sitewise

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

func generatePropertyHistoryTestData(t *testing.T, client client.Client) interface{} {
	var (
		ctx = context.Background()
	)

	// hard coded values from my account
	query := models.AssetPropertyValueQuery{}
	query.AssetId = testAssetId
	query.PropertyId = testPropIdAvgWind
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
