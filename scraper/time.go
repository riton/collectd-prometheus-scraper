package scraper

import (
	"time"

	timemilli "github.com/Tigraine/go-timemilli"
)

func promTimestampToTime(pts *int64) time.Time {
	if pts == nil {
		return time.Time{}
	}

	return timemilli.FromUnixMilli(*pts)
}
