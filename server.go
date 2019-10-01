package main

import (
	"context"

	"github.com/prometheus/common/log"
	"github.com/wdullaer/whoami/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedWhoamiServiceServer
	health codes.Code
}

func (s *server) Bench(ctx context.Context, req *pb.BenchRequest) (*pb.BenchResponse, error) {
	if ctx.Err() != nil {
		return nil, status.FromContextError(ctx.Err()).Err()
	}

	log.Infof("Bench is called with: %s", req.String())

	return &pb.BenchResponse{
		Result: 1,
	}, nil
}

func (s *server) GetHealth(ctx context.Context, req *pb.GetHealthRequest) (*pb.GetHealthResponse, error) {
	if ctx.Err() != nil {
		return nil, status.FromContextError(ctx.Err()).Err()
	}

	log.Infof("GetHealth is called with: %s", req.String())

	if s.health == codes.OK {
		return &pb.GetHealthResponse{
			Ok: true,
		}, nil
	}

	return nil, status.Error(s.health, s.health.String())
}

func (s *server) SetHealth(ctx context.Context, req *pb.SetHealthRequest) (*pb.SetHealthResponse, error) {
	if ctx.Err() != nil {
		return nil, status.FromContextError(ctx.Err()).Err()
	}

	log.Infof("SetHealth is called with: %s", req.String())

	s.health = codes.Code(req.GetStatus())

	return &pb.SetHealthResponse{}, nil
}
