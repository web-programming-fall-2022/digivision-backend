package productmeta

import (
	"context"
	v1 "github.com/arimanius/digivision-backend/pkg/api/v1"
)

type Fetcher interface {
	Fetch(ctx context.Context, productId string) (*v1.Product, error)
}
