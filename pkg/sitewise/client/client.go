//go:generate mockery --name SitewiseClient

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/aws/aws-sdk-go/service/iotsitewise/iotsitewiseiface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type SitewiseClient interface {
	iotsitewiseiface.IoTSiteWiseAPI
	GetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyValueHistoryOutput, error)
	GetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyAggregatesOutput, error)
}

type sitewiseClient struct {
	iotsitewiseiface.IoTSiteWiseAPI
}

// NewSitewiseClient is mainly for testing in this case
func NewSitewiseClientForRegion(region string) SitewiseClient {
	sesh := session.Must(session.NewSession())
	sw := iotsitewise.New(sesh, aws.NewConfig().WithRegion(region))
	return &sitewiseClient{
		sw,
	}
}

func (c *sitewiseClient) GetAssetPropertyValueHistoryPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyValueHistoryInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyValueHistoryOutput, error) {
	var (
		numPages  = 0
		values    []*iotsitewise.AssetPropertyValue
		nextToken *string
	)

	err := c.GetAssetPropertyValueHistoryPagesWithContext(ctx, req, func(output *iotsitewise.GetAssetPropertyValueHistoryOutput, b bool) bool {
		numPages++
		values = append(values, output.AssetPropertyValueHistory...)
		nextToken = output.NextToken
		return numPages < maxPages && len(values) <= maxResults
	})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.GetAssetPropertyValueHistoryOutput{
		AssetPropertyValueHistory: values,
		NextToken:                 nextToken,
	}, nil
}

func (c *sitewiseClient) GetAssetPropertyAggregatesPageAggregation(ctx context.Context, req *iotsitewise.GetAssetPropertyAggregatesInput, maxPages int, maxResults int) (*iotsitewise.GetAssetPropertyAggregatesOutput, error) {

	var (
		numPages  = 0
		values    []*iotsitewise.AggregatedValue
		nextToken *string
	)

	err := c.GetAssetPropertyAggregatesPagesWithContext(ctx, req, func(output *iotsitewise.GetAssetPropertyAggregatesOutput, b bool) bool {
		numPages++
		values = append(values, output.AggregatedValues...)
		nextToken = output.NextToken
		return numPages < maxPages && len(values) <= maxResults
	})

	if err != nil {
		return nil, err
	}

	return &iotsitewise.GetAssetPropertyAggregatesOutput{
		AggregatedValues: values,
		NextToken:        nextToken,
	}, nil
}

func GetClient(region string, settings models.AWSSiteWiseDataSourceSetting, provider awsds.AmazonSessionProvider) (client SitewiseClient, err error) {
	sess, err := provider(region, settings.ToAWSDatasourceSettings())
	if err != nil {
		return nil, err
	}

	swcfg := &aws.Config{}
	if settings.Endpoint != "" {
		swcfg.Endpoint = aws.String(settings.Endpoint)
	}

	if settings.Region == models.EDGE_REGION {
		pool, _ := x509.SystemCertPool()
		if pool == nil {
			pool = x509.NewCertPool()
		}

		if settings.Cert == "" {
			return nil, errors.New("certificate cannot be null")
		}

		block, _ := pem.Decode([]byte(settings.Cert))
		if block == nil || block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			return nil, fmt.Errorf("decode certificate failed: %s", settings.Cert)
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
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
							return errors.New("tls: failed to parse certificate from server: " + err.Error())
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
		httpClient := &http.Client{Transport: tr}

		swcfg = swcfg.WithHTTPClient(httpClient)
		swcfg = swcfg.WithDisableEndpointHostPrefix(true)

	}

	c := iotsitewise.New(sess, swcfg)

	c.Handlers.Send.PushFront(func(r *request.Request) {
		r.HTTPRequest.Header.Set("User-Agent", userAgentString())
	})
	return &sitewiseClient{c}, nil
}

// TODO, move to https://github.com/grafana/grafana-plugin-sdk-go
func userAgentString() string {
	return fmt.Sprintf("%s/%s (%s; %s) Grafana/%s", aws.SDKName, aws.SDKVersion, runtime.Version(), runtime.GOOS, os.Getenv("GF_VERSION"))
}
