package server

import (
	"context"
	"github.com/web-programming-fall-2022/digivision-backend/internal/errors"
	"github.com/web-programming-fall-2022/digivision-backend/internal/productmeta"
	"github.com/web-programming-fall-2022/digivision-backend/internal/rank"
	"github.com/web-programming-fall-2022/digivision-backend/internal/storage"
	pb "github.com/web-programming-fall-2022/digivision-backend/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type FavoriteServiceServer struct {
	pb.UnimplementedFavoriteServiceServer
	Storage *storage.Storage
	fetcher productmeta.Fetcher
}

func NewFavoriteServiceServer(
	storage *storage.Storage,
	fetcher productmeta.Fetcher,
) *FavoriteServiceServer {
	return &FavoriteServiceServer{
		Storage: storage,
		fetcher: fetcher,
	}
}

func (s *FavoriteServiceServer) AddItemToFavorites(
	ctx context.Context, req *pb.AddItemToFavoritesRequest,
) (*pb.AddItemToFavoritesResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := GetContextUser(ctx)
	if user == nil {
		return nil, errors.NotLoggedIn
	}

	list, err := s.Storage.GetFavoriteListByUserIDAndName(user.ID, req.ListName)
	if err != nil {
		return nil, errors.NotFound
	}
	err = s.Storage.AddItemToList(list.ID, uint(req.ProductId))
	if err != nil {
		return nil, errors.Internal
	}
	return &pb.AddItemToFavoritesResponse{Success: true}, nil
}

func (s *FavoriteServiceServer) RemoveItemFromFavorites(
	ctx context.Context, req *pb.RemoveItemFromFavoritesRequest,
) (*pb.RemoveItemFromFavoritesResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := GetContextUser(ctx)
	if user == nil {
		return nil, errors.NotLoggedIn
	}

	list, err := s.Storage.GetFavoriteListByUserIDAndName(user.ID, req.ListName)
	if err != nil {
		return nil, errors.NotFound
	}
	err = s.Storage.RemoveItemFromList(list.ID, uint(req.ProductId))
	if err != nil {
		return nil, errors.NotFound
	}
	return &pb.RemoveItemFromFavoritesResponse{Success: true}, nil
}

func (s *FavoriteServiceServer) GetFavorites(
	ctx context.Context, req *pb.GetFavoritesRequest,
) (*pb.GetFavoritesResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := GetContextUser(ctx)
	if user == nil {
		return nil, errors.NotLoggedIn
	}

	list, err := s.Storage.GetFavoriteListByUserIDAndName(user.ID, req.ListName)
	if err != nil {
		return nil, errors.NotFound
	}
	uniqueProductIDs := make(map[uint]interface{})
	products := make([]rank.Product, 0)
	for _, item := range list.Items {
		_, ok := uniqueProductIDs[item.ProductID]
		if ok {
			continue
		}
		uniqueProductIDs[item.ProductID] = nil
		products = append(products, rank.Product{
			Id:    strconv.Itoa(int(item.ProductID)),
			Score: 0,
		})
	}
	productsChan := s.fetcher.AsyncFetch(ctx, products, len(products))
	productsList := make([]*pb.Product, 0)
	for product := range productsChan {
		if product == nil {
			continue
		}
		productsList = append(productsList, product.Product)
	}

	return &pb.GetFavoritesResponse{Products: productsList}, nil
}
