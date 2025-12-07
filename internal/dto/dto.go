package dto

import "time"

type CheckoutSessionRequest struct {
	MerchantID    int    `json:"merchant_id,omitempty"`
	Amount        int64  `json:"amount,omitempty"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerID    int    `json:"customer_id,omitempty"`
}

type CheckoutSessionResponse struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type CheckoutPayRequest struct {
	SessionID     int    `json:"session_id,omitempty"`
	PaymentLinkID int    `json:"payment_link_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	TokenID       string `json:"token_id,omitempty"`
	MerchantID    int    `json:"merchant_id,omitempty"`
	Amount        int64  `json:"amount,omitempty"`
	Currency      string `json:"currency,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerID    int    `json:"customer_id,omitempty"`
	CustomerName  string `json:"customer_name,omitempty"`
	Description   string `json:"description,omitempty"`
	Reference     string `json:"reference,omitempty"`
}

type CheckoutPayResponse struct {
	TransactionReference int    `json:"transaction_reference,omitempty"`
	Status               string `json:"status"`
}

// TransactionCreateRequest DTO for creating a new transaction in transaction-service
type TransactionCreateRequest struct {
	Reference     int    `json:"reference,omitempty"`
	MerchantID    int    `json:"merchant_id"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerName  string `json:"customer_name,omitempty"`
	CustomerID    int    `json:"customer_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Description   string `json:"description,omitempty"`
	Status        string `json:"status,omitempty"` // status should be handled by transaction service
}

// TransactionResponse DTO for transaction-service response
type TransactionResponse struct {
	ID            int       `json:"id"`
	Reference     int       `json:"reference"`
	MerchantID    int       `json:"merchant_id"`
	CustomerEmail string    `json:"customer_email"`
	CustomerName  string    `json:"customer_name,omitempty"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateWalletRequest DTO for creating a new wallet in wallet-ledger-service
type CreateWalletRequest struct {
	UserID   int    `json:"user_id"`
	Currency string `json:"currency"`
}

// UpdateBalanceRequest DTO for updating a wallet's balance in wallet-ledger-service
type UpdateBalanceRequest struct {
	Amount      int64  `json:"amount"`
	Reference   int    `json:"reference"`
	Description string `json:"description"`
	Type        string `json:"type"` // "credit" or "debit"
}

// WalletResponse DTO for returning wallet information from wallet-ledger-service
type WalletResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Currency  string    `json:"currency"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Fee quote DTOs
type FeeQuoteRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Channel  string  `json:"channel"`
}

type FeeQuoteResponse struct {
	TotalFee   float64 `json:"total_fee"`
	BaseAmount float64 `json:"base_amount"`
	Currency   string  `json:"currency"`
	Rate       float64 `json:"rate"`
	Flat       float64 `json:"flat"`
	Capped     bool    `json:"capped"`
	Cap        float64 `json:"cap,omitempty"`
	Channel    string  `json:"channel,omitempty"`
}
