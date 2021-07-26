package sitewise

import (
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSimpleAuth(t *testing.T) {
	t.Skip()
	// This test is only meant for the testing of the edge authentication
	// during dev work, hence is skipped. To test with this, enter the
	// appropriate authMode, username, password, endpoint and cert
	settings := models.AWSSiteWiseDataSourceSetting{
		EdgeAuthMode: "linux",
		EdgeAuthUser: "username",
		EdgeAuthPass: "password",
	}

	settings.Endpoint = "https://localhost:80"
	settings.Cert = ``

	a := EdgeAuthenticator{Settings: settings}
	info, err := a.Authenticate()

	require.NoError(t, err)
	require.Equal(t, settings.EdgeAuthMode, info.AuthMechanism)
	require.Equal(t, settings.EdgeAuthUser, info.Username)
}
