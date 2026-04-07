package util

import (
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// ToRFC3339UTC formats DB UTC datetime as RFC3339 without timezone shift.
// It keeps wall-clock value and appends "Z" (e.g. "2026-03-23T08:00:00Z").
// Returns "" if t is nil.
func ToRFC3339UTC(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	s := t.Format("Y-m-d H:i:s")
	if s == "" {
		return ""
	}
	return s[:10] + "T" + s[11:] + "Z"
}

// ToRFC3339UTCPtr is like ToRFC3339UTC but returns nil when t is nil, otherwise a pointer to the formatted string.
func ToRFC3339UTCPtr(t *gtime.Time) *string {
	if t == nil {
		return nil
	}
	s := ToRFC3339UTC(t)
	return &s
}

// ToRFC3339UTCShift converts time instant to UTC timezone and formats as RFC3339.
// Use this for runtime "now" values that represent local timezone instant.
func ToRFC3339UTCShift(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Time.Format(time.RFC3339)
}
