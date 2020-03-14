package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pb "github.com/Stoakes/grpc-gateway-example/echopb"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func prepareGateway(ctx context.Context) (http.Handler, error) {
	// gRPC dialup options
	opts := []grpc.DialOption{
		grpc.WithTimeout(10 * time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(DemoAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("Fail to dial: %v", err)
	}

	// changes json serializer to include empty fields with default values
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)

	// Register Gateway endpoints
	err = pb.RegisterEchoServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, err
	}

	return gwMux, nil
}
