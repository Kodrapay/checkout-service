package dto

type PaymentLinkCreateRequest struct {
	MerchantID  int     `json:"merchant_id"`
	Mode        string  `json:"mode"` // fixed or open
	Amount      *int64  `json:"amount,omitempty"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	Reference   string  `json:"reference"`
}

type PaymentLinkResponse struct {
	ID          int     `json:"id"`
	MerchantID  int     `json:"merchant_id"`
	Mode        string  `json:"mode"`
	Amount      *int64  `json:"amount,omitempty"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	Reference   string  `json:"reference"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type PaymentLinkListResponse struct {
	Links []PaymentLinkResponse `json:"links"`
}
