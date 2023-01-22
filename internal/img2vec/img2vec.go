package img2vec

import "context"

type Img2Vec interface {
	Vectorize(ctx context.Context, image []byte) ([]float32, error)
}
