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

func (sh *authHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &auth.LoginResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	// Process Login
	resp, err := sh.authService.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (sh *authHandler) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &auth.LogoutResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	// Process Logout
	resp, err := sh.authService.Logout(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (sh *authHandler) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	// Process Change Password
	resp, err := sh.authService.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (sh *authHandler) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	resp, err := sh.authService.GetProfile(ctx, req)
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
