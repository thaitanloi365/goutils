package goutils

import (
	"time"

	"github.com/labstack/gommon/log"
)

const fallbackTimezone = "America/Chicago"

// TimeIn tz
func TimeIn(t time.Time, tz string, fallbackTZ ...string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err == nil {
		t = t.In(loc)

		return t
	}
	log.Errorf("Parse timezone=%s error: %v", tz, err)

	var fallback = "America/Chicago"
	if len(fallbackTZ) > 0 {
		fallback = fallbackTZ[0]
	}

	loc, err = time.LoadLocation(fallback)
	if err == nil {
		t = t.In(loc)
		return t
	}

	log.Errorf("Parse timezone=%s error: %v", fallback, err)

	return t
}

// NowIn now
func NowIn(tz string) time.Time {
	var t = TimeIn(time.Now(), tz)
	return t
}
