package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/web-programming-fall-2022/digivision-backend/internal/errors"
	"github.com/web-programming-fall-2022/digivision-backend/internal/storage"
	"github.com/web-programming-fall-2022/digivision-backend/internal/token"
	pb "github.com/web-programming-fall-2022/digivision-backend/pkg/api/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"time"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	AuthTokenExpire    int64
	RefreshTokenExpire int64
	TokenManager       token.Manager
	Storage            *storage.Storage
}

func NewAuthServiceServer(
	tokenManager token.Manager,
	storage *storage.Storage,
	authTokenExpire int64,
	refreshTokenExpire int64,
) *AuthServiceServer {
	return &AuthServiceServer{
		TokenManager:       tokenManager,
		Storage:            storage,
		AuthTokenExpire:    authTokenExpire,
		RefreshTokenExpire: refreshTokenExpire,
	}
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.Storage.GetUserByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	authToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not generate tokens")
	}
	return &pb.LoginResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = s.Storage.GetUserByEmail(req.Email)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}
	_, err = s.Storage.GetUserByPhoneNumber(req.PhoneNumber)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "phone number already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not hash password")
	}
	user := storage.UserAccount{
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Gender:       req.Gender,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: string(hash),
	}
	err = s.Storage.CreateUser(&user)
	if err != nil {
		logrus.Errorln(err)
		return nil, status.Error(codes.Internal, "could not create user")
	}
	authToken, refreshToken, err := s.generateTokens(&user)
	if err != nil {
		logrus.Errorln(err)
		return nil, status.Error(codes.Internal, "could not generate tokens")
	}
	return &pb.RegisterResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	claims, err := s.TokenManager.Validate(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	userId, ok := claims["userID"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no userId in token")
	}
	authToken, refreshToken, err := s.generateTokensWithID(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not generate tokens")
	}
	return &pb.RefreshTokenResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	user := GetContextUser(ctx)
	if user == nil {
		return nil, errors.NotLoggedIn
	}
	return &pb.UserInfoResponse{
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}, nil
}

func (s *AuthServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	err := s.TokenManager.InvalidateToken(ctx, req.AuthToken)
	if err != nil {
		logrus.Errorln(err)
		return nil, status.Error(codes.Internal, "could not invalidate auth token")
	}
	err = s.TokenManager.InvalidateToken(ctx, req.RefreshToken)
	if err != nil {
		logrus.Errorln(err)
		return nil, status.Error(codes.Internal, "could not invalidate refresh token")
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthServiceServer) generateTokensWithID(userId string) (string, string, error) {
	authToken, err := s.TokenManager.Generate(map[string]string{
		"userID": userId,
	}, time.Now().Add(time.Second*time.Duration(s.AuthTokenExpire)))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.TokenManager.Generate(map[string]string{
		"userID": userId,
	}, time.Now().Add(time.Second*time.Duration(s.RefreshTokenExpire)))
	if err != nil {
		return "", "", err
	}
	return authToken, refreshToken, nil
}

func (s *AuthServiceServer) generateTokens(user *storage.UserAccount) (string, string, error) {
	return s.generateTokensWithID(strconv.Itoa(int(user.ID)))
}
