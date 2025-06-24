package grpcmiddleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()
	resp, err = handler(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return resp, err
}
