package sitewise

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSimpleAuth(t *testing.T) {
	//t.Skip()
	// This test is only meant for the testing of the edge authentication
	// during dev work, hence is skipped. To test with this, enter the
	// appropriate authMode, username, password, endpoint and cert
	settings := models.AWSSiteWiseDataSourceSetting{
		EdgeAuthMode: "linux",
		EdgeAuthUser: "username",
		EdgeAuthPass: "password",
	}

	settings.Endpoint = "https://localhost:80"
	settings.Cert = `
-----BEGIN CERTIFICATE-----
MIICsDCCAhmgAwIBAgIJALwzrJEIBOaeMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTEwOTMwMTUyNjM2WhcNMjEwOTI3MTUyNjM2WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQC88Ckwru9VR2p2KJ1WQyqesLzr95taNbhkYfsd0j8Tl0MGY5h+dczCaMQz0YY3
xHXuU5yAQQTZjiks+D3KA3cx+iKDf2p1q77oXxQcx5CkrXBWTaX2oqVtHm3aX23B
AIORGuPk00b4rT3cld7VhcEFmzRNbyI0EqLMAxIwceUKSQIDAQABo4GnMIGkMB0G
A1UdDgQWBBSGmOdvSXKXclic5UOKPW35JLMEEjB1BgNVHSMEbjBsgBSGmOdvSXKX
clic5UOKPW35JLMEEqFJpEcwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUt
U3RhdGUxITAfBgNVBAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZIIJALwzrJEI
BOaeMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAcPfWn49pgAX54ji5
SiUPFFNCuQGSSTHh2I+TMrs1G1Mb3a0X1dV5CNLRyXyuVxsqhiM/H2veFnTz2Q4U
wdY/kPxE19Auwcz9AvCkw7ol1LIlLfJvBzjzOjEpZJNtkXTx8ROSooNrDeJl3HyN
cciS5hf80XzIFqwhzaVS9gmiyM8=
-----END CERTIFICATE-----
	`

	a := EdgeAuthenticator{Settings: settings}
	info, err := a.Authenticate()

	require.NoError(t, err)
	require.Equal(t, settings.EdgeAuthMode, info.AuthMechanism)
	require.Equal(t, settings.EdgeAuthUser, info.Username)
}

func TestAuthWithServer(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	settings := models.AWSSiteWiseDataSourceSetting{
		EdgeAuthMode: "linux",
		EdgeAuthUser: "username",
		EdgeAuthPass: "password",
	}

	settings.Endpoint = ts.URL
	settings.Cert = ``

	a := EdgeAuthenticator{Settings: settings}
	info, err := a.Authenticate()

	require.NoError(t, err)
	require.Equal(t, settings.EdgeAuthMode, info.AuthMechanism)
	require.Equal(t, settings.EdgeAuthUser, info.Username)
}
