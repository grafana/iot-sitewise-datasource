package testdata

import (
	"time"
)

var (
	Now         = time.Now()
	TwoHours    = Now.Add(-2 * time.Hour)
	FiveMinutes = Now.Add(-5 * time.Minute)
	OneDay      = Now.Add(-24 * time.Hour)
	FifteenDays = Now.Add(-15 * 24 * time.Hour)
	OneMonth    = Now.Add(-24 * 31 * time.Hour)
)
