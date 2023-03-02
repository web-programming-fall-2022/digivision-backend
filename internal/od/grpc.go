package od

import (
	"context"
	"github.com/web-programming-fall-2022/digivision-backend/internal/api/od"
)

// GrpcObjectDetector implements ObjectDetector interface{}
type GrpcObjectDetector struct {
	stub od.ObjectDetectorClient
}

// NewGrpcObjectDetector returns a new GrpcObjectDetector
func NewGrpcObjectDetector(stub od.ObjectDetectorClient) *GrpcObjectDetector {
	return &GrpcObjectDetector{stub: stub}
}

// Detect implements ObjectDetector interface{}
func (v *GrpcObjectDetector) Detect(ctx context.Context, image []byte) (*Position, *Position, error) {
	res, err := v.stub.Detect(ctx, &od.Image{Image: image})
	if err != nil {
		return nil, nil, err
	}
	return &Position{
			X: int(res.TopLeft.X),
			Y: int(res.TopLeft.Y),
		}, &Position{
			X: int(res.BottomRight.X),
			Y: int(res.BottomRight.Y),
		}, nil
}
