package transport

import (
	"context"

	"github.com/nacobas/customer/pb"
)

type grpcServer struct {
	pb.UnimplementedCustomerRegistryServer
}

func (gs *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	return nil, nil
}

func (gs *grpcServer) New(ctx context.Context, req *pb.NewRequest) (*pb.NewResponse, error) {
	return nil, nil
}

func (gs *grpcServer) Update(ctx context.Context, req *pb.UpdateInfoRequest) (*pb.UpdateInfoResponse, error) {
	return nil, nil
}

func (gs *grpcServer) SetState(ctx context.Context, req *pb.SetStateRequest) (*pb.SetStateResponse, error) {
	return nil, nil
}

func (gs *grpcServer) mustEmbedUnimplementedCustomerRegistryServer()
