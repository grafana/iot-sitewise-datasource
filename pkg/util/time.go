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

// GetFormattedTimeRange returns the time.Time values
// formatted as UTC strings in the "TIMESTAMP 'YYYY-MM-DD HH:MM:SS'" format.
func GetFormattedTimeRange(time time.Time) string {
	const sqlFormat = "2006-01-02 15:04:05"
	return time.UTC().Format(sqlFormat)
}
