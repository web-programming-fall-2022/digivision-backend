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
		p, err := s.fetcher.Fetch(ctx, product)
		if err != nil {
			return nil, err
		}
		p.Score = product.Score
		return p, nil
	}, products[:req.TopK])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch products: %v", err)
	}
	return &pb.SearchResponse{
		Products: fetchedProducts,
	}, nil
}

func (s *SearchServiceServer) AsyncSearch(req *pb.SearchRequest, stream pb.SearchService_AsyncSearchServer) error {
	vector, err := s.img2vec.Vectorize(stream.Context(), req.Image)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to vectorize the image: %v", err)
	}
	productImages, err := s.searchHandler.Search(stream.Context(), vector, int(req.TopK))
	if err != nil {
		return status.Errorf(codes.Internal, "failed to search: %v", err)
	}
	products := s.ranker.Rank(productImages)
	respChan := s.fetcher.AsyncFetch(stream.Context(), products, int(req.TopK))
	for {
		select {
		case <-stream.Context().Done():
			return status.Errorf(codes.Canceled, "client canceled the request")
		case resp := <-respChan:
			if resp == nil {
				return nil
			}
			if resp.Product != nil {
				if err := stream.Send(&pb.AsyncSearchResponse{Product: resp.Product}); err != nil {
					return status.Errorf(codes.Internal, "failed to send product: %v", err)
				}
			}
		}
	}
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
