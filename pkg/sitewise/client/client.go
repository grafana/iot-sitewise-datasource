//go:generate mockery --name Client

package client

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/aws/aws-sdk-go/service/iotsitewise/iotsitewiseiface"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	gaws "github.com/grafana/iot-sitewise-datasource/pkg/common/aws"
)

type Client interface {
	iotsitewiseiface.IoTSiteWiseAPI
}

type clientCache map[string]Client

var cache = make(clientCache)
var lock sync.RWMutex

func GetClient(ctx backend.PluginContext, region string) (client Client, err error) {
	lock.Lock()
	if sw, ok := cache[region]; ok {
		client = sw
	} else {
		client, err = initClient(ctx, region)
		if client != nil {
			cache[region] = client
		}
	}
	lock.Unlock()
	return
}

func initClient(ctx backend.PluginContext, region string) (Client, error) {
	settings, err := gaws.LoadSettings(*ctx.DataSourceInstanceSettings)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	if region == "" {
		region = settings.DefaultRegion
	}

	cfg, err := gaws.GetAwsConfig(settings, region)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	swcfg := &aws.Config{}
	if settings.Endpoint != "" {
		swcfg.Endpoint = aws.String(settings.Endpoint)
	}

	return iotsitewise.New(sess, swcfg), nil
}
