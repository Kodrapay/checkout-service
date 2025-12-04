package dto

type CheckoutSessionRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	Description string `json:"description"`
}

type CheckoutSessionResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Amount  int64  `json:"amount"`
	Currency string `json:"currency"`
}

type CheckoutPayRequest struct {
	SessionID     string `json:"session_id"`
	PaymentMethod string `json:"payment_method"`
}

type CheckoutPayResponse struct {
	Status string `json:"status"`
}
