package services

import (
	"context"
	"fmt"
	"math"
	"time"

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

func (s *PaymentLinkService) Create(ctx context.Context, req dto.PaymentLinkCreateRequest) (dto.PaymentLinkResponse, error) {
	var amountKobo *int64
	if req.Amount != nil {
		val := int64(math.Round(*req.Amount * 100))
		amountKobo = &val
	}

	pl := &models.PaymentLink{
		MerchantID:  req.MerchantID,
		Mode:        req.Mode,
		Amount:      amountKobo,
		Currency:    req.Currency,
		Description: req.Description,
		Status:      "active",
	}
	if err := s.repo.Create(ctx, pl); err != nil {
		return dto.PaymentLinkResponse{}, fmt.Errorf("failed to create payment link in repository: %w", err)
	}

	return dto.PaymentLinkResponse{
		ID:          pl.ID,
		MerchantID:  pl.MerchantID,
		Mode:        pl.Mode,
		Amount:      toCurrency(pl.Amount),
		Currency:    pl.Currency,
		Description: pl.Description,
		Status:      pl.Status,
		CreatedAt:   pl.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *PaymentLinkService) ListByMerchant(ctx context.Context, merchantID int) dto.PaymentLinkListResponse { // merchantID changed to int
	links, _ := s.repo.ListByMerchant(ctx, merchantID, 50)
	resp := dto.PaymentLinkListResponse{}
	for _, pl := range links {
		resp.Links = append(resp.Links, dto.PaymentLinkResponse{
			ID:          pl.ID,
			MerchantID:  pl.MerchantID,
			Mode:        pl.Mode,
			Amount:      toCurrency(pl.Amount),
			Currency:    pl.Currency,
			Description: pl.Description,
			Status:      pl.Status,
			CreatedAt:   pl.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp
}

func (s *PaymentLinkService) Get(ctx context.Context, id int) (*dto.PaymentLinkResponse, error) {
	pl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := &dto.PaymentLinkResponse{
		ID:          pl.ID,
		MerchantID:  pl.MerchantID,
		Mode:        pl.Mode,
		Amount:      toCurrency(pl.Amount),
		Currency:    pl.Currency,
		Description: pl.Description,
		Status:      pl.Status,
		CreatedAt:   pl.CreatedAt.Format(time.RFC3339),
	}
	return resp, nil
}

func toCurrency(amountKobo *int64) *float64 {
	if amountKobo == nil {
		return nil
	}
	val := float64(*amountKobo) / 100
	return &val
}
