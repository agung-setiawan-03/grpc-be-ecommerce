package handler

import (
	"context"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/service"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/cart"
)

type cartHandler struct {
	cart.UnimplementedCartServiceServer

	cartService service.ICartService
}

func (ch *cartHandler) AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &cart.AddProductToCartResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	res, err := ch.cartService.AddProductToCart(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *cartHandler) ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &cart.ListCartResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	res, err := ch.cartService.ListCart(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *cartHandler) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &cart.DeleteCartResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	res, err := ch.cartService.DeleteCart(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ch *cartHandler) UpdateCartQuantity(ctx context.Context, req *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error) {
	validationErors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErors != nil {
		return &cart.UpdateCartQuantityResponse{
			Base: utils.ValidationErrorResponse(validationErors),
		}, nil
	}

	res, err := ch.cartService.UpdateCartQuantity(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewCartHandler(cartService service.ICartService) *cartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}
