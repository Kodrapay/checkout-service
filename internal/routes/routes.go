package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/checkout-service/internal/config"
	"github.com/kodra-pay/checkout-service/internal/handlers"
	"github.com/kodra-pay/checkout-service/internal/repositories"
	"github.com/kodra-pay/checkout-service/internal/services"
)

func Register(app *fiber.App, cfg config.Config, repo *repositories.CheckoutRepository) {
	health := handlers.NewHealthHandler(cfg.ServiceName)
	health.Register(app)

	svc := services.NewCheckoutService(repo)
	h := handlers.NewCheckoutHandler(svc)
	api := app.Group("/checkout")
	api.Post("/sessions", h.CreateSession)
	api.Get("/sessions/:id", h.GetSession)
	api.Post("/pay", h.Pay)
}
