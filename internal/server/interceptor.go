package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/web-programming-fall-2022/digivision-backend/internal/errors"
	"github.com/web-programming-fall-2022/digivision-backend/internal/storage"
	"github.com/web-programming-fall-2022/digivision-backend/internal/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
)

type AuthInterceptor struct {
	storage      *storage.Storage
	tokenManager token.Manager
}

func NewAuthInterceptor(storage *storage.Storage, tokenManager token.Manager) *AuthInterceptor {
	return &AuthInterceptor{
		storage:      storage,
		tokenManager: tokenManager,
	}
}

type authUserKey struct{}

func (i *AuthInterceptor) InterceptServer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		accessToken := getMetadataValue(ctx, "x-access-token")
		if len(accessToken) != 0 {
			logrus.Debug("Using access token: ", accessToken)
			claims, err := i.tokenManager.Validate(ctx, accessToken)
			if err != nil {
				return nil, errors.InvalidAccessToken
			}
			userId, ok := claims["userID"]
			if !ok {
				return nil, errors.InvalidAccessToken
			}
			id, err := strconv.Atoi(userId)
			if err != nil {
				return nil, errors.InvalidAccessToken
			}
			user, err := i.storage.GetUserByID(uint(id))

			ctx = AttachUserToCtx(ctx, user)
		}

		return handler(ctx, req)
	}
}

func AttachUserToCtx(ctx context.Context, user *storage.UserAccount) context.Context {
	return context.WithValue(ctx, authUserKey{}, *user)
}

func GetContextUser(ctx context.Context) *storage.UserAccount {
	val := ctx.Value(authUserKey{})
	if val == nil {
		return nil
	}
	user, ok := val.(storage.UserAccount)
	if !ok {
		return nil
	}
	return &user
}

func getMetadataValue(ctx context.Context, key string) string {
	md, _ := metadata.FromIncomingContext(ctx)
	if md == nil {
		return ""
	}
	values := metadata.ValueFromIncomingContext(ctx, key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
