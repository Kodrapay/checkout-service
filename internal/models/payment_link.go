package models

import "time"

type PaymentLink struct {
	ID          int        `json:"id"`
	MerchantID  int        `json:"merchant_id"`
	Mode        string     `json:"mode"`             // fixed or open
	Amount      *int64     `json:"amount,omitempty"` // stored in kobo
	Currency    string     `json:"currency"`
	Description string     `json:"description"`
	Reference   string     `json:"reference"`
	Status      string     `json:"status"`
	Signature   string     `json:"signature,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
