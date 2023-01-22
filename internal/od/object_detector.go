package od

import "context"

type Position struct {
	X int
	Y int
}

type ObjectDetector interface {
	Detect(ctx context.Context, image []byte) (*Position, *Position, error)
}
