package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Adapter 存储适配：预签名 GET、直传/直读、删除对象。
type Adapter interface {
	PresignGet(ctx context.Context, bucket, objectKey string, expire time.Duration) (string, error)
	PutObject(ctx context.Context, bucket, objectKey string, body io.Reader, size int64, contentType string) error
	GetObject(ctx context.Context, bucket, objectKey string) (body io.ReadCloser, size int64, contentType string, err error)
	Delete(ctx context.Context, bucket, path string) error
}

// NewAdapter 按当前活动存储配置构造适配器。
func NewAdapter() Adapter {
	cfg, _ := GetActiveConfig(context.Background())
	t := strings.ToLower(strings.TrimSpace(cfg.Type))
	switch t {
	case "oss", "s3", "minio":
		if a := newS3Adapter(cfg); a != nil {
			return a
		}
		fallthrough
	default:
		return &localAdapter{basePath: cfg.BasePath, publicBase: strings.TrimRight(cfg.PublicBaseURL, "/")}
	}
}

type localAdapter struct {
	basePath   string
	publicBase string
}

func (l *localAdapter) fullPath(objectKey string) string {
	key := strings.TrimPrefix(objectKey, "/")
	if l.basePath == "" {
		l.basePath = "./storage"
	}
	return filepath.Join(l.basePath, key)
}

func (l *localAdapter) PresignGet(ctx context.Context, bucket, objectKey string, expire time.Duration) (string, error) {
	if strings.TrimSpace(l.publicBase) != "" {
		p := strings.Trim(objectKey, "/")
		return strings.TrimRight(l.publicBase, "/") + "/" + p, nil
	}
	return "file://" + filepath.ToSlash(l.fullPath(objectKey)), nil
}

func (l *localAdapter) PutObject(ctx context.Context, bucket, objectKey string, body io.Reader, size int64, contentType string) error {
	p := l.fullPath(objectKey)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, body); err != nil {
		_ = os.Remove(p)
		return err
	}
	return nil
}

func (l *localAdapter) GetObject(ctx context.Context, bucket, objectKey string) (io.ReadCloser, int64, string, error) {
	p := l.fullPath(objectKey)
	st, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, "", err
		}
		return nil, 0, "", err
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, 0, "", err
	}
	return f, st.Size(), "", nil
}

func (l *localAdapter) Delete(ctx context.Context, bucket, path string) error {
	p := l.fullPath(path)
	if p == "" || p == "." || p == "/" {
		return nil
	}
	return os.Remove(p)
}

type s3Adapter struct {
	svc *s3.S3
	cfg Config
}

func newS3Adapter(cfg Config) *s3Adapter {
	region := strings.TrimSpace(cfg.Region)
	if region == "" {
		region = "us-east-1"
	}
	ep := strings.TrimSpace(cfg.Endpoint)
	awsCfg := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
	}
	if ep != "" {
		awsCfg.Endpoint = aws.String(ep)
		// false：virtual-hosted，预签名 URL 路径为 /key；true：path-style，路径为 /bucket/key。
		awsCfg.S3ForcePathStyle = aws.Bool(cfg.S3ForcePathStyle)
	}
	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil
	}
	return &s3Adapter{svc: s3.New(sess), cfg: cfg}
}

func (a *s3Adapter) PresignGet(ctx context.Context, bucket, objectKey string, expire time.Duration) (string, error) {
	if a == nil || a.svc == nil {
		return "", fmt.Errorf("s3 adapter not initialized")
	}
	b := bucket
	k := strings.TrimPrefix(objectKey, "/")
	// PresignSignatureVersion 已在 GetActiveConfig 中归一化为 v2 / v3 / v4。
	ver := a.cfg.PresignSignatureVersion
	if ver == "" {
		ver = "v4"
	}
	var raw string
	var err error
	switch ver {
	case "v2":
		raw, err = presignS3GetV2(a.cfg, b, k, expire, time.Now())
	default:
		// v3 与 v4 均走 AWS SDK SigV4 预签名。
		req, _ := a.svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(b),
			Key:    aws.String(k),
		})
		raw, err = req.Presign(expire)
	}
	if err != nil {
		return "", err
	}
	// public_base_url：替换预签名 URL 的协议与主机（含端口），不改变 Path/Query；与 SigV4 Host 绑定，换域可能导致 403。
	return applyPublicBaseURL(raw, a.cfg.PublicBaseURL)
}

func (a *s3Adapter) PutObject(ctx context.Context, bucket, objectKey string, body io.Reader, size int64, contentType string) error {
	if a == nil || a.svc == nil {
		return fmt.Errorf("s3 adapter not initialized")
	}
	var seeker io.ReadSeeker
	if rs, ok := body.(io.ReadSeeker); ok {
		seeker = rs
		if size <= 0 {
			n, err := rs.Seek(0, io.SeekEnd)
			if err != nil {
				return err
			}
			if _, err := rs.Seek(0, io.SeekStart); err != nil {
				return err
			}
			size = n
		}
	} else {
		b, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		seeker = bytes.NewReader(b)
		size = int64(len(b))
	}
	in := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(strings.TrimPrefix(objectKey, "/")),
		Body:   seeker,
	}
	if size > 0 {
		in.ContentLength = aws.Int64(size)
	}
	if strings.TrimSpace(contentType) != "" {
		in.ContentType = aws.String(contentType)
	}
	_, err := a.svc.PutObjectWithContext(ctx, in)
	return err
}

func (a *s3Adapter) GetObject(ctx context.Context, bucket, objectKey string) (io.ReadCloser, int64, string, error) {
	if a == nil || a.svc == nil {
		return nil, 0, "", fmt.Errorf("s3 adapter not initialized")
	}
	out, err := a.svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(strings.TrimPrefix(objectKey, "/")),
	})
	if err != nil {
		return nil, 0, "", err
	}
	var n int64
	if out.ContentLength != nil {
		n = *out.ContentLength
	}
	ct := ""
	if out.ContentType != nil {
		ct = *out.ContentType
	}
	return out.Body, n, ct, nil
}

func (a *s3Adapter) Delete(ctx context.Context, bucket, path string) error {
	if a == nil || a.svc == nil {
		return fmt.Errorf("s3 adapter not initialized")
	}
	_, err := a.svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	})
	return err
}
