package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity"
	jwtentity "github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity/jwt"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/repository"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/cart"
	"github.com/google/uuid"
)

type ICartService interface {
	AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
	ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error)
	DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error)
	UpdateCartQuantity(ctx context.Context, req *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error)
}

type cartService struct {
	productRepository repository.IProductRepository
	cartRepository    repository.ICartRepository
}

func (cs *cartService) AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	productEntity, err := cs.productRepository.GetProductById(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("Produk tidak ditemukan"),
		}, nil
	}

	cartEntity, err := cs.cartRepository.GetCartByProductAndUserId(ctx, req.ProductId, claims.Subject)
	if err != nil {
		return nil, err
	}

	if cartEntity != nil {
		now := time.Now()
		cartEntity.Quantity += 1
		cartEntity.UpdatedAt = &now
		cartEntity.UpdatedBy = &claims.FullName

		err = cs.cartRepository.UpdateCart(ctx, cartEntity)
		if err != nil {
			return nil, err
		}

		return &cart.AddProductToCartResponse{
			Base: utils.SuccessResponse("Berhasil menambahkan produk ke keranjang belanja"),
			Id:   cartEntity.Id,
		}, nil
	}

	newCartEntity := entity.UserCart{
		Id:        uuid.NewString(),
		UserId:    claims.Subject,
		ProductId: req.ProductId,
		Quantity:  1,
		CreatedAt: time.Now(),
		CreatedBy: claims.FullName,
	}

	err = cs.cartRepository.CreateNewCart(ctx, &newCartEntity)
	if err != nil {
		return nil, err
	}

	return &cart.AddProductToCartResponse{
		Base: utils.SuccessResponse("Berhasil menambahkan produk ke keranjang belanja"),
		Id:   newCartEntity.Id,
	}, nil
}

func (cs *cartService) ListCart(ctx context.Context, req *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	carts, err := cs.cartRepository.GetListCart(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	var items []*cart.ListCartResponseItem = make([]*cart.ListCartResponseItem, 0)
	for _, cartEntity := range carts {
		item := cart.ListCartResponseItem{
			CartId:          cartEntity.Id,
			ProductId:       cartEntity.ProductId,
			ProductName:     cartEntity.Product.Name,
			ProductImageUrl: fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), cartEntity.Product.ImageFileName),
			ProductPrice:    cartEntity.Product.Price,
			Quantity:        int64(cartEntity.Quantity),
		}

		items = append(items, &item)
	}

	return &cart.ListCartResponse{
		Base:  utils.SuccessResponse("Berhasil mendapatkan data keranjang belanja"),
		Items: items,
	}, nil
}

func (cs *cartService) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	// dapatkan data user id
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// dapatkan data cart
	cartEntity, err := cs.cartRepository.GetCartById(ctx, req.CartId)
	if err != nil {
		return nil, err
	}

	if cartEntity == nil {
		return &cart.DeleteCartResponse{
			Base: utils.NotFoundResponse("Keranjang belanja tidak ditemukan"),
		}, nil
	}

	// cocokan data user id dengan auth
	if cartEntity.UserId != claims.Subject {
		// return bad request
		return &cart.DeleteCartResponse{
			Base: utils.BadRequestResponse("Anda tidak memiliki akses untuk menghapus keranjang belanja ini"),
		}, nil
	}

	// delete
	err = cs.cartRepository.DeleteCart(ctx, req.CartId)
	if err != nil {
		return nil, err
	}

	// kirim response
	return &cart.DeleteCartResponse{
		Base: utils.SuccessResponse("Berhasil menghapus keranjang belanja"),
	}, nil
}

func (cs *cartService) UpdateCartQuantity(ctx context.Context, req *cart.UpdateCartQuantityRequest) (*cart.UpdateCartQuantityResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// dapatkan data cart by id
	cartEntity, err := cs.cartRepository.GetCartById(ctx, req.CartId)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.UpdateCartQuantityResponse{
			Base: utils.NotFoundResponse("Keranjang belanja tidak ditemukan"),
		}, nil
	}
	// cocokkan user id
	if cartEntity.UserId != claims.Subject {
		return &cart.UpdateCartQuantityResponse{
			Base: utils.BadRequestResponse("Anda tidak memiliki akses untuk mengubah keranjang belanja ini"),
		}, nil
	}

	// update new quantity
	if req.NewQuantity == 0 {
		cs.cartRepository.DeleteCart(ctx, cartEntity.Id)
		return &cart.UpdateCartQuantityResponse{
			Base: utils.SuccessResponse("Berhasil menghapus keranjang belanja, keranjang belanja anda kosong"),
		}, nil
	}
	now := time.Now()
	cartEntity.Quantity = int(req.NewQuantity)
	cartEntity.UpdatedAt = &now
	cartEntity.UpdatedBy = &claims.FullName

	// update ke db
	err = cs.cartRepository.UpdateCart(ctx, cartEntity)
	if err != nil {
		return nil, err
	}

	// return response
	return &cart.UpdateCartQuantityResponse{
		Base: utils.SuccessResponse("Berhasil mengubah jumlah keranjang belanja"),
	}, nil
}

func NewCartService(productRepository repository.IProductRepository, cartRepository repository.ICartRepository) ICartService {
	return &cartService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
	}
}
