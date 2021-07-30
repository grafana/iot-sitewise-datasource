package sitewise

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	ts := httptest.NewUnstartedServer(http.HandlerFunc(dummyHandler))

	// https TLS cert
	rootCertPEM, rootTLSCert, err := createTLSCert()
	if err != nil {
		t.Fatal(fmt.Errorf("generating TLS certificate: %v", err))
	}

	// Configure the server to present the TLS certificate we created
	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{rootTLSCert},
	}

	ts.StartTLS()
	defer ts.Close()

	// test client
	settings := models.AWSSiteWiseDataSourceSetting{
		EdgeAuthMode: "linux",
		EdgeAuthUser: "username",
		EdgeAuthPass: "password",
	}

	settings.Endpoint = ts.URL
	settings.Cert = string(rootCertPEM)

	a := EdgeAuthenticator{Settings: settings}
	info, err := a.Authenticate()

	require.NoError(t, err)
	require.Equal(t, settings.EdgeAuthMode, info.AuthMechanism)
	require.Equal(t, settings.EdgeAuthUser, info.Username)
}

// helper function for the https request handler
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	// dummy /authenticate endpoint
	if r.Method == "POST" && r.RequestURI == "/authenticate" {
		authReq := AuthRequest{}
		err := json.NewDecoder(r.Body).Decode(&authReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// dummy response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["username"] = authReq.Username
		resp["accessKeyId"] = "dummyAccessKeyIdYEWZb5yVBl9llM9TQvn10hD4wmXKlUCNgXeCQY5YmssV55Fz"
		resp["secretAccessKey"] = "dummySecretAccessKey2wH5XvUVv2FKIxvvj3YNCblvMJkI67KbXZV6ZHiy2w16"
		resp["sessionToken"] = "dummySessionTokenglPPiSJwMx3iDuLm5BsVJVA0t5wXVhMNHFyaOkh68yz48V9"
		resp["sessionExpiryTime"] = "2019-07-29T20:29:41.176Z"
		resp["authMechanism"] = authReq.AuthMechanism
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

// helper function to create a certificate template with a serial number and other required fields
func certTemplate() (*x509.Certificate, error) {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Test Inc."}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour), // valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

//helper function to create a certificate from a template and public key plus a parent certificate and private key
func createCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	*x509.Certificate, []byte, error) {

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return nil, nil, errors.New("failed to create a certificate: " + err.Error())
	}

	// parse the resulting certificate so we can use it again
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, errors.New("failed to parse the certificate: " + err.Error())
	}

	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM := pem.EncodeToMemory(&b)

	return cert, certPEM, err
}

// helper function to create a TLS Certificate and Root Key Combination
func createTLSCert() ([]byte, tls.Certificate, error) {
	// generate a new key-pair for the test server TLS
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, tls.Certificate{}, fmt.Errorf("generating random key: %v", err)
	}

	// generate a certificate template for the test server TLS
	rootCertTmpl, err := certTemplate()
	if err != nil {
		return nil, tls.Certificate{}, fmt.Errorf("creating cert template: %v", err)
	}

	// set the certificate to be used for TLS handshake authentication
	rootCertTmpl.IsCA = true
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	rootCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")} // set the host address here

	// create a self-signed certificate for the test server TLS
	_, rootCertPEM, err := createCert(rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey)
	if err != nil {
		return nil, tls.Certificate{}, fmt.Errorf("error creating cert: %v", err)
	}

	// PEM encode the private key for TLS server handshake
	rootKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey),
	})

	// Create a TLS certificate using the private key and certificate
	rootTLSCert, err := tls.X509KeyPair(rootCertPEM, rootKeyPEM)
	if err != nil {
		return nil, tls.Certificate{}, fmt.Errorf("invalid key pair: %v", err)
	}

	return rootCertPEM, rootTLSCert, nil
}
