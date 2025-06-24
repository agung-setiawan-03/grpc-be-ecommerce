package handler

import (
	"context"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/service"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/auth"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer

	authService service.IAuthService
}

func (sh *authHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	// Process Register
	resp, err := sh.authService.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
