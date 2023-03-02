package productmeta

import (
	"context"
	"github.com/web-programming-fall-2022/digivision-backend/internal/rank"
	v1 "github.com/web-programming-fall-2022/digivision-backend/pkg/api/v1"
)

type Fetcher interface {
	Fetch(ctx context.Context, product rank.Product) (*v1.Product, error)
	AsyncFetch(ctx context.Context, products []rank.Product, count int) chan *ProductWithError
}
