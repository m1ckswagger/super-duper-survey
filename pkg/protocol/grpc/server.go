package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	v1 "github.com/m1ckswagger/super-duper-survey/pkg/api/v1"
)

func RunServer(ctx context.Context, v1Catalog v1.CatalogServiceServer, v1User v1.UserServiceServer, v1Answer v1.AnswerServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterCatalogServiceServer(server, v1Catalog)
	v1.RegisterUserServiceServer(server, v1User)
	v1.RegisterAnswerServiceServer(server, v1Answer)

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
