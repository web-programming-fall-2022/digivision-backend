package token

import (
	"context"
	"time"
)

type Manager interface {
	Generate(claims map[string]string, expiration time.Time) (string, error)
	Validate(ctx context.Context, tokenString string) (map[string]string, error)
	InvalidateToken(ctx context.Context, tokenString string) error
}
