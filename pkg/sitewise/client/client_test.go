package client

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("Uses region", func(t *testing.T) {
		var cfg awsds.SessionConfig
		// deliberately error to avoid the rest of the function
		mockProvider := func(c awsds.SessionConfig) (*session.Session, error) {
			cfg = c
			return nil, fmt.Errorf("break")
		}

		region := "us-east-1"
		baseConfig := models.AWSSiteWiseDataSourceSetting{
			AWSDatasourceSettings: awsds.AWSDatasourceSettings{Region: "us-west-1"},
		}
		_, err := GetClient(region, baseConfig, mockProvider)
		assert.Error(t, err)
		assert.Equal(t, cfg.Settings.Region, region)
	})
}
