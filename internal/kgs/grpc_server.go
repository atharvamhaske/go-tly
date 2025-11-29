package kgs

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

// grpcServer implements the generated KGS gRPC server interface.
// It must embed UnimplementedKGSServer by value to satisfy the KGSServer
// interface generated in service_grpc.pb.go.
type grpcServer struct {
	UnimplementedKGSServer
	svc Service
}

// NewGRPCServer creates a gRPC server instance bound to the given Service.
func NewGRPCServer(svc Service) *grpcServer {
	return &grpcServer{svc: svc}
}

// GetKey handles the gRPC request by delegating to the Service.
func (s *grpcServer) GetKey(ctx context.Context, _ *GetKeyRequest) (*GetKeyResponse, error) {
	key, err := s.svc.GenerateKey(ctx)
	if err != nil {
		return nil, err
	}
	return &GetKeyResponse{Key: key}, nil
}

// Run starts a standalone gRPC server for the KGS on the given address, e.g. ":50051".
func Run(addr string, svc Service) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterKGSServer(server, NewGRPCServer(svc))

	return server.Serve(lis)
}
