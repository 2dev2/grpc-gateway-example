package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/Stoakes/grpc-gateway-example/echopb"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// MicroServer represents a microservice server instance
type MicroServer struct {
	serverName string
	lis        net.Listener
	httpServer *http.Server
	grpcServer *grpc.Server
}

// New returns a microserver instance
func New(serverName string, l net.Listener) *MicroServer {
	return &MicroServer{
		serverName: serverName,
		lis:        l,
	}
}

// Echo implementation
// Implement EchoService interface
func (m *MicroServer) Echo(c context.Context, s *pb.EchoMessage) (*pb.EchoMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}

// -----------------------------------------------------------------------------

// Start starts the microserver
func (m *MicroServer) Start() error {
	var err error

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// tcpMuxer
	tcpMux := cmux.New(m.lis)

	// Connection dispatcher rules
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())

	go func(ms *MicroServer, grpcL net.Listener) {
		// initialize gRPC server instance
		m.grpcServer, err = prepareGRPC(ctx)
		if err != nil {
			log.Fatalf("Unable to initialize gRPC server instance: %s", err.Error())
		}
		if err = m.grpcServer.Serve(grpcL); err != nil {
			log.Fatalf("Unable to start external gRPC server: %s", err.Error())
		}
	}(m, grpcL)

	go func(ms *MicroServer, httpL net.Listener) {
		// initialize HTTP server
		m.httpServer, err = prepareHTTP(ctx, ms.serverName)
		if err != nil {
			log.Fatalf("Unable to initialize HTTP server instance: %s", err.Error())
		}
		if err = m.httpServer.Serve(httpL); err != nil {
			log.Fatalf("Unable to start HTTP server: %s", err.Error())
		}
	}(m, httpL)

	return tcpMux.Serve()
}
