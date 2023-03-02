package img2vec

import (
	"context"
	"github.com/web-programming-fall-2022/digivision-backend/internal/api/img2vec"
)

// GrpcImg2Vec implements Img2Vec interface{}
type GrpcImg2Vec struct {
	stub img2vec.Img2VecClient
}

// NewGrpcImg2Vec returns a new GrpcImg2Vec
func NewGrpcImg2Vec(stub img2vec.Img2VecClient) *GrpcImg2Vec {
	return &GrpcImg2Vec{stub: stub}
}

// Vectorize implements Img2Vec interface{}
func (v *GrpcImg2Vec) Vectorize(ctx context.Context, image []byte) ([]float32, error) {
	res, err := v.stub.Vectorize(ctx, &img2vec.Image{Image: image})
	if err != nil {
		return nil, err
	}
	return res.Vector, nil
}
