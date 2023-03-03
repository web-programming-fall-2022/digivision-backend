package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidEmail       = status.Error(codes.InvalidArgument, "Invalid email")
	TokenNotProvided   = status.Error(codes.Unauthenticated, "Token not provided")
	InvalidAccessToken = status.Error(codes.Unauthenticated, "Invalid access token")
	Internal           = status.Error(codes.Internal, "Internal error")
	NotLoggedIn        = status.Error(codes.Unauthenticated, "You are not logged in")
)
