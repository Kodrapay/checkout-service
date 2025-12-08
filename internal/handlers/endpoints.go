package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodra-pay/checkout-service/internal/dto"
	"github.com/kodra-pay/checkout-service/internal/services"
)

type CheckoutHandler struct {
	svc *services.CheckoutService
}

func NewCheckoutHandler(svc *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{svc: svc}
}

func (h *CheckoutHandler) CreateSession(c *fiber.Ctx) error {
	var req dto.CheckoutSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	return c.JSON(h.svc.CreateSession(c.Context(), req))
}

func (h *CheckoutHandler) GetSession(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id") // Use c.ParamsInt
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid session ID")
	}
	return c.JSON(h.svc.GetSession(c.Context(), id))
}

func (h *CheckoutHandler) Pay(c *fiber.Ctx) error {
	var req dto.CheckoutPayRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	req.Origin = c.IP() // Set the client IP from Fiber context
	resp, err := h.svc.Pay(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(resp)
}

type PaymentLinkHandler struct {
	svc *services.PaymentLinkService
}

func NewPaymentLinkHandler(svc *services.PaymentLinkService) *PaymentLinkHandler {
	return &PaymentLinkHandler{svc: svc}
}

func (h *PaymentLinkHandler) Create(c *fiber.Ctx) error {
	var req dto.PaymentLinkCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	resp, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(resp)
}

func (h *PaymentLinkHandler) List(c *fiber.Ctx) error {
	merchantID := c.QueryInt("merchant_id", 0) // Use c.QueryInt for query parameters
	// If merchantID is mandatory, you might add an error check here:
	// if merchantID == 0 {
	//    return fiber.NewError(fiber.StatusBadRequest, "merchant_id query parameter is required")
	// }
	return c.JSON(h.svc.ListByMerchant(c.Context(), merchantID))
}

func (h *PaymentLinkHandler) Get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid payment link id")
	}
	pl, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Payment link not found")
	}
	return c.JSON(pl)
}
