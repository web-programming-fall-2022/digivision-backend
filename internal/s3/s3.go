package s3

import (
	"context"
	"io"
)

type S3Client interface {
	Upload(ctx context.Context, bucket, path string, file io.Reader, size int64) error
}
