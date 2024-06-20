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
		var cfg awsds.GetSessionConfig
		var as awsds.AuthSettings
		// deliberately error to avoid the rest of the function
		mockProvider := func(c awsds.GetSessionConfig, authSettings awsds.AuthSettings) (*session.Session, error) {
			cfg = c
			as = authSettings
			return nil, fmt.Errorf("break")
		}

		region := "us-east-1"
		baseConfig := models.AWSSiteWiseDataSourceSetting{
			AWSDatasourceSettings: awsds.AWSDatasourceSettings{Region: "us-west-1"},
		}
		authSettings := awsds.AuthSettings{
			AllowedAuthProviders: []string{"keys"},
		}
		_, err := GetClient(region, baseConfig, mockProvider, &authSettings)
		assert.Error(t, err)
		assert.Equal(t, cfg.Settings.Region, region)
		assert.Equal(t, []string{"keys"}, as.AllowedAuthProviders)
	})
}
