package sitewise

import (
	"context"
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSimpleAuth(t *testing.T) {
	t.Skip()
	settings := models.AWSSiteWiseDataSourceSetting{
		EdgeAuthMode: "linux",
		EdgeAuthUser: "username",
		EdgeAuthPass: "password",
	}

	settings.Endpoint = "https://localhost:80"
	settings.Cert = ``

	a := EdgeAuthenticator{Settings: settings}
	info, err := a.Authorize(context.Background())

	require.NoError(t, err)
	require.Equal(t, settings.EdgeAuthMode, info.AuthMechanism)
	require.Equal(t, settings.EdgeAuthUser, info.Username)
}
