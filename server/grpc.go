package server

import (
	"context"

	pb "github.com/Stoakes/grpc-gateway-example/echopb"

	"google.golang.org/grpc"
)

func prepareGRPC(context context.Context) (*grpc.Server, error) {

	grpcServer := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpcServer, new(MicroServer))

	return grpcServer, nil
}
