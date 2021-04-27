package propvals

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"math"
)

const (
	maxResponseSize = 250
	maxPagesToLoad  = 3 // TODO: how should this be optimized/selected?

	ResolutionRaw    = "RAW"
	ResolutionSecond = "1s"
	ResolutionMinute = "1m"
	ResolutionHour   = "1h"
	ResolutionDay    = "1d"
)

func roundUp(num float64) int64 {
	return int64(math.Ceil(num))
}

func durationForTimeRange(resolution string, timeRange backend.TimeRange) float64 {
	if ResolutionSecond == resolution {
		return timeRange.Duration().Seconds()
	} else if ResolutionMinute == resolution {
		return timeRange.Duration().Minutes()
	} else if ResolutionHour == resolution {
		return timeRange.Duration().Hours()
	} else {
		return timeRange.Duration().Hours() / 24
	}
}

// Takes the ceil of the pages. Ex: a duration that takes 2.1 pages to load all data takes 3 requests/pages to load
func pagesForResolution(resolution string, timeRange backend.TimeRange) int64 {
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
		pages := pagesForResolution(resolution, timeRange)
		dps := dataPointsForResolution(resolution, timeRange)
		// TODO: once '1s' resolution is supported, will need to add threshold for determining
		if dps <= maxDp && pages <= maxPagesToLoad {
			return resolution
		}
	}

	return ResolutionDay
}
