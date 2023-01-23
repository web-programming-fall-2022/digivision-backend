package rank

import (
	"github.com/arimanius/digivision-backend/internal/search"
	"math"
	"sort"
)

type DistCountRanker struct {
}

func NewDistCountRanker() *DistCountRanker {
	return &DistCountRanker{}
}

func (r *DistCountRanker) Rank(productImages []search.ProductImage) []Product {
	productScore := make(map[string]float64)
	for i, productImage := range productImages {
		// make decay factor float64
		decayFactor := 1.0 / (1.0 + float64(i))
		score := math.Exp(float64(-0.002*productImage.Distance)) * decayFactor
		p, ok := productScore[productImage.ProductId]
		if !ok {
			productScore[productImage.ProductId] = score
		} else {
			productScore[productImage.ProductId] = p + score
		}
	}

	result := make([]Product, 0)
	for id, score := range productScore {
		result = append(result, Product{
			Id:    id,
			Score: float32(score),
		})
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	return result
}
