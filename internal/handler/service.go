package handler

import (
	"context"
	"fmt"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh *serviceHandler) HelloWorld(ctx context.Context, request *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	validationErors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &service.HelloWorldResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello " + request.Name),
		Base:    utils.SuccessResponse("Success"),
	}, nil

}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}
