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
	mux "github.com/gorilla/mux"
	assetfs "github.com/philips/go-bindata-assetfs"
)

func prepareHTTP(ctx context.Context, serverName string) (*http.Server, error) {
	// gorilla mux instead of default net/http handler
	router := mux.NewRouter()
	router.HandleFunc("/bonjour/{name}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bonjour " + vars["name"]))
	})
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, strings.NewReader(pb.Swagger))
	})
	serveSwagger(router)

	// initialize grpc-gateway
	gw, err := prepareGateway(ctx)
	if err != nil {
		return nil, err
	}

	// Add a default matcher for everything that begins with /
	router.PathPrefix("/").Handler(gw)

	/**
	 * Need to debug a router ? https://github.com/gorilla/mux/issues/186#issuecomment-249347077
	 */

	// Return HTTP Server instance
	return &http.Server{
		Addr:         serverName,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, nil
}

func serveSwagger(router *mux.Router) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, fileServer))
}
