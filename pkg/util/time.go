package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TimeRangeToUnix(tr backend.TimeRange) (from *time.Time, to *time.Time) {
	from = aws.Time(time.Unix(tr.From.Unix(), 0))
	to = aws.Time(time.Unix(tr.To.Unix(), 0))
	return
}

// TimestampToDate replaces epoch timestamps in rawSQL with formatted TIMESTAMP strings
// based on the provided timeRange.
func TimestampToDate(rawSQL string, timeRange *backend.TimeRange) string {
	if timeRange != nil {
		fromEpoch, toEpoch := GetEpochRange(timeRange.From, timeRange.To)
		fromStr, toStr := GetFormattedTimeRange(timeRange.From, timeRange.To)
		rawSQL = ReplaceEpochWithTimestamp(rawSQL, fromEpoch, fromStr)
		rawSQL = ReplaceEpochWithTimestamp(rawSQL, toEpoch, toStr)
	}
	return rawSQL
}

// GetEpochRange returns the Unix epoch timestamps (in seconds)
// for the provided 'from' and 'to' time.Time values.
func GetEpochRange(from, to time.Time) (int64, int64) {
	return from.Unix(), to.Unix()
}

// GetFormattedTimeRange returns the 'from' and 'to' time.Time values
// formatted as UTC strings in the "TIMESTAMP 'YYYY-MM-DD HH:MM:SS'" format.
func GetFormattedTimeRange(from, to time.Time) (string, string) {
	const sqlFormat = "2006-01-02 15:04:05"
	return from.UTC().Format(sqlFormat), to.UTC().Format(sqlFormat)
}

// ReplaceEpochWithTimestamp replaces all occurrences of the given epoch
// timestamp (in seconds) in the provided SQL string with a formatted
// SQL TIMESTAMP literal.
func ReplaceEpochWithTimestamp(sql string, epoch int64, timestampStr string) string {
	epochStr := fmt.Sprintf("%d", epoch)
	timestampSQL := fmt.Sprintf("TIMESTAMP '%s'", timestampStr)
	return strings.ReplaceAll(sql, epochStr, timestampSQL)
}
