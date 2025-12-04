package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/kodra-pay/checkout-service/internal/dto"
)

type CheckoutService struct{}

func NewCheckoutService() *CheckoutService {
	return &CheckoutService{}
}

func (s *CheckoutService) CreateSession(_ context.Context, req dto.CheckoutSessionRequest) dto.CheckoutSessionResponse {
	return dto.CheckoutSessionResponse{
		ID:       "chk_" + uuid.NewString(),
		Status:   "pending",
		Amount:   req.Amount,
		Currency: req.Currency,
	}
}

func (s *CheckoutService) GetSession(_ context.Context, id string) dto.CheckoutSessionResponse {
	return dto.CheckoutSessionResponse{
		ID:     id,
		Status: "pending",
	}
}

func (s *CheckoutService) Pay(_ context.Context, req dto.CheckoutPayRequest) dto.CheckoutPayResponse {
	return dto.CheckoutPayResponse{
		Status: "paid",
	}
}
