package rank

import "github.com/web-programming-fall-2022/digivision-backend/internal/search"

// FirstImageRanker implements Ranker interface{}
type FirstImageRanker struct {
}

// NewFirstImageRanker returns a new FirstImageRanker
func NewFirstImageRanker() *FirstImageRanker {
	return &FirstImageRanker{}
}

// Rank implements Ranker interface{}
func (r *FirstImageRanker) Rank(productImages []search.ProductImage) []Product {
	productRankById := make(map[string]int)
	result := make([]Product, 0)
	for _, productImage := range productImages {
		if _, ok := productRankById[productImage.ProductId]; !ok {
			productRankById[productImage.ProductId] = len(result)
			result = append(result, Product{
				Id:    productImage.ProductId,
				Score: -productImage.Distance,
			})
		}
	}
	return result
}
