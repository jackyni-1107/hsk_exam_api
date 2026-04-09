package config

import "strings"

// JoinHTTPPath 拼接 URL 路径：prefix 可为空（表示由反向代理已剥掉 /api）；relative 为去掉首尾 / 的一段或多段。
func JoinHTTPPath(prefix, relative string) string {
	prefix = strings.TrimSpace(prefix)
	prefix = strings.TrimSuffix(prefix, "/")
	relative = strings.Trim(strings.TrimSpace(relative), "/")
	if relative == "" {
		if prefix == "" {
			return "/"
		}
		return prefix
	}
	if prefix == "" {
		return "/" + relative
	}
	return prefix + "/" + relative
}
