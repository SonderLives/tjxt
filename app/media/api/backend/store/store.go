// Package store 封装 S3 协议对象存储访问（MinIO / 腾讯云 COS / 阿里云 OSS 兼容网关）。
//
// 上传走预签名 PUT URL（客户端直传），播放/预览走预签名 GET URL（带时效防盗链），
// 持久访问地址为公共读桶的稳定 URL。换云厂商只需改 Endpoint/凭证，业务代码零改动。
package store

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config 对象存储连接配置。
type Config struct {
	Endpoint  string // 如 127.0.0.1:9000 或 cos.ap-guangzhou.myqcloud.com
	AccessKey string
	SecretKey string
	Bucket    string
	Secure    bool // 是否 HTTPS
}

// Store 对象存储客户端。
type Store struct {
	client *minio.Client
	bucket string
}

// New 建立客户端并确保桶存在（不存在则创建）。
func New(ctx context.Context, cfg Config) (*Store, error) {
	cli, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("create object store client: %w", err)
	}
	s := &Store{client: cli, bucket: cfg.Bucket}
	exists, err := cli.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket %s: %w", cfg.Bucket, err)
	}
	if !exists {
		if err := cli.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("create bucket %s: %w", cfg.Bucket, err)
		}
		// 开发环境桶策略：对象公共可读，持久 URL 可直接访问；写操作仍必须走预签名
		policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, cfg.Bucket)
		if err := cli.SetBucketPolicy(ctx, cfg.Bucket, policy); err != nil {
			return nil, fmt.Errorf("set bucket policy: %w", err)
		}
	}
	return s, nil
}

// PresignUpload 生成客户端直传用的预签名 PUT URL。
func (s *Store) PresignUpload(ctx context.Context, key string, expiry time.Duration) (string, error) {
	u, err := s.client.Presign(ctx, http.MethodPut, s.bucket, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("presign upload %s: %w", key, err)
	}
	return u.String(), nil
}

// PresignGet 生成带时效的播放/预览地址。
func (s *Store) PresignGet(ctx context.Context, key string, expiry time.Duration) (string, error) {
	u, err := s.client.Presign(ctx, http.MethodGet, s.bucket, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("presign get %s: %w", key, err)
	}
	return u.String(), nil
}

// PublicURL 返回对象在公共读桶中的持久访问地址。
func (s *Store) PublicURL(key string) string {
	return s.client.EndpointURL().String() + "/" + s.bucket + "/" + key
}
