package service

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	operatingsystem "os"
	"runtime/debug"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity"
	jwtentity "github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity/jwt"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/repository"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/order"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
	ListOrderAdmin(ctx context.Context, req *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error)
	ListOrder(ctx context.Context, req *order.ListOrderRequest) (*order.ListOrderResponse, error)
}

type orderService struct {
	db                *sql.DB
	orderRepository   repository.IOrderRepository
	productRepository repository.IProductRepository
}

func (os *orderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := os.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				tx.Rollback()
			}

			debug.PrintStack()
			panic(e)
		}
	}()

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	orederRepo := os.orderRepository.WithTransaction(tx)
	productRepo := os.productRepository.WithTransaction(tx)

	numbering, err := orederRepo.GetNumbering(ctx, "order")
	if err != nil {
		return nil, err
	}

	var productIds = make([]string, len(req.Products))
	for i := range req.Products {
		productIds[i] = req.Products[i].Id
	}

	products, err := productRepo.GetProducstByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]*entity.Product)
	for i := range products {
		productMap[products[i].Id] = products[i]
	}

	var total float64 = 0
	for _, p := range req.Products {
		if productMap[p.Id] == nil {
			return &order.CreateOrderResponse{
				Base: utils.BadRequestResponse(fmt.Sprintf("Product dengan ID %s tidak ditemukan", p.Id)),
			}, nil
		}
		total += productMap[p.Id].Price * float64(p.Quantity)
	}

	now := time.Now()
	prefix := "INV"
	store := "RK"
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(90000000) + 10000000
	expiretAt := now.Add(24 * time.Hour)
	orderEntity := entity.Order{
		Id:              uuid.NewString(),
		Number:          fmt.Sprintf("%s/%s/%s/%d%d", prefix, now.Format("20060102"), store, randomNumber, numbering.Number),
		UserId:          claims.Subject,
		OrderStatusCode: entity.OrderStatusCodeUnpaid,
		UserFullName:    req.FullName,
		Address:         req.Address,
		PhoneNumber:     req.PhoneNumber,
		Notes:           &req.Notes,
		Total:           total,
		ExpiredAt:       &expiretAt,
		CreatedAt:       now,
		CreatedBy:       claims.FullName,
	}

	invoiceItems := make([]xendit.InvoiceItem, 0)
	for _, p := range req.Products {
		prod := productMap[p.Id]
		if prod != nil {
			invoiceItems = append(invoiceItems, xendit.InvoiceItem{
				Name:     prod.Name,
				Price:    prod.Price,
				Quantity: int(p.Quantity),
			})
		}
	}

	xenditInvoice, xenditErr := invoice.CreateWithContext(ctx, &invoice.CreateParams{
		ExternalID: orderEntity.Id,
		Amount:     total,
		Customer: xendit.InvoiceCustomer{
			GivenNames: req.FullName,
		},
		Currency:           "IDR",
		SuccessRedirectURL: fmt.Sprintf("%s/checkout/%s/success", operatingsystem.Getenv("FRONTEND_BASE_URL"), orderEntity.Id),
		Items:              invoiceItems,
	})

	if xenditErr != nil {
		err = xenditErr
		return nil, err
	}

	orderEntity.XenditInvoiceId = &xenditInvoice.ID
	orderEntity.XenditInvoiceUrl = &xenditInvoice.InvoiceURL

	err = orederRepo.CreateOrder(ctx, &orderEntity)
	if err != nil {
		return nil, err
	}

	for _, p := range req.Products {
		var orderItem = entity.OrderItem{
			Id:                   uuid.NewString(),
			ProductId:            p.Id,
			ProductName:          productMap[p.Id].Name,
			ProductImageFileName: productMap[p.Id].ImageFileName,
			ProductPrice:         productMap[p.Id].Price,
			Quantity:             p.Quantity,
			OrderId:              orderEntity.Id,
			CreatedAt:            now,
			CreatedBy:            claims.FullName,
		}

		err = orederRepo.CreateOrderItem(ctx, &orderItem)
		if err != nil {
			return nil, err
		}
	}

	numbering.Number++
	err = orederRepo.UpdateNumbering(ctx, numbering)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Base: utils.SuccessResponse("Order Berhasil Dibuat"),
		Id:   orderEntity.Id,
	}, nil
}

func (os *orderService) ListOrderAdmin(ctx context.Context, req *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	orders, metadata, err := os.orderRepository.GetListOrderAdminPagination(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderAdminResponseItem, 0)
	for _, o := range orders {
		products := make([]*order.ListOrderAdminResponseItemProduct, 0)
		for _, oi := range o.Items {
			products = append(products, &order.ListOrderAdminResponseItemProduct{
				Id:       oi.ProductId,
				Name:     oi.ProductName,
				Price:    oi.ProductPrice,
				Quantity: oi.Quantity,
			})
		}

		orderStatusCode := o.OrderStatusCode
		if o.OrderStatusCode == entity.OrderStatusCodeUnpaid && time.Now().After(*o.ExpiredAt) {
			orderStatusCode = entity.OrderStatusExpired
		}

		items = append(items, &order.ListOrderAdminResponseItem{
			Id:         o.Id,
			Number:     o.Number,
			Customer:   o.UserFullName,
			StatusCode: orderStatusCode,
			Total:      o.Total,
			CreatedAt:  timestamppb.New(o.CreatedAt),
			Products:   products,
		})
	}

	return &order.ListOrderAdminResponse{
		Base:       utils.SuccessResponse("Berhasil Mendapatkan List Order"),
		Pagination: metadata,
		Items:      items,
	}, nil
}

func (os *orderService) ListOrder(ctx context.Context, req *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	orders, metadata, err := os.orderRepository.GetListOrderPagination(ctx, req.Pagination, claims.Subject)
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderResponseItem, 0)
	for _, o := range orders {
		products := make([]*order.ListOrderResponseItemProduct, 0)
		for _, oi := range o.Items {
			products = append(products, &order.ListOrderResponseItemProduct{
				Id:       oi.ProductId,
				Name:     oi.ProductName,
				Price:    oi.ProductPrice,
				Quantity: oi.Quantity,
			})
		}

		orderStatusCode := o.OrderStatusCode
		if o.OrderStatusCode == entity.OrderStatusCodeUnpaid && time.Now().After(*o.ExpiredAt) {
			orderStatusCode = entity.OrderStatusExpired
		}

		xenditInvoiceUrl := ""
		if o.XenditInvoiceUrl != nil {
			xenditInvoiceUrl = *o.XenditInvoiceUrl
		}
		items = append(items, &order.ListOrderResponseItem{
			Id:                o.Id,
			Number:            o.Number,
			Customer:          o.UserFullName,
			StatusCode:        orderStatusCode,
			Total:             o.Total,
			CreatedAt:         timestamppb.New(o.CreatedAt),
			Products:          products,
			XenditInvoinceUrl: xenditInvoiceUrl,
		})
	}

	return &order.ListOrderResponse{
		Base:       utils.SuccessResponse("Berhasil Mendapatkan List Order"),
		Pagination: metadata,
		Items:      items,
	}, nil
}

func NewOrderService(db *sql.DB, orderRepository repository.IOrderRepository, productRepository repository.IProductRepository) IOrderService {
	return &orderService{
		db:                db,
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}
