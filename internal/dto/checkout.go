package dto

type CheckoutSessionRequest struct {
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	MerchantID    string `json:"merchant_id"`
	Description   string `json:"description"`
	CustomerEmail string `json:"customer_email"`
}

type CheckoutSessionResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type CheckoutPayRequest struct {
	SessionID string `json:"session_id"`
	TokenID   string `json:"token_id"`
}

type CheckoutPayResponse struct {
	TransactionReference string `json:"transaction_reference"`
	Status               string `json:"status"`
}
