package handler

import (
	"context"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/service"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/order"
)

type orderHandler struct {
	order.UnimplementedOrderServiceServer

	orderService service.IOrderService
}

func (oh *orderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &order.CreateOrderResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := oh.orderService.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (oh *orderHandler) ListOrderAdmin(ctx context.Context, req *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &order.ListOrderAdminResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := oh.orderService.ListOrderAdmin(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (oh *orderHandler) ListOrder(ctx context.Context, req *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &order.ListOrderResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := oh.orderService.ListOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (oh *orderHandler) DetailOrder(ctx context.Context, req *order.DetailOrderRequest) (*order.DetailOrderResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &order.DetailOrderResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := oh.orderService.DetailOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (oh *orderHandler) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	validationErorrs, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErorrs != nil {
		return &order.UpdateOrderStatusResponse{
			Base: utils.ValidationErrorResponse(validationErorrs),
		}, nil
	}

	resp, err := oh.orderService.UpdateOrderStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewOrderHandler(orderService service.IOrderService) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}
