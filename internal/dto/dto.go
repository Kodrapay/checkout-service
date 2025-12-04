package dto

type CheckoutSessionRequest struct {
	MerchantID    string `json:"merchant_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
	CustomerEmail string `json:"customer_email,omitempty"`
}

type CheckoutSessionResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type CheckoutPayRequest struct {
	SessionID     string `json:"session_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	TokenID       string `json:"token_id,omitempty"`
}

type CheckoutPayResponse struct {
	TransactionReference string `json:"transaction_reference,omitempty"`
	Status               string `json:"status"`
}
