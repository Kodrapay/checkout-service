package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/google/uuid" // Import uuid

	"github.com/kodra-pay/checkout-service/internal/clients"
	"github.com/kodra-pay/checkout-service/internal/dto"
	"github.com/kodra-pay/checkout-service/internal/models"
)

type CheckoutService struct {
	transactionClient  clients.TransactionClient
	walletLedgerClient clients.WalletLedgerClient
	feeClient          clients.FeeClient
	fraudClient        clients.FraudClient // Add FraudClient
	paymentLinkRepo    PaymentLinkRepository
}

type PaymentLinkRepository interface {
	GetByID(ctx context.Context, id int) (*models.PaymentLink, error)
}

func NewCheckoutService(txClient clients.TransactionClient, wlClient clients.WalletLedgerClient, feeClient clients.FeeClient, fraudClient clients.FraudClient, plRepo PaymentLinkRepository) *CheckoutService {
	return &CheckoutService{
		transactionClient:  txClient,
		walletLedgerClient: wlClient,
		feeClient:          feeClient,
		fraudClient:        fraudClient, // Inject FraudClient
		paymentLinkRepo:    plRepo,
	}
}

func (s *CheckoutService) CreateSession(_ context.Context, req dto.CheckoutSessionRequest) dto.CheckoutSessionResponse {
	// In a real scenario, this would store session details and return a unique ID.
	// For now, it's a placeholder. Since we are using int IDs, we can use a simple counter
	// or, more realistically, this ID would come from a database insertion.
	// For this placeholder, we'll return a static ID or 0.
	return dto.CheckoutSessionResponse{
		ID:       1, // Placeholder for an auto-generated int ID
		Status:   "pending",
		Amount:   req.Amount,
		Currency: req.Currency,
	}
}

