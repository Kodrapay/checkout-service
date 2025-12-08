package dto

// FraudCheckRequest represents the request body for the fraud service's CheckTransaction endpoint.
type FraudCheckRequest struct {
	TransactionReference string                 `json:"transaction_reference,omitempty"`
	Amount               float64                `json:"amount"`
	Currency             string                 `json:"currency"`
	CustomerID           string                 `json:"customer_id,omitempty"`
	MerchantID           string                 `json:"merchant_id,omitempty"`
	Origin               string                 `json:"origin,omitempty"` // e.g., IP address
	PaymentMethod        string                 `json:"payment_method,omitempty"`
	CustomData           map[string]interface{} `json:"custom_data,omitempty"`
}

// FraudDecision represents the outcome of a fraud evaluation, copied from fraud-service.
type FraudDecision struct {
	OverallScore float64  `json:"overall_score"`
	Decision     string   `json:"decision"` // "approve", "flag", "deny"
	Reasons      []string `json:"reasons"`
}