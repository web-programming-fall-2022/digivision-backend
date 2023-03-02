package token

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"github.com/web-programming-fall-2022/digivision-backend/internal/storage"
	"strconv"
	"time"
)

type JWTManager struct {
	secret  []byte
	Storage *storage.Storage
	RDB     *redis.Client
}

func NewJWTManager(secret string, store *storage.Storage, rdb *redis.Client) *JWTManager {
	return &JWTManager{
		secret:  []byte(secret),
		Storage: store,
		RDB:     rdb,
	}
}

func (m *JWTManager) Generate(claims map[string]string, expiration time.Time) (string, error) {
	mapClaims := jwt.MapClaims{
		"exp": expiration.Unix(),
	}
	for key, value := range claims {
		mapClaims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) Validate(ctx context.Context, tokenString string) (map[string]string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	err = m.CheckUnauthorizedToken(ctx, tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	result := make(map[string]string)
	for key, value := range claims {
		if val, ok := value.(string); ok {
			result[key] = val
		}
	}
	return result, nil
}

func (m *JWTManager) InvalidateToken(ctx context.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token")
	}
	userId, _ := strconv.Atoi(claims["userID"].(string))
	var expiration time.Time
	switch exp := claims["exp"].(type) {
	case float64:
		expiration = time.Unix(int64(exp), 0)
	case json.Number:
		v, _ := exp.Int64()
		expiration = time.Unix(v, 0)
	}
	err = m.Storage.CreateUnauthorizedToken(&storage.UnauthorizedToken{
		UserID:     uint(userId),
		Token:      tokenString,
		Expiration: expiration,
	})
	m.RDB.SetEx(ctx, tokenString, "false", time.Until(expiration))
	if err != nil {
		return err
	}
	return nil
}

func (m *JWTManager) CheckUnauthorizedToken(ctx context.Context, tokenString string) error {
	resp := m.RDB.Get(ctx, tokenString)
	if resp.Err() == redis.Nil {
		_, err := m.Storage.GetUnauthorizedToken(tokenString)
		if err != nil {
			m.RDB.SetEx(ctx, tokenString, "true", time.Minute*10)
			return nil
		}
		m.RDB.SetEx(ctx, tokenString, "false", time.Minute*10)
	} else if resp.Err() != nil {
		_, err := m.Storage.GetUnauthorizedToken(tokenString)
		if err != nil {
			m.RDB.SetEx(ctx, tokenString, "true", time.Minute*10)
			return nil
		}
		m.RDB.SetEx(ctx, tokenString, "false", time.Minute*10)
	} else if resp.Val() == "true" {
		return nil
	}
	return errors.New("token is unauthorized")
}
