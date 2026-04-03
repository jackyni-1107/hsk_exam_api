package util

import (
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// ToRFC3339UTC formats t in UTC as RFC3339 (e.g. 2026-03-23T08:00:00Z). Returns "" if t is nil.
func ToRFC3339UTC(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Time.Format(time.RFC3339)
}

// ToRFC3339UTCPtr is like ToRFC3339UTC but returns nil when t is nil, otherwise a pointer to the formatted string.
func ToRFC3339UTCPtr(t *gtime.Time) *string {
	if t == nil {
		return nil
	}
	s := ToRFC3339UTC(t)
	return &s
}
