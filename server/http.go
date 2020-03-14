package server

import (
	"context"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	pb "github.com/Stoakes/grpc-gateway-example/echopb"
	"github.com/Stoakes/grpc-gateway-example/pkg/ui/data/swagger"
	assetfs "github.com/philips/go-bindata-assetfs"
)

func prepareHTTP(ctx context.Context, serverName string) (*http.Server, error) {
	// HTTP router
	router := http.NewServeMux()
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, strings.NewReader(pb.Swagger))
	})
	serveSwagger(router)

	// initialize grpc-gateway
	gw, err := prepareGateway(ctx)
	if err != nil {
		return nil, err
	}
	router.Handle("/", gw)

	// Return HTTP Server instance
	return &http.Server{
		Addr:         serverName,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, nil
}

func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
