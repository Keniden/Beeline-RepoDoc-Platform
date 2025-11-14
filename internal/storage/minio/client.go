package minio

import (
    "bytes"
    "context"
    "errors"
    "fmt"

    "github.com/beeline/repodoc/configs"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
    Client *minio.Client
    Bucket string
}

func New(ctx context.Context, cfg configs.StorageConfig) (*Client, error) {
    mc, err := minio.New(cfg.MinioEndpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
        Secure: false,
    })
    if err != nil {
        return nil, fmt.Errorf("minio new: %w", err)
    }
    if err := mc.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{}); err != nil {
        var resp minio.ErrorResponse
        if errors.As(err, &resp) && resp.Code == "BucketAlreadyOwnedByYou" {
        } else {
            return nil, fmt.Errorf("create bucket: %w", err)
        }
    }
    return &Client{Client: mc, Bucket: cfg.MinioBucket}, nil
}

func (c *Client) Upload(ctx context.Context, object string, data []byte) error {
    _, err := c.Client.PutObject(ctx, c.Bucket, object, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
    return err
}
