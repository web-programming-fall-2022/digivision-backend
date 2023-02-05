package server

import (
	"context"
	"fmt"
	"github.com/arimanius/digivision-backend/internal/img2vec"
	"github.com/arimanius/digivision-backend/internal/od"
	"github.com/arimanius/digivision-backend/internal/productmeta"
	"github.com/arimanius/digivision-backend/internal/rank"
	"github.com/arimanius/digivision-backend/internal/search"
	pb "github.com/arimanius/digivision-backend/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type SearchServiceServer struct {
	pb.UnimplementedSearchServiceServer
	img2vec        img2vec.Img2Vec
	searchHandler  search.Handler
	rankers        map[pb.Ranker]rank.Ranker
	objectDetector od.ObjectDetector
	fetcher        productmeta.Fetcher
	logSearchImage bool
}

func NewSearchServiceServer(
	i2v img2vec.Img2Vec,
	searchHandler search.Handler,
	fetcher productmeta.Fetcher,
	rankers map[pb.Ranker]rank.Ranker,
	objectDetector od.ObjectDetector,
	logSearchImage bool,
) *SearchServiceServer {
	return &SearchServiceServer{
		img2vec:        i2v,
		searchHandler:  searchHandler,
		fetcher:        fetcher,
		rankers:        rankers,
		objectDetector: objectDetector,
	}
}

func (s *SearchServiceServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	s.saveImageToFile(req.Image)
	vector, err := s.img2vec.Vectorize(ctx, req.Image)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to vectorize the image: %v", err)
	}
	productImages, err := s.searchHandler.Search(ctx, vector, int(req.Params.TopK))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search: %v", err)
	}
	products := s.rankers[req.Params.Ranker].Rank(productImages)
	respChan := s.fetcher.AsyncFetch(ctx, products, int(req.Params.TopK))
	var resultProducts []*pb.Product
	for {
		select {
		case <-ctx.Done():
			return nil, status.Errorf(codes.Canceled, "client canceled the request")
		case resp := <-respChan:
			if resp == nil {
				return &pb.SearchResponse{Products: resultProducts}, nil
			}
			if resp.Product != nil {
				resultProducts = append(resultProducts, resp.Product)
			}
		}
	}
}

func (s *SearchServiceServer) AsyncSearch(req *pb.SearchRequest, stream pb.SearchService_AsyncSearchServer) error {
	s.saveImageToFile(req.Image)
	vector, err := s.img2vec.Vectorize(stream.Context(), req.Image)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to vectorize the image: %v", err)
	}
	productImages, err := s.searchHandler.Search(stream.Context(), vector, int(req.Params.TopK))
	if err != nil {
		return status.Errorf(codes.Internal, "failed to search: %v", err)
	}
	products := s.rankers[req.Params.Ranker].Rank(productImages)
	respChan := s.fetcher.AsyncFetch(stream.Context(), products, int(req.Params.TopK))
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

func (s *SearchServiceServer) saveImageToFile(image []byte) error {
	if !s.logSearchImage {
		return nil
	}
	time := time.Now().Format("2006-01-02_15-04-05")
	return os.WriteFile(fmt.Sprintf("./data/%s.png", time), image, 0644)
}
