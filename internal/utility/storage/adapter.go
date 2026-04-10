package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Adapter 存储适配：预签名 GET 与删除对象。
type Adapter interface {
	PresignGet(ctx context.Context, bucket, objectKey string, expire time.Duration) (string, error)
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

func (l *localAdapter) Delete(ctx context.Context, bucket, path string) error {
	p := l.fullPath(path)
	if p == "" || p == "." || p == "/" {
		return nil
	}
	return os.Remove(p)
}

type s3Adapter struct {
	svc *s3.S3
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
		awsCfg.S3ForcePathStyle = aws.Bool(true)
	}
	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil
	}
	return &s3Adapter{svc: s3.New(sess)}
}

func (a *s3Adapter) PresignGet(ctx context.Context, bucket, objectKey string, expire time.Duration) (string, error) {
	if a == nil || a.svc == nil {
		return "", fmt.Errorf("s3 adapter not initialized")
	}
	b := bucket
	k := strings.TrimPrefix(objectKey, "/")
	req, _ := a.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(b),
		Key:    aws.String(k),
	})
	return req.Presign(expire)
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
