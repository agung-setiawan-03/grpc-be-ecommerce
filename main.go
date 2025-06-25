package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/grpcmiddleware"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/handler"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/repository"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/service"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/auth"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pkg/database"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
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

	db := database.ConnectDB(ctx, os.Getenv("DB_URL"))
	log.Println("Connected to database")

	cacheService := gocache.New(time.Hour*24, time.Hour)

	authMiddleware := grpcmiddleware.NewAuthMiddleware(cacheService)

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrMiddleware,
			authMiddleware.AuthMiddleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)

	if os.Getenv("ENVIRONMENT") == "development" {
		reflection.Register(serv)
		log.Println("Reflection is registered")
	}

	log.Println("Server is running on port 8080")
	if err := serv.Serve(list); err != nil {
		log.Panicf("Server is error: %v", err)
	}

}
