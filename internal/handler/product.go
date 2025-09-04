package handler

import (
	"context"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/service"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/product"
)

type productHandler struct {
	product.UnimplementedProductServiceServer

	productService service.IProductService
}

func (ph *productHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &product.CreateProductResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := ph.productService.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ph *productHandler) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &product.DetailProductResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := ph.productService.DetailProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ph *productHandler) EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &product.EditProductResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := ph.productService.EditProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ph *productHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &product.DeleteProductResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := ph.productService.DeleteProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ph *productHandler) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &product.ListProductResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := ph.productService.ListProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ph *productHandler) ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &product.ListProductAdminResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := ph.productService.ListProductAdmin(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ph *productHandler) HighlightProducts(ctx context.Context, req *product.HighlightProductsRequest) (*product.HighlightProductsResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &product.HighlightProductsResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := ph.productService.HighlightProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{
		productService: productService,
	}
}
