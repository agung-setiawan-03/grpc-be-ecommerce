package service

import (
	"context"
	"errors"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/dto"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/repository"
)

type IWebhookService interface {
	ReceiveInvoice(ctx context.Context, request *dto.XenditInvoinceRequest) error
}

type webhookService struct {
	orderRepository repository.IOrderRepository
}

func (ws *webhookService) ReceiveInvoice(ctx context.Context, request *dto.XenditInvoinceRequest) error {
	// find orer di db
	orderEntity, err := ws.orderRepository.GetOrderById(ctx, request.ExternalID)
	if err != nil {
		return err
	}

	if orderEntity == nil {
		return errors.New("order tidak ditemukan")
	}

	// gangi / update status order
	now := time.Now()
	updatedBy := "System Xendit"
	orderEntity.OrderStatusCode = entity.OrderStatusCodePaid
	orderEntity.UpdatedAt = &now
	orderEntity.UpdatedBy = &updatedBy
	orderEntity.XenditPaidAt = &now
	orderEntity.XenditPaymentChannel = &request.PaymentChannel
	orderEntity.XenditPaymentMethod = &request.PaymentMethod

	// update ke db

	err = ws.orderRepository.UpdateOrder(ctx, orderEntity)
	if err != nil {
		return err
	}

	return nil
}

func NewWebhookService(orderRepository repository.IOrderRepository) IWebhookService {
	return &webhookService{
		orderRepository: orderRepository,
	}
}
