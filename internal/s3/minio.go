package s3

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type MinioClient struct {
	minioClient *minio.Client
}

func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*MinioClient, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioClient{minioClient: minioClient}, nil
}

func (m *MinioClient) Upload(ctx context.Context, bucket, path string, file io.Reader, size int64) error {
	exists, err := m.minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		err = m.minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	_, err = m.minioClient.PutObject(ctx, bucket, path, file, size, minio.PutObjectOptions{})
	return err
}
