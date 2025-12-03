package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/kodra-pay/checkout-service/internal/dto"
	"github.com/kodra-pay/checkout-service/internal/repositories"
)

type CheckoutService struct {
	repo *repositories.CheckoutRepository
}

func NewCheckoutService(repo *repositories.CheckoutRepository) *CheckoutService {
	return &CheckoutService{repo: repo}
}

func (s *CheckoutService) CreateSession(_ context.Context, req dto.CheckoutSessionRequest) dto.CheckoutSessionResponse {
	return dto.CheckoutSessionResponse{
		ID:       "chk_" + uuid.NewString(),
		Status:   "created",
		Amount:   req.Amount,
		Currency: req.Currency,
	}
}

func (s *CheckoutService) GetSession(_ context.Context, id string) dto.CheckoutSessionResponse {
	return dto.CheckoutSessionResponse{
		ID:       id,
		Status:   "created",
		Amount:   0,
		Currency: "NGN",
	}
}

func (s *CheckoutService) Pay(_ context.Context, req dto.CheckoutPayRequest) dto.CheckoutPayResponse {
	return dto.CheckoutPayResponse{
		TransactionReference: "txn_" + uuid.NewString(),
		Status:               "processing",
	}
}
