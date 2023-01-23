package rank

import "github.com/arimanius/digivision-backend/internal/search"

type Product struct {
	Id       string
	Distance float32
}

type Ranker interface {
	Rank(productImages []search.ProductImage) []Product
}
