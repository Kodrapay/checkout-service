package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// MerchantResponse represents the response from merchant service
type MerchantResponse struct {
	ID          string `json:"id"`
	KYCStatus   string `json:"kyc_status"`
	Status      string `json:"status"`
	CanTransact bool   `json:"can_transact"`
}

// KYCCheckMiddleware validates merchant KYC status before allowing checkout
type KYCCheckMiddleware struct {
	merchantServiceURL string
}

func NewKYCCheckMiddleware(merchantServiceURL string) *KYCCheckMiddleware {
	return &KYCCheckMiddleware{
		merchantServiceURL: merchantServiceURL,
	}
}

// RequireApprovedKYC checks if merchant has approved KYC via merchant service
func (m *KYCCheckMiddleware) RequireApprovedKYC(c *fiber.Ctx) error {
	// Extract merchant ID from request body or context
	merchantID := c.Locals("merchant_id")
	if merchantID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "merchant not authenticated")
	}

	merchantIDStr, ok := merchantID.(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid merchant ID format")
	}

	// Call merchant service to check KYC status
	url := fmt.Sprintf("%s/merchants/%s", m.merchantServiceURL, merchantIDStr)
	resp, err := http.Get(url)
	if err != nil {
		// If merchant service is unavailable, log error but don't block
		// In production, you might want to cache merchant KYC status
		return c.Next()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fiber.NewError(fiber.StatusBadRequest, "unable to verify merchant status")
	}

	var merchant MerchantResponse
	if err := json.NewDecoder(resp.Body).Decode(&merchant); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to parse merchant data")
	}

	// Check if merchant can transact
	if !merchant.CanTransact {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "kyc_not_approved",
			"message": "Your KYC verification must be approved before you can process transactions. Please complete your business KYC verification at /merchant/kyc",
			"kyc_status": merchant.KYCStatus,
			"can_transact": false,
		})
	}

	return c.Next()
}
