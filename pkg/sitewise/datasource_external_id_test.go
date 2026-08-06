package sitewise

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/grafana/grafana-aws-sdk/pkg/awsauth"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type spyConfigProvider struct {
	captured awsauth.Settings
}

func (s *spyConfigProvider) GetConfig(_ context.Context, settings awsauth.Settings) (aws.Config, error) {
	s.captured = settings
	return aws.Config{}, nil
}

func TestGetClient_passesGrafanaExternalIDFields(t *testing.T) {
	spy := &spyConfigProvider{}
	origProvider := newAWSConfigProvider
	newAWSConfigProvider = func() awsauth.ConfigProvider { return spy }
	t.Cleanup(func() { newAWSConfigProvider = origProvider })

	usePerDS := true
	ds := &Datasource{
		Cfg: models.AWSSiteWiseDataSourceSetting{
			AWSDatasourceSettings: awsds.AWSDatasourceSettings{
				AuthType:                   awsds.AuthTypeGrafanaAssumeRole,
				Region:                     "us-east-1",
				AssumeRoleARN:              "arn:aws:iam::123456789012:role/test",
				GrafanaExternalID:          "stackABC-dsUid1",
				UsePerDatasourceExternalID: &usePerDS,
			},
		},
	}
	_, err := ds.getClient(context.Background(), "us-east-1")
	require.NoError(t, err)
	assert.Equal(t, "stackABC-dsUid1", spy.captured.GrafanaExternalID)
	require.NotNil(t, spy.captured.UsePerDatasourceExternalID)
	assert.True(t, *spy.captured.UsePerDatasourceExternalID)
}
