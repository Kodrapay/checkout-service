package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kodra-pay/checkout-service/internal/dto"
)

// FraudClient defines the interface for interacting with the Fraud Service.
type FraudClient interface {
	CheckTransaction(ctx context.Context, req dto.FraudCheckRequest) (dto.FraudDecision, error)
}

// HTTPFraudClient is an HTTP implementation of the FraudClient interface.
type HTTPFraudClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewHTTPFraudClient creates a new HTTPFraudClient.
func NewHTTPFraudClient(baseURL, apiKey string) *HTTPFraudClient {
	return &HTTPFraudClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

// CheckTransaction calls the fraud service to evaluate a transaction.
func (c *HTTPFraudClient) CheckTransaction(ctx context.Context, req dto.FraudCheckRequest) (dto.FraudDecision, error) {
	url := fmt.Sprintf("%s/fraud/check-transaction", c.baseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return dto.FraudDecision{}, fmt.Errorf("failed to marshal fraud check request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return dto.FraudDecision{}, fmt.Errorf("failed to create http request for fraud check: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("X-API-Key", c.apiKey)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return dto.FraudDecision{}, fmt.Errorf("failed to call fraud service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return dto.FraudDecision{}, fmt.Errorf("fraud service returned non-ok status: %d, body: %s", resp.StatusCode, respBody)
	}

	var decision dto.FraudDecision
	if err := json.NewDecoder(resp.Body).Decode(&decision); err != nil {
		return dto.FraudDecision{}, fmt.Errorf("failed to decode fraud decision response: %w", err)
	}

	return decision, nil
}