func (s *CheckoutService) GetSession(_ context.Context, id int) dto.CheckoutSessionResponse {
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
	amount := req.Amount // currency units (e.g., NGN)
	currency := req.Currency
	description := req.Description
	customerIDStr := strconv.Itoa(req.CustomerID) // Convert CustomerID to string for fraud service

	// If payment link ID is provided, fetch payment link details
	if req.PaymentLinkID != 0 {
		paymentLink, err := s.paymentLinkRepo.GetByID(ctx, req.PaymentLinkID)
		if err != nil {
			return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("failed to get payment link: %w", err)
		}

		// Verify payment link signature to detect tampering
		if !verifyPaymentLinkSignature(paymentLink) {
			return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("payment link signature verification failed - possible tampering detected")
		}

		// Use payment link values where appropriate
		merchantID = paymentLink.MerchantID
		currency = paymentLink.Currency

		// For open links, honor the client-provided amount when present; fall back to link amount only if none was supplied.
		if paymentLink.Mode == "fixed" {
			if paymentLink.Amount != nil {
				amount = float64(*paymentLink.Amount) / 100
			}
		} else {
			if amount == 0 && paymentLink.Amount != nil {
				amount = float64(*paymentLink.Amount) / 100
			}
		}

		if paymentLink.Description != "" {
			description = paymentLink.Description
		}
	}

	// Validate required fields
	if merchantID == 0 || amount <= 0 || currency == "" {
		return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("merchant_id, amount, and currency are required")
	}

	customerID := req.CustomerID
	if customerID == 0 {
		// If customerID is 0, we can potentially use customerEmail for wallet lookup
		// However, for consistency with int IDs, we'll assume 0 means no specific ID.
	}
	if description == "" {
		description = "Payment Link Transaction"
	}

	// Generate a robust transaction reference if not provided or just "2"
	transactionReference := req.Reference
	if transactionReference == "" || transactionReference == strconv.Itoa(req.PaymentLinkID) { // If it's empty or just the payment link ID
		if req.PaymentLinkID != 0 {
			transactionReference = fmt.Sprintf("PL_%d_%s", req.PaymentLinkID, uuid.New().String()[:8])
		} else {
			transactionReference = fmt.Sprintf("TXN_%s", uuid.New().String())
		}
	}

	// === FRAUD CHECK ===
	fraudReq := dto.FraudCheckRequest{
		TransactionReference: transactionReference, // Use the generated/prefixed reference
		Amount:               amount,               // already in currency units
		Currency:             currency,
		CustomerID:           customerIDStr,
		MerchantID:           strconv.Itoa(merchantID),
		Origin:               req.Origin, // Use req.Origin which is passed from handler
		PaymentMethod:        req.PaymentMethod,
		// CustomData:           nil, // Add any custom data if needed
	}

	fraudDecision, err := s.fraudClient.CheckTransaction(ctx, fraudReq)
	if err != nil {
		return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("fraud check failed: %w", err)
	}

	if fraudDecision.Decision == "deny" {
		// Use req.Reference (string) here
		return dto.CheckoutPayResponse{Status: "denied_by_fraud", TransactionReference: transactionReference}, fmt.Errorf("transaction denied by fraud rules: %v", fraudDecision.Reasons)
	}
	// === END FRAUD CHECK ===

	// 1. Quote fees (best-effort; fall back to zero on error)
	var feeAmount float64
	if s.feeClient != nil {
		quote, err := s.feeClient.Quote(ctx, dto.FeeQuoteRequest{
			Amount:   float64(amount),
			Currency: currency,
			Channel:  req.PaymentMethod,
		})
		if err != nil {
			fmt.Printf("Warning: fee quote failed: %v\n", err)
		} else {
			feeAmount = quote.TotalFee
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
		Status:        "successful",         // Assuming immediate success for now
		Reference:     transactionReference, // Use the generated/prefixed reference
	}

	// If fraud decision was "flag", update transaction status accordingly
	if fraudDecision.Decision == "flag" {
		transactionReq.Status = "pending_review" // or "flagged"
	}

	txResp, err := s.transactionClient.CreateTransaction(ctx, transactionReq)
	if err != nil {
		return dto.CheckoutPayResponse{Status: "failed"}, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 3. Update Wallet Balance in Wallet-Ledger Service (optional for payment links)
	// This is primarily for customer wallet management, not required for payment processing
	// We only attempt wallet operations if a valid customer ID is provided.
	if customerID != 0 {
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
				fmt.Printf("Warning: failed to get or create wallet for customer %d: %v\n", customerID, createErr)
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

		netCreditKobo := int64(math.Round(netCredit * 100))

		updateBalanceReq := dto.UpdateBalanceRequest{
			Amount:      netCreditKobo,
			Reference:   txResp.Reference, // Link to the transaction (string)
			Description: fmt.Sprintf("Credit for transaction %s (fee: %.2f)", txResp.Reference, feeAmount),
			Type:        "credit",
		}

		_, err = s.walletLedgerClient.UpdateWalletBalance(ctx, wallet.ID, updateBalanceReq)
		if err != nil {
			// Log the error but don't fail the transaction
			fmt.Printf("Warning: failed to update wallet balance: %v\n", err)
		}
	} else {
		fmt.Printf("Info: Skipping wallet operations for customer 0 as no valid customer ID was provided.\n")
	}

	return dto.CheckoutPayResponse{
		TransactionReference: txResp.Reference,
		Status:               "paid",
	}, nil
}

// verifyPaymentLinkSignature verifies that a payment link's parameters haven't been tampered with
func verifyPaymentLinkSignature(link *models.PaymentLink) bool {
	if link.Signature == "" {
		// Old links without signature - allow them but log warning
		fmt.Printf("Warning: Payment link %d has no signature\n", link.ID)
		return true
	}

	secret := os.Getenv("PAYMENT_LINK_SECRET")
	if secret == "" {
		secret = "kodrapay-default-secret-change-in-production"
	}

	// Create canonical string from link parameters (same as generation)
	amountStr := "null"
	if link.Amount != nil {
		amountStr = strconv.FormatInt(*link.Amount, 10)
	}

	data := fmt.Sprintf("%d|%s|%s|%s|%s",
		link.MerchantID,
		link.Mode,
		amountStr,
		link.Currency,
		link.Description,
	)

	// Generate expected HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(link.Signature), []byte(expectedSignature))
}
