package sitewise

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type EdgeAuthenticator struct {
	Settings models.AWSSiteWiseDataSourceSetting
}

func (a *EdgeAuthenticator) Authorize(ctx context.Context) (models.AuthInfo, error) {
	reqBody := map[string]string{
		"username":      a.Settings.EdgeAuthUser,
		"password":      a.Settings.EdgeAuthPass,
		"authMechanism": a.Settings.EdgeAuthMode,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return models.AuthInfo{}, err
	}

	pool, _ := x509.SystemCertPool()
	if pool == nil {
		pool = x509.NewCertPool()
	}

	if a.Settings.Cert == "" {
		return models.AuthInfo{}, fmt.Errorf("certificate cannot be null")
	}

	block, _ := pem.Decode([]byte(a.Settings.Cert))
	if block == nil || block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return models.AuthInfo{}, fmt.Errorf("decode certificate failed: %s", a.Settings.Cert)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return models.AuthInfo{}, err
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
						return fmt.Errorf("tls: failed to parse certificate from server: " + err.Error())
					}
					certs[i] = cert
				}

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

	client := &http.Client{Transport: tr}

	authEndpoint := a.Settings.AWSDatasourceSettings.Endpoint + "authenticate"
	resp, err := client.Post(authEndpoint, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return models.AuthInfo{}, err
	}

	log.DefaultLogger.Debug("edge auth response: %s", resp)

	res := make(map[string]string)
	json.NewDecoder(resp.Body).Decode(&res)

	timeLayout := time.RFC3339
	sessionExpiryTime, err := time.Parse(timeLayout, res["sessionExpiryTime"])
	if err != nil {
		return models.AuthInfo{}, err
	}

	authInfo := models.AuthInfo{
		Username:          res["username"],
		AccessKeyId:       res["accessKeyId"],
		SecretAccessKey:   res["secretAccessKey"],
		SessionToken:      res["sessionToken"],
		SessionExpiryTime: sessionExpiryTime,
		AuthMechanism:     res["authMechanism"],
	}
	return authInfo, nil
}
