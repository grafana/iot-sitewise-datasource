package sitewise

import (
	"context"
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSimpleAuth(t *testing.T) {
	settings := models.AWSSiteWiseDataSourceSetting{
		EdgeAuthMode: "linux",
		EdgeAuthUser: "test-swe-admin",
		EdgeAuthPass: "password",
	}

	settings.Endpoint = "https://54.213.46.117:443/"

	a := EdgeAuthenticator{Settings: settings}
	info, err := a.Authorize(context.Background())

	require.NoError(t, err)
	require.Equal(t, settings.EdgeAuthMode, info.AuthMechanism)
	require.Equal(t, settings.EdgeAuthUser, info.Username)
}
