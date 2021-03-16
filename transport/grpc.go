package transport

import (
	"context"

	"github.com/nacobas/customer/pb"
)

type grpcServer struct {
	pb.UnimplementedCustomerRegistryServer
}

func (gs *grpcServer) Get(context.Context, *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	return nil, nil
}

func (gs *grpcServer) New(context.Context, *pb.NewCustomerRequest) (*pb.NewCustomerResponse, error) {
	return nil, nil
}

func (gs *grpcServer) Update(context.Context, *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {
	return nil, nil
}

func (gs *grpcServer) Close(context.Context, *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	return nil, nil
}

func (gs *grpcServer) mustEmbedUnimplementedCustomerRegistryServer()
