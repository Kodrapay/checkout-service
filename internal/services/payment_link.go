package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kodra-pay/checkout-service/internal/dto"
	"github.com/kodra-pay/checkout-service/internal/models"
	"github.com/kodra-pay/checkout-service/internal/repositories"
)

type PaymentLinkService struct {
	repo *repositories.PaymentLinkRepository
}

func NewPaymentLinkService(repo *repositories.PaymentLinkRepository) *PaymentLinkService {
	return &PaymentLinkService{repo: repo}
}

func (s *PaymentLinkService) Create(ctx context.Context, req dto.PaymentLinkCreateRequest) dto.PaymentLinkResponse {
	ref := req.Reference
	if ref == "" {
		ref = "pl_" + uuid.NewString()
	}

	pl := &models.PaymentLink{
		MerchantID:  req.MerchantID,
		Mode:        req.Mode,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		Status:      "active",
		Reference:   ref,
	}
	_ = s.repo.Create(ctx, pl)

	return dto.PaymentLinkResponse{
		ID:          pl.ID,
		MerchantID:  pl.MerchantID,
		Mode:        pl.Mode,
		Amount:      pl.Amount,
		Currency:    pl.Currency,
		Description: pl.Description,
		Reference:   pl.Reference,
		Status:      pl.Status,
		CreatedAt:   pl.CreatedAt.Format(time.RFC3339),
	}
}

func (s *PaymentLinkService) ListByMerchant(ctx context.Context, merchantID string) dto.PaymentLinkListResponse {
	links, _ := s.repo.ListByMerchant(ctx, merchantID, 50)
	resp := dto.PaymentLinkListResponse{}
	for _, pl := range links {
		resp.Links = append(resp.Links, dto.PaymentLinkResponse{
			ID:          pl.ID,
			MerchantID:  pl.MerchantID,
			Mode:        pl.Mode,
			Amount:      pl.Amount,
			Currency:    pl.Currency,
			Description: pl.Description,
			Reference:   pl.Reference,
			Status:      pl.Status,
			CreatedAt:   pl.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp
}
