package server

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
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
		return nil, status.Error(codes.NotFound, errors.New("user not found").Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.New("user not found").Error())
	}
	authToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.New("could not generate tokens").Error())
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
		return nil, status.Error(codes.AlreadyExists, errors.New("email already exists").Error())
	}
	_, err = s.Storage.GetUserByPhoneNumber(req.PhoneNumber)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, errors.New("phone number already exists").Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.New("could not hash password").Error())
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
		return nil, status.Error(codes.Internal, errors.New("could not create user").Error())
	}
	authToken, refreshToken, err := s.generateTokens(&user)
	if err != nil {
		logrus.Errorln(err)
		return nil, status.Error(codes.Internal, errors.New("could not generate tokens").Error())
	}
	return &pb.RegisterResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	claims, err := s.TokenManager.Validate(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.New("invalid token").Error())
	}
	userId, ok := claims["userID"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, errors.New("no userId in token").Error())
	}
	authToken, refreshToken, err := s.generateTokensWithID(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.New("could not generate tokens").Error())
	}
	return &pb.RefreshTokenResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	claims, err := s.TokenManager.Validate(ctx, req.AuthToken)
	if err != nil {
		logrus.Errorln(err)
		return nil, status.Error(codes.Unauthenticated, errors.New("invalid token").Error())
	}
	userId, ok := claims["userID"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, errors.New("no userId in token").Error())
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.New("wrong userId").Error())
	}
	user, err := s.Storage.GetUserByID(uint(id))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.New("user not found").Error())
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
		return nil, status.Error(codes.Internal, errors.New("could not invalidate auth token").Error())
	}
	err = s.TokenManager.InvalidateToken(ctx, req.RefreshToken)
	if err != nil {
		logrus.Errorln(err)
		return nil, status.Error(codes.Internal, errors.New("could not invalidate refresh token").Error())
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
