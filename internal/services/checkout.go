package services

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/kodra-pay/checkout-service/internal/clients"
	"github.com/kodra-pay/checkout-service/internal/dto"
	"github.com/kodra-pay/checkout-service/internal/models"
)

type CheckoutService struct {
	transactionClient  clients.TransactionClient
	walletLedgerClient clients.WalletLedgerClient
	feeClient          clients.FeeClient
	paymentLinkRepo    PaymentLinkRepository
}

type PaymentLinkRepository interface {
	GetByID(ctx context.Context, id string) (*models.PaymentLink, error)
}

func NewCheckoutService(txClient clients.TransactionClient, wlClient clients.WalletLedgerClient, feeClient clients.FeeClient, plRepo PaymentLinkRepository) *CheckoutService {
	return &CheckoutService{
		transactionClient:  txClient,
		walletLedgerClient: wlClient,
		feeClient:          feeClient,
		paymentLinkRepo:    plRepo,
	}
}

func (s *CheckoutService) CreateSession(_ context.Context, req dto.CheckoutSessionRequest) dto.CheckoutSessionResponse {
	// In a real scenario, this would store session details and return a unique ID.
	// For now, it's a placeholder.
	return dto.CheckoutSessionResponse{
		ID:       "chk_" + uuid.NewString(),
		Status:   "pending",
		Amount:   req.Amount,
		Currency: req.Currency,
	}
}

func (s *CheckoutService) GetSession(_ context.Context, id string) dto.CheckoutSessionResponse {
	// In a real scenario, this would retrieve session details by ID.
	// For now, it's a placeholder.
	return dto.CheckoutSessionResponse{
		ID:     id,
		Status: "pending",
	}
}

func (s *CheckoutService) Pay(ctx context.Context, req dto.CheckoutPayRequest) (dto.CheckoutPayResponse, error) {
	// Simulate payment processing (e.g., call a payment gateway)
	// For this exercise, we assume payment is successful.

	merchantID := req.MerchantID
	amount := req.Amount
	currency := req.Currency
	description := req.Description

	// If payment link ID is provided, fetch payment link details
	if req.PaymentLinkID != "" {
		paymentLink, err := s.paymentLinkRepo.GetByID(ctx, req.PaymentLinkID)
		if err != nil {
			return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("failed to get payment link: %w", err)
		}

		// Use payment link values if not provided in request
		merchantID = paymentLink.MerchantID
		currency = paymentLink.Currency
		if paymentLink.Amount != nil {
			amount = *paymentLink.Amount
		}
		if paymentLink.Description != "" {
			description = paymentLink.Description
		}
	}

	// Validate required fields
	if merchantID == "" || amount <= 0 || currency == "" {
		return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("merchant_id, amount, and currency are required")
	}

	customerID := req.CustomerID
	if customerID == "" {
		customerID = req.CustomerEmail
	}
	if description == "" {
		description = "Payment Link Transaction"
	}

	// 1. Quote fees (best-effort; fall back to zero on error)
	var feeAmount int64
	if s.feeClient != nil {
		quote, err := s.feeClient.Quote(ctx, dto.FeeQuoteRequest{
			Amount:   float64(amount),
			Currency: currency,
			Channel:  req.PaymentMethod,
		})
		if err != nil {
			fmt.Printf("Warning: fee quote failed: %v\n", err)
		} else {
			feeAmount = int64(math.Round(quote.TotalFee))
		}
	}

	// 2. Create Transaction in Transaction Service (gross amount)
	transactionReq := dto.TransactionCreateRequest{
		MerchantID:    merchantID,
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
		CustomerID:    customerID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: req.PaymentMethod,
		Description:   description,
		Status:        "successful", // Assuming immediate success for now
		Reference:     req.Reference,
	}

	txResp, err := s.transactionClient.CreateTransaction(ctx, transactionReq)
	if err != nil {
		return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 3. Update Wallet Balance in Wallet-Ledger Service (optional for payment links)
	// This is primarily for customer wallet management, not required for payment processing
	// First, try to get the customer's wallet
	wallet, err := s.walletLedgerClient.GetWalletByUserIDAndCurrency(ctx, customerID, currency)
	if err != nil {
		// If wallet not found, try to create one
		createWalletReq := dto.CreateWalletRequest{
			UserID:   customerID,
			Currency: currency,
		}
		newWallet, createErr := s.walletLedgerClient.CreateWallet(ctx, createWalletReq)
		if createErr != nil {
			// Wallet ledger service unavailable - log but don't fail the transaction
			fmt.Printf("Warning: failed to get or create wallet for customer %s: %v\n", customerID, createErr)
			// Transaction was already created successfully, return success
			return dto.CheckoutPayResponse{
				TransactionReference: txResp.Reference,
				Status:               "paid",
			}, nil
		}
		wallet = newWallet
	}

	// Credit the wallet
	netCredit := amount - feeAmount
	if netCredit < 0 {
		netCredit = 0
	}

	updateBalanceReq := dto.UpdateBalanceRequest{
		Amount:      netCredit,
		Reference:   txResp.Reference, // Link to the transaction
		Description: fmt.Sprintf("Credit for transaction %s (fee: %d)", txResp.Reference, feeAmount),
		Type:        "credit",
	}

	_, err = s.walletLedgerClient.UpdateWalletBalance(ctx, wallet.ID, updateBalanceReq)
	if err != nil {
		// Log the error but don't fail the transaction
		fmt.Printf("Warning: failed to update wallet balance: %v\n", err)
	}

	return dto.CheckoutPayResponse{
		TransactionReference: txResp.Reference,
		Status:               "paid",
	}, nil
}
