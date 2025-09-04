package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity"
	jwtentity "github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity/jwt"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/repository"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/product"
	"github.com/google/uuid"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error)
	DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error)
	EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error)
	DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error)
	ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error)
	ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error)
	HighlightProducts(ctx context.Context, req *product.HighlightProductsRequest) (*product.HighlightProductsResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
}

func (ps *productService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	imagePath := filepath.Join("storage", "product", req.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &product.CreateProductResponse{
				Base: utils.BadRequestResponse("Gambar produk tidak ditemukan"),
			}, nil
		}
		return nil, err
	}

	productEntity := entity.Product{
		Id:            uuid.NewString(),
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		ImageFileName: req.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.FullName,
	}
	err = ps.productRepository.CreateNewProduct(ctx, &productEntity)
	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Produk Berhasil dibuat"),
		Id:   productEntity.Id,
	}, nil
}

func (ps *productService) DetailProduct(ctx context.Context, req *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	productentity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productentity == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Produk tidak ditemukan"),
		}, nil
	}

	return &product.DetailProductResponse{
		Base:        utils.SuccessResponse("Detail produk berhasil didapatkan"),
		Id:          productentity.Id,
		Name:        productentity.Name,
		Description: productentity.Description,
		Price:       productentity.Price,
		ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), productentity.ImageFileName),
	}, nil

}

func (ps *productService) EditProduct(ctx context.Context, req *product.EditProductRequest) (*product.EditProductResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.EditProductResponse{
			Base: utils.NotFoundResponse("Produk tidak ditemukan"),
		}, nil
	}

	if productEntity.ImageFileName != req.ImageFileName {
		newImagePath := filepath.Join("storage", "product", req.ImageFileName)
		_, err := os.Stat(newImagePath)
		if err != nil {
			if os.IsNotExist(err) {
				return &product.EditProductResponse{
					Base: utils.BadRequestResponse("Gambar tidak ditemukan"),
				}, nil
			}

			return nil, err
		}

		oldImagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
		err = os.Remove(oldImagePath)
		if err != nil {
			return nil, err
		}
	}

	newProduct := entity.Product{
		Id:            req.Id,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		ImageFileName: req.ImageFileName,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &claims.FullName,
	}

	err = ps.productRepository.UpdateProduct(ctx, &newProduct)
	if err != nil {
		return nil, err
	}

	return &product.EditProductResponse{
		Base: utils.SuccessResponse("Edit produk berhasil"),
		Id:   req.Id,
	}, nil

}

func (ps *productService) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	productEntity, err := ps.productRepository.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.DeleteProductResponse{
			Base: utils.NotFoundResponse("Produk tidak ditemukan"),
		}, nil
	}

	err = ps.productRepository.DeleteProduct(ctx, req.Id, time.Now(), claims.FullName)
	if err != nil {
		return nil, err
	}

	imagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
	err = os.Remove(imagePath)
	if err != nil {
		return nil, err
	}

	return &product.DeleteProductResponse{
		Base: utils.SuccessResponse("Produk berhasil dihapus"),
	}, nil

}

func (ps *productService) ListProduct(ctx context.Context, req *product.ListProductRequest) (*product.ListProductResponse, error) {
	products, paginationResponse, err := ps.productRepository.GetProductsPagination(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*product.ListProductResponseItem = make([]*product.ListProductResponseItem, 0)
	for _, prod := range products {
		data = append(data, &product.ListProductResponseItem{
			Id:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("SOTRAGE_SERVICE_URL"), prod.ImageFileName),
		})
	}

	return &product.ListProductResponse{
		Base:       utils.SuccessResponse("Produk berhasil diambil"),
		Pagination: paginationResponse,
		Data:       data,
	}, nil

}

func (ps *productService) ListProductAdmin(ctx context.Context, req *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	products, paginationResponse, err := ps.productRepository.GetProductsPaginationAdmin(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*product.ListProductAdminResponseItem = make([]*product.ListProductAdminResponseItem, 0)
	for _, prod := range products {
		data = append(data, &product.ListProductAdminResponseItem{
			Id:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("SOTRAGE_SERVICE_URL"), prod.ImageFileName),
		})
	}

	return &product.ListProductAdminResponse{
		Base:       utils.SuccessResponse("Produk berhasil diambil"),
		Pagination: paginationResponse,
		Data:       data,
	}, nil

}

func (ps *productService) HighlightProducts(ctx context.Context, req *product.HighlightProductsRequest) (*product.HighlightProductsResponse, error) {
	products, err := ps.productRepository.GetProductsHighlight(ctx)
	if err != nil {
		return nil, err
	}

	var data []*product.HighlightProductsResponseItem = make([]*product.HighlightProductsResponseItem, 0)
	for _, prod := range products {
		data = append(data, &product.HighlightProductsResponseItem{
			Id:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("SOTRAGE_SERVICE_URL"), prod.ImageFileName),
		})
	}

	return &product.HighlightProductsResponse{
		Base: utils.SuccessResponse("Highlight Produk berhasil diambil"),
		Data: data,
	}, nil

}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
