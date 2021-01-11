package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "sipwise.com/shardik/pkg/api/v1"
)

func RunServer(ctx context.Context, v1API v1.CatalogServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grcp.NewServer()
	v1.RegisterCatalogServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
