package storage

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// presignS3GetV2 生成 S3 Signature Version 2 的 GET 预签名 URL（query string auth）。
// 与 newS3Adapter 一致：配置了 endpoint 时使用 path-style；否则使用 virtual-hosted 形式。
func presignS3GetV2(cfg Config, bucket, objectKey string, expire time.Duration, now time.Time) (string, error) {
	bucket = strings.TrimSpace(bucket)
	key := strings.TrimPrefix(strings.TrimSpace(objectKey), "/")
	if bucket == "" || key == "" {
		return "", fmt.Errorf("presign v2: empty bucket or key")
	}
	if expire <= 0 {
		return "", fmt.Errorf("presign v2: expire must be positive")
	}
	ak := strings.TrimSpace(cfg.AccessKey)
	sk := strings.TrimSpace(cfg.SecretKey)
	if ak == "" || sk == "" {
		return "", fmt.Errorf("presign v2: missing access or secret key")
	}

	expires := now.Add(expire).Unix()
	if expires <= now.Unix() {
		return "", fmt.Errorf("presign v2: invalid expires")
	}

	region := strings.TrimSpace(cfg.Region)
	if region == "" {
		region = "us-east-1"
	}
	ep := strings.TrimSpace(cfg.Endpoint)

	var reqURL string
	var canonicalResource string

	keyPath := encodeS3PathSegments(key)

	if ep != "" {
		epu, err := url.Parse(ep)
		if err != nil {
			return "", fmt.Errorf("presign v2: endpoint: %w", err)
		}
		if cfg.S3ForcePathStyle {
			base := strings.TrimRight(ep, "/")
			reqPath := "/" + bucket + "/" + keyPath
			reqURL = base + reqPath
			canonicalResource = reqPath
		} else {
			// virtual-hosted：scheme://bucket.host/keyPath，路径不含 /bucket/ 前缀。
			if epu.Scheme == "" || epu.Host == "" {
				return "", fmt.Errorf("presign v2: endpoint must include scheme and host")
			}
			virtualHost := bucket + "." + epu.Host
			reqURL = epu.Scheme + "://" + virtualHost + "/" + keyPath
			canonicalResource = "/" + bucket + "/" + keyPath
		}
	} else {
		host := s3VirtualHostedHost(bucket, region)
		reqPath := "/" + keyPath
		reqURL = "https://" + host + reqPath
		canonicalResource = "/" + bucket + "/" + keyPath
	}

	stringToSign := "GET\n\n\n" + strconv.FormatInt(expires, 10) + "\n" + canonicalResource

	mac := hmac.New(sha1.New, []byte(sk))
	_, _ = mac.Write([]byte(stringToSign))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	u, err := url.Parse(reqURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("AWSAccessKeyId", ak)
	q.Set("Expires", strconv.FormatInt(expires, 10))
	q.Set("Signature", sig)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func s3VirtualHostedHost(bucket, region string) string {
	if region == "us-east-1" {
		return bucket + ".s3.amazonaws.com"
	}
	return bucket + ".s3." + region + ".amazonaws.com"
}

// encodeS3PathSegments 对 object key 按 “/” 分段做 PathEscape，保留层级。
func encodeS3PathSegments(key string) string {
	if key == "" {
		return ""
	}
	parts := strings.Split(key, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}
