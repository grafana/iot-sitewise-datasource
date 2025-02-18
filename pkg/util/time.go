package util

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TimeRangeToUnix(tr backend.TimeRange) (from *time.Time, to *time.Time) {
	from = aws.Time(time.Unix(tr.From.Unix(), 0))
	to = aws.Time(time.Unix(tr.To.Unix(), 0))
	return
}
