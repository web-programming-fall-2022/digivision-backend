package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/web-programming-fall-2022/digivision-backend/internal/img2vec"
	"github.com/web-programming-fall-2022/digivision-backend/internal/od"
	"github.com/web-programming-fall-2022/digivision-backend/internal/productmeta"
	"github.com/web-programming-fall-2022/digivision-backend/internal/rank"
	"github.com/web-programming-fall-2022/digivision-backend/internal/s3"
	"github.com/web-programming-fall-2022/digivision-backend/internal/search"
	"github.com/web-programming-fall-2022/digivision-backend/internal/storage"
	pb "github.com/web-programming-fall-2022/digivision-backend/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServiceServer struct {
	pb.UnimplementedSearchServiceServer
	img2vec        img2vec.Img2Vec
	searchHandler  search.Handler
	rankers        map[pb.Ranker]rank.Ranker
	objectDetector od.ObjectDetector
	fetcher        productmeta.Fetcher
	s3Client       s3.S3Client
	storage        *storage.Storage
}

func NewSearchServiceServer(
	i2v img2vec.Img2Vec,
	searchHandler search.Handler,
	fetcher productmeta.Fetcher,
	rankers map[pb.Ranker]rank.Ranker,
	objectDetector od.ObjectDetector,
	s3Client s3.S3Client,
) *SearchServiceServer {
	return &SearchServiceServer{
		img2vec:        i2v,
		searchHandler:  searchHandler,
		fetcher:        fetcher,
		rankers:        rankers,
		objectDetector: objectDetector,
		s3Client:       s3Client,
	}
}

func (s *SearchServiceServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var history *storage.SearchHistory
	user := GetContextUser(ctx)
	if user != nil {
		path := fmt.Sprintf("%s.jpg", uuid.New().String())
		err := s.s3Client.Upload(ctx, "history-images", path, bytes.NewReader(req.Image), int64(len(req.Image)))
		if err != nil {
			logrus.Errorf("failed to upload image to s3: %v", err)
		} else {
			history = &storage.SearchHistory{
				UserID:       user.ID,
				QueryAddress: path,
			}
			err = s.storage.CreateSearchHistory(history)
		}
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload the image: %v", err)
	}
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
				if history != nil && history.ID != 0 {
					err := s.storage.CreateSearchHistoryResult(&storage.SearchHistoryResult{
						SearchHistoryID: history.ID,
						ProductID:       uint(resp.Product.Id),
					})
					if err != nil {
						logrus.Errorf("failed to create search history result: %v", err)
					}
				}
			}
		}
	}
}

func (s *SearchServiceServer) AsyncSearch(req *pb.SearchRequest, stream pb.SearchService_AsyncSearchServer) error {
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
