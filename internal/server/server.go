package server

import (
	"context"
	pb "github.com/arimanius/digivision-backend/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServiceServer struct {
	pb.UnimplementedSearchServiceServer
}

func NewSearchServiceServer() *SearchServiceServer {
	return &SearchServiceServer{}
}

func (s *SearchServiceServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (s *SearchServiceServer) Crop(ctx context.Context, req *pb.CropRequest) (*pb.CropResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
