package middleware

import "strings"

type permissionRule struct {
	Method     string
	Pattern    []string
	Permission string
}

var adminPermissionRules = []permissionRule{
	{Method: "POST", Pattern: []string{"exam", "paper", "import"}, Permission: "exam:import"},
	{Method: "GET", Pattern: []string{"exam", "attempt", "*"}, Permission: "exam:result:list"},
	{Method: "PUT", Pattern: []string{"exam", "attempt", "*", "subjective-scores"}, Permission: "exam:result:grade"},
	{Method: "POST", Pattern: []string{"file", "upload"}, Permission: "file:list"},
	{Method: "POST", Pattern: []string{"task", "run"}, Permission: "task:run"},
	{Method: "GET", Pattern: []string{"task", "log"}, Permission: "task:log"},
	{Method: "POST", Pattern: []string{"user", "*", "kick-sessions"}, Permission: "user:update"},
}

func matchPermissionRule(method string, route adminRoute) (string, bool) {
	for _, rule := range adminPermissionRules {
		if rule.matches(method, route) {
			return rule.Permission, true
		}
	}
	return "", false
}

func (r permissionRule) matches(method string, route adminRoute) bool {
	if r.Method != "" && !strings.EqualFold(r.Method, method) {
		return false
	}
	return route.matches(r.Pattern)
}
