package storage

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestRewriteURLOrigin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		raw     string
		base    string
		want    string
		wantErr bool
	}{
		{
			name: "replace host keeps path and query",
			raw:  "https://internal:9000/mybucket/obj?AWSAccessKeyId=x&Expires=1&Signature=abc%2B",
			base: "https://cdn.example.com",
			want: "https://cdn.example.com/mybucket/obj?AWSAccessKeyId=x&Expires=1&Signature=abc%2B",
		},
		{
			name: "with port",
			raw:  "http://minio.local:9000/a/b?k=v",
			base: "https://files.example.org:8443",
			want: "https://files.example.org:8443/a/b?k=v",
		},
		{
			name:    "invalid public base",
			raw:     "https://a/b",
			base:    "not-a-url",
			wantErr: true,
		},
		{
			name: "empty public base returns raw",
			raw:  "https://x/y",
			base: "",
			want: "https://x/y",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := rewriteURLOrigin(tc.raw, tc.base)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestApplyPublicBaseURL(t *testing.T) {
	t.Parallel()
	raw := "http://127.0.0.1:9000/bucket/key?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	got, err := applyPublicBaseURL(raw, "https://cdn.example.com/")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://cdn.example.com/bucket/key?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNormalizePresignSigVersion(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	if g := normalizePresignSigVersion(ctx, ""); g != "v4" {
		t.Fatalf("empty -> v4, got %q", g)
	}
	if g := normalizePresignSigVersion(ctx, "SIGV4"); g != "v4" {
		t.Fatalf("sigv4 -> v4, got %q", g)
	}
	if g := normalizePresignSigVersion(ctx, "3"); g != "v3" {
		t.Fatalf("3 -> v3, got %q", g)
	}
	if g := normalizePresignSigVersion(ctx, "v2"); g != "v2" {
		t.Fatalf("v2 -> v2, got %q", g)
	}
	if g := normalizePresignSigVersion(ctx, "bogus"); g != "v4" {
		t.Fatalf("bogus -> v4, got %q", g)
	}
}

func TestPresignS3GetV2_pathStyle(t *testing.T) {
	t.Parallel()
	cfg := Config{
		Endpoint:         "http://127.0.0.1:9000",
		AccessKey:        "MYACCESSKEY",
		SecretKey:        "MYSECRETKEY",
		Region:           "us-east-1",
		S3ForcePathStyle: true,
	}
	now := time.Unix(1700000000, 0)
	expire := time.Hour
	u, err := presignS3GetV2(cfg, "mybucket", "folder/object.txt", expire, now)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(u, "http://127.0.0.1:9000/mybucket/folder/object.txt?") {
		t.Fatalf("unexpected url prefix: %s", u)
	}
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}
	q := parsed.Query()
	if q.Get("AWSAccessKeyId") != "MYACCESSKEY" {
		t.Fatalf("AWSAccessKeyId: %q", q.Get("AWSAccessKeyId"))
	}
	if q.Get("Expires") != "1700003600" {
		t.Fatalf("Expires: %q", q.Get("Expires"))
	}
	sig := q.Get("Signature")
	if sig == "" {
		t.Fatal("empty Signature")
	}
	// 固定输入下签名稳定
	u2, _ := presignS3GetV2(cfg, "mybucket", "folder/object.txt", expire, now)
	if u2 != u {
		t.Fatalf("presign not stable: %s vs %s", u, u2)
	}
}

func TestPresignS3GetV2_customEndpointVirtualHosted(t *testing.T) {
	t.Parallel()
	cfg := Config{
		Endpoint:         "https://minio.example.com:9000",
		AccessKey:        "AK",
		SecretKey:        "SK",
		Region:           "us-east-1",
		S3ForcePathStyle: false,
	}
	u, err := presignS3GetV2(cfg, "mybucket", "a/b.ts", time.Minute, time.Unix(1700000000, 0))
	if err != nil {
		t.Fatal(err)
	}
	wantPrefix := "https://mybucket.minio.example.com:9000/a/b.ts?"
	if !strings.HasPrefix(u, wantPrefix) {
		t.Fatalf("want prefix %q got %q", wantPrefix, u)
	}
}

func TestPresignS3GetV2_virtualHosted(t *testing.T) {
	t.Parallel()
	cfg := Config{
		Endpoint:  "",
		AccessKey: "AK",
		SecretKey: "SK",
		Region:    "eu-west-1",
	}
	now := time.Unix(1600000000, 0)
	u, err := presignS3GetV2(cfg, "mybucket", "k", time.Minute, now)
	if err != nil {
		t.Fatal(err)
	}
	wantPrefix := "https://mybucket.s3.eu-west-1.amazonaws.com/k?"
	if !strings.HasPrefix(u, wantPrefix) {
		t.Fatalf("want prefix %q got %q", wantPrefix, u)
	}
}

func TestPresignS3GetV2_usEast1Host(t *testing.T) {
	t.Parallel()
	cfg := Config{AccessKey: "A", SecretKey: "S", Region: "us-east-1"}
	u, err := presignS3GetV2(cfg, "b", "k", time.Minute, time.Unix(1, 0))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(u, "https://b.s3.amazonaws.com/k?") {
		t.Fatalf("unexpected: %s", u)
	}
}

func TestEncodeS3PathSegments(t *testing.T) {
	t.Parallel()
	if encodeS3PathSegments("a/b c") != "a/b%20c" {
		t.Fatalf("got %q", encodeS3PathSegments("a/b c"))
	}
}
