package productmeta

import (
	"context"
	"github.com/arimanius/digivision-backend/internal/rank"
	v1 "github.com/arimanius/digivision-backend/pkg/api/v1"
)

type Fetcher interface {
	Fetch(ctx context.Context, product rank.Product) (*v1.Product, error)
	AsyncFetch(ctx context.Context, products []rank.Product, count int) chan *ProductWithError
}
