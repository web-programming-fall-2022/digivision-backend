package server

import (
	"context"
	"github.com/arimanius/digivision-backend/internal/img2vec"
	"github.com/arimanius/digivision-backend/internal/od"
	"github.com/arimanius/digivision-backend/internal/productmeta"
	"github.com/arimanius/digivision-backend/internal/rank"
	"github.com/arimanius/digivision-backend/internal/search"
	"github.com/arimanius/digivision-backend/internal/utils"
	pb "github.com/arimanius/digivision-backend/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServiceServer struct {
	pb.UnimplementedSearchServiceServer
	img2vec        img2vec.Img2Vec
	searchHandler  search.Handler
	ranker         rank.Ranker
	objectDetector od.ObjectDetector
	fetcher        productmeta.Fetcher
}

func NewSearchServiceServer(
	i2v img2vec.Img2Vec,
	searchHandler search.Handler,
	fetcher productmeta.Fetcher,
	ranker rank.Ranker,
	objectDetector od.ObjectDetector,
) *SearchServiceServer {
	return &SearchServiceServer{
		img2vec:        i2v,
		searchHandler:  searchHandler,
		fetcher:        fetcher,
		ranker:         ranker,
		objectDetector: objectDetector,
	}
}

func (s *SearchServiceServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	vector, err := s.img2vec.Vectorize(ctx, req.Image)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to vectorize the image: %v", err)
	}
	productImages, err := s.searchHandler.Search(ctx, vector, int(req.TopK))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search: %v", err)
	}
	products := s.ranker.Rank(productImages)
	fetchedProducts, err := utils.ConcurrentMap(func(product rank.Product) (*pb.Product, error) {
		p, err := s.fetcher.Fetch(ctx, product.Id)
		if err != nil {
			return nil, err
		}
		p.Score = product.Distance
		return p, nil
	}, products)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch products: %v", err)
	}
	return &pb.SearchResponse{
		Products: fetchedProducts,
	}, nil
}

func (s *SearchServiceServer) Crop(ctx context.Context, req *pb.CropRequest) (*pb.CropResponse, error) {
	topLeft, bottomRight, err := s.objectDetector.Detect(ctx, req.Image)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to detect object: %v", err)
	}
	return &pb.CropResponse{
		TopLeft: &pb.Position{
			X: int32(topLeft.X),
			Y: int32(topLeft.Y),
		},
		BottomRight: &pb.Position{
			X: int32(bottomRight.X),
			Y: int32(bottomRight.Y),
		},
	}, nil
}
