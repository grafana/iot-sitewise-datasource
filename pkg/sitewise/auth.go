package sitewise

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type EdgeAuthenticator struct {
	Settings models.AWSSiteWiseDataSourceSetting
	authInfo *models.AuthInfo
}

type AuthRequest struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	AuthMechanism string `json:"authMechanism,omitempty"`
}

func (a *EdgeAuthenticator) GetAuthInfo() (*models.AuthInfo, error) {
	if a == nil {
		return nil, nil
	}
	if a.authInfo == nil || time.Now().After(a.authInfo.SessionExpiryTime) {
		err := a.Authenticate()
		if err != nil {
			return nil, err
		}
	}
	return a.authInfo, nil

}

func (a *EdgeAuthenticator) Authenticate() error {
	if a == nil {
		return nil
	}
	reqBodyJson, err := json.Marshal(
		&AuthRequest{
			Username:      a.Settings.EdgeAuthUser,
			Password:      a.Settings.EdgeAuthPass,
			AuthMechanism: a.Settings.EdgeAuthMode,
		})
	if err != nil {
		return err
	}

	pool, _ := x509.SystemCertPool()
	if pool == nil {
		pool = x509.NewCertPool()
	}

	if a.Settings.Cert == "" {
		return fmt.Errorf("certificate cannot be null")
	}

	block, _ := pem.Decode([]byte(a.Settings.Cert))
	if block == nil || block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return fmt.Errorf("decode certificate failed: %s", a.Settings.Cert)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	pool.AddCert(cert)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //Not actually skipping, check the cert in VerifyPeerCertificate
			RootCAs:            pool,
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				// If this is the first handshake on a connection, process and
				// (optionally) verify the server's certificates.
				certs := make([]*x509.Certificate, len(rawCerts))
				for i, asn1Data := range rawCerts {
					cert, err := x509.ParseCertificate(asn1Data)
					if err != nil {
						return fmt.Errorf("tls: failed to parse certificate from server: %w", err)
					}
					certs[i] = cert
				}

				// see: https://github.com/golang/go/issues/21971
				opts := x509.VerifyOptions{
					Roots:         pool,
					CurrentTime:   time.Now(),
					DNSName:       "", // <- skip hostname verification
					Intermediates: x509.NewCertPool(),
				}

				for i, cert := range certs {
					if i == 0 {
						continue
					}
					opts.Intermediates.AddCert(cert)
				}
				_, err := certs[0].Verify(opts)
				return err
			},
		},
	}

	client := &http.Client{Transport: tr, Timeout: time.Second * 5}

	u, err := url.Parse(a.Settings.Endpoint)
	if err != nil {
		log.DefaultLogger.Error("error parsing edge endpoint url.", "endpoint url:", a.Settings.Endpoint)
		return fmt.Errorf("cannot parse edge endpoint url. url: %v", a.Settings.Endpoint)
	}
	u.Path = path.Join(u.Path, "authenticate")
	authEndpoint := u.String()

	resp, err := client.Post(authEndpoint, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.DefaultLogger.Error("edge auth response not ok:", "response code:", strconv.Itoa(resp.StatusCode))
		return fmt.Errorf("request not ok. returned code: %v", resp.StatusCode)
	}

	log.DefaultLogger.Debug("edge auth response ok.")

	authInfo := models.AuthInfo{}
	err = json.NewDecoder(resp.Body).Decode(&authInfo)
	if err != nil {
		return err
	}
	a.authInfo = &authInfo
	return nil
}

type DummyAuthenticator struct {
	Settings models.AWSSiteWiseDataSourceSetting
}

func (a *DummyAuthenticator) Authenticate() (models.AuthInfo, error) {
	if rand.Float64() > .8 {
		return models.AuthInfo{}, fmt.Errorf("dummy auth failed (1/5) chance of that")
	}

	return models.AuthInfo{
		AccessKeyId:       a.Settings.AccessKey,
		SecretAccessKey:   a.Settings.SecretKey,
		SessionToken:      "",
		SessionExpiryTime: time.Now().Add(20 * time.Second),
	}, nil
}
