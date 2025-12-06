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
	id := c.Params("id")
	return c.JSON(h.svc.GetSession(c.Context(), id))
}

func (h *CheckoutHandler) Pay(c *fiber.Ctx) error {
	var req dto.CheckoutPayRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
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
	return c.JSON(h.svc.Create(c.Context(), req))
}

func (h *PaymentLinkHandler) List(c *fiber.Ctx) error {
	merchantID := c.Query("merchant_id")
	return c.JSON(h.svc.ListByMerchant(c.Context(), merchantID))
}
