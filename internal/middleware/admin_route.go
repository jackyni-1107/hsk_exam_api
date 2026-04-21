package middleware

import "strings"

type adminRoute struct {
	segments []string
}

func parseAdminRoute(path string) (adminRoute, bool) {
	path = strings.Trim(path, "/")
	if path == "" {
		return adminRoute{}, false
	}
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] == "api" {
		parts = parts[1:]
	}
	if len(parts) < 2 || parts[0] != "admin" {
		return adminRoute{}, false
	}
	segments := parts[1:]
	if len(segments) == 0 {
		return adminRoute{}, false
	}
	return adminRoute{segments: segments}, true
}

func (r adminRoute) module() string {
	if len(r.segments) == 0 {
		return ""
	}
	return r.segments[0]
}

func (r adminRoute) matches(pattern []string) bool {
	if len(r.segments) != len(pattern) {
		return false
	}
	for i, expect := range pattern {
		if expect == "*" {
			continue
		}
		if r.segments[i] != expect {
			return false
		}
	}
	return true
}
