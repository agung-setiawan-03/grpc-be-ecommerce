package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/handler"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/service"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pkg/database"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pkg/grpcmiddleware"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panicf("Error when listening server: %v", err)
	}

	database.ConnectDB(ctx, os.Getenv("DB_URL"))
	log.Println("Connected to database")

	serviceHandler := handler.NewServiceHandler()

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrMiddleware,
		),
	)

	service.RegisterHelloWorldServiceServer(serv, serviceHandler)

	if os.Getenv("ENVIRONMENT") == "development" {
		reflection.Register(serv)
		log.Println("Reflection is registered")
	}

	log.Println("Server is running on port 8080")
	if err := serv.Serve(list); err != nil {
		log.Panicf("Server is error: %v", err)
	}

}
