package storage

import (
	"fmt"
	"net/url"
	"strings"
)

// rewriteURLOrigin 将预签名 URL 的协议、主机名（含端口）替换为 publicBase 中的对应部分，
// 保留 Path 与 RawQuery。SigV4 签名与请求 Host 绑定，替换域名后若与服务端校验不一致可能 403。
func rewriteURLOrigin(rawURL, publicBase string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	publicBase = strings.TrimSpace(publicBase)
	if publicBase == "" {
		return rawURL, nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	b, err := url.Parse(publicBase)
	if err != nil {
		return "", err
	}
	if b.Scheme == "" || b.Host == "" {
		return "", fmt.Errorf("public_base_url must include scheme and host")
	}
	u.Scheme = b.Scheme
	u.Host = b.Host
	return u.String(), nil
}
