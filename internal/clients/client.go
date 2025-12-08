package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kodra-pay/checkout-service/internal/dto"
)

// TransactionClient defines the interface for interacting with the Transaction Service
type TransactionClient interface {
	CreateTransaction(ctx context.Context, req dto.TransactionCreateRequest) (*dto.TransactionResponse, error)
}

// HTTPTransactionClient implements TransactionClient using HTTP
type HTTPTransactionClient struct {
	baseURL string
	client  *http.Client
}

// NewHTTPTransactionClient creates a new HTTPTransactionClient
func NewHTTPTransactionClient(baseURL string) *HTTPTransactionClient {
	return &HTTPTransactionClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPTransactionClient) CreateTransaction(ctx context.Context, req dto.TransactionCreateRequest) (*dto.TransactionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction create request: %w", err)
	}

	url := fmt.Sprintf("%s/transactions", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request to transaction service: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("transaction service returned non-201 status: %d", httpResp.StatusCode)
	}

	var resp dto.TransactionResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode transaction response: %w", err)
	}
	return &resp, nil
}

// WalletLedgerClient defines the interface for interacting with the Wallet-Ledger Service
type WalletLedgerClient interface {
	GetWalletByUserIDAndCurrency(ctx context.Context, userID int, currency string) (*dto.WalletResponse, error)
	CreateWallet(ctx context.Context, req dto.CreateWalletRequest) (*dto.WalletResponse, error)
	UpdateWalletBalance(ctx context.Context, walletID int, req dto.UpdateBalanceRequest) (*dto.WalletResponse, error)
}

// FeeClient defines interface for fee-service
type FeeClient interface {
	Quote(ctx context.Context, req dto.FeeQuoteRequest) (*dto.FeeQuoteResponse, error)
}

// HTTPWalletLedgerClient implements WalletLedgerClient using HTTP
type HTTPWalletLedgerClient struct {
	baseURL string
	client  *http.Client
}

// NewHTTPWalletLedgerClient creates a new HTTPWalletLedgerClient
func NewHTTPWalletLedgerClient(baseURL string) *HTTPWalletLedgerClient {
	return &HTTPWalletLedgerClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPWalletLedgerClient) GetWalletByUserIDAndCurrency(ctx context.Context, userID int, currency string) (*dto.WalletResponse, error) {
	url := fmt.Sprintf("%s/wallets?user_id=%d&currency=%s", c.baseURL, userID, currency)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request to wallet-ledger service: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet-ledger service returned non-200 status: %d", httpResp.StatusCode)
	}

	var resp dto.WalletResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode wallet response: %w", err)
	}
	return &resp, nil
}

func (c *HTTPWalletLedgerClient) CreateWallet(ctx context.Context, req dto.CreateWalletRequest) (*dto.WalletResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal wallet create request: %w", err)
	}

	url := fmt.Sprintf("%s/wallets", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request to wallet-ledger service: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("wallet-ledger service returned non-201 status for create wallet: %d", httpResp.StatusCode)
	}

	var resp dto.WalletResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode wallet response: %w", err)
	}
	return &resp, nil
}

func (c *HTTPWalletLedgerClient) UpdateWalletBalance(ctx context.Context, walletID int, req dto.UpdateBalanceRequest) (*dto.WalletResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update balance request: %w", err)
	}

	url := fmt.Sprintf("%s/wallets/%d/update-balance", c.baseURL, walletID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request to wallet-ledger service: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wallet-ledger service returned non-200 status for update balance: %d", httpResp.StatusCode)
	}

	var resp dto.WalletResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode wallet response: %w", err)
	}
	return &resp, nil
}

// HTTPFeeClient implements FeeClient
type HTTPFeeClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPFeeClient(baseURL string) *HTTPFeeClient {
	return &HTTPFeeClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPFeeClient) Quote(ctx context.Context, req dto.FeeQuoteRequest) (*dto.FeeQuoteResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal fee quote request: %w", err)
	}

	url := fmt.Sprintf("%s/fees/quote", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create fee quote request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("fee service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("fee service returned status: %d", resp.StatusCode)
	}

	var quote dto.FeeQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, fmt.Errorf("decode fee quote: %w", err)
	}
	return &quote, nil
}
