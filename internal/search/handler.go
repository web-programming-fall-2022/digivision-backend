package search

import "context"

type ProductImage struct {
	ProductId string
	ImageId   string
	Distance  float32
}

type Handler interface {
	Search(ctx context.Context, query []float32, topK int) ([]ProductImage, error)
}
