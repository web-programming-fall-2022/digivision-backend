package rank

import "github.com/web-programming-fall-2022/digivision-backend/internal/search"

type Product struct {
	Id    string
	Score float32
}

type Ranker interface {
	Rank(productImages []search.ProductImage) []Product
}
