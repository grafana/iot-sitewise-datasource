//go:generate mockery --name Client

package client

import (
	"github.com/aws/aws-sdk-go/service/iotsitewise/iotsitewiseiface"
)

type Client interface {
	iotsitewiseiface.IoTSiteWiseAPI
}
