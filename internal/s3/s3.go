package s3

import (
	"context"
	"io"
)

type Client interface {
	Upload(ctx context.Context, bucket, path string, file io.Reader, size int64) error
	Download(ctx context.Context, bucket, path string) (io.ReadCloser, error)
}
