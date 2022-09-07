package propvals

import (
	"math"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

const (
	maxHistoryResponseSize = 250
	maxHistoryPagesToLoad  = 4 // 100-200ms * 4 = 800ms max on average ?

	maxInterpolatedResponseSize = 10
	maxInterpolatedPagesToLoad  = 10 // 100-200ms * 10 = 2s max on average ?

	ResolutionRaw        = "RAW"
	ResolutionSecond     = "1s"
	ResolutionTenSeconds = "10s"
	ResolutionMinute     = "1m"
	ResolutionTenMinutes = "10m"
	ResolutionHour       = "1h"
	ResolutionTenHours   = "10h"
	ResolutionDay        = "1d"
)

func roundUp(num float64) int64 {
	return int64(math.Ceil(num))
}

func durationForTimeRange(resolution string, timeRange backend.TimeRange) float64 {
	if ResolutionSecond == resolution {
		return timeRange.Duration().Seconds()
	} else if ResolutionTenSeconds == resolution {
		return timeRange.Duration().Seconds() / 10
	} else if ResolutionMinute == resolution {
		return timeRange.Duration().Minutes()
	} else if ResolutionTenMinutes == resolution {
		return timeRange.Duration().Minutes() / 10
	} else if ResolutionHour == resolution {
		return timeRange.Duration().Hours()
	} else if ResolutionTenHours == resolution {
		return timeRange.Duration().Hours() / 10
	} else {
		return timeRange.Duration().Hours() / 24
	}
}

// Takes the ceil of the pages. Ex: a duration that takes 2.1 pages to load all data takes 3 requests/pages to load
func pagesForResolution(resolution string, timeRange backend.TimeRange, maxResponseSize float64) int64 {
	duration := durationForTimeRange(resolution, timeRange)
	return roundUp(duration / maxResponseSize)
}

// Takes the floor of the duration - ex: duration of 10.5 minutes would load 10 data points
func dataPointsForResolution(resolution string, timeRange backend.TimeRange) int64 {
	return int64(durationForTimeRange(resolution, timeRange))
}

func Resolution(query models.BaseQuery) string {

	timeRange := query.TimeRange
	maxDp := query.MaxDataPoints

	for _, resolution := range []string{ResolutionSecond, ResolutionMinute, ResolutionHour} {
		pages := pagesForResolution(resolution, timeRange, maxHistoryResponseSize)
		dps := dataPointsForResolution(resolution, timeRange)
		// TODO: once '1s' resolution is supported, will need to add threshold for determining
		if dps <= maxDp && pages <= maxHistoryPagesToLoad {
			return resolution
		}
	}

	return ResolutionDay
}

func InterpolatedResolution(query models.AssetPropertyValueQuery) string {
	timeRange := query.TimeRange
	maxDp := query.MaxDataPoints

	for _, resolution := range []string{ResolutionSecond, ResolutionTenSeconds, ResolutionMinute, ResolutionTenMinutes, ResolutionHour, ResolutionTenHours} {
		pages := pagesForResolution(resolution, timeRange, maxInterpolatedResponseSize)
		dps := dataPointsForResolution(resolution, timeRange)
		if dps <= maxDp && pages <= maxInterpolatedPagesToLoad {
			return resolution
		}
	}

	return ResolutionDay
}

func ResolutionToDuration(resolution string) time.Duration {
	switch resolution {
	case ResolutionSecond:
		return time.Second
	case ResolutionTenSeconds:
		return 10 * time.Second
	case ResolutionMinute:
		return time.Minute
	case ResolutionTenMinutes:
		return 10 * time.Minute
	case ResolutionHour:
		return time.Hour
	case ResolutionTenHours:
		return 10 * time.Hour
	case ResolutionDay:
		fallthrough
	default:
		return 24 * time.Hour
	}
}
