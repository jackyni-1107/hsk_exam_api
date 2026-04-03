package utility

import (
	"strings"
)

// ParseDeviceInfo 从 User-Agent 提取简短可读描述（登录审计用）。
func ParseDeviceInfo(userAgent string) string {
	ua := strings.TrimSpace(userAgent)
	if ua == "" {
		return ""
	}
	if len(ua) > 512 {
		ua = ua[:512]
	}
	browser := "Unknown"
	osName := "Unknown"
	lower := strings.ToLower(ua)
	switch {
	case strings.Contains(lower, "edg/"):
		browser = "Edge"
	case strings.Contains(lower, "chrome/") && !strings.Contains(lower, "chromium"):
		browser = "Chrome"
	case strings.Contains(lower, "firefox/"):
		browser = "Firefox"
	case strings.Contains(lower, "safari/") && !strings.Contains(lower, "chrome"):
		browser = "Safari"
	case strings.Contains(lower, "micromessenger"):
		browser = "WeChat"
	}
	switch {
	case strings.Contains(lower, "windows"):
		osName = "Windows"
	case strings.Contains(lower, "mac os") || strings.Contains(lower, "macintosh"):
		osName = "macOS"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad"):
		osName = "iOS"
	case strings.Contains(lower, "android"):
		osName = "Android"
	case strings.Contains(lower, "linux"):
		osName = "Linux"
	}
	return browser + " / " + osName
}
