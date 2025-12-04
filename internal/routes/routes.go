package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/checkout-service/internal/handlers"
	"github.com/kodra-pay/checkout-service/internal/config"
	"github.com/kodra-pay/checkout-service/internal/repositories"
	"github.com/kodra-pay/checkout-service/internal/services"
)

func Register(app *fiber.App, serviceName string) {
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	cfg := config.Load(serviceName, "7005")
	repo, err := repositories.NewPaymentLinkRepository(cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}
	plSvc := services.NewPaymentLinkService(repo)
	plHandler := handlers.NewPaymentLinkHandler(plSvc)
	checkoutSvc := services.NewCheckoutService()
	checkoutHandler := handlers.NewCheckoutHandler(checkoutSvc)

	app.Post("/payment-links", plHandler.Create)
	app.Get("/payment-links", plHandler.List)

	app.Post("/checkout/session", checkoutHandler.CreateSession)
	app.Get("/checkout/session/:id", checkoutHandler.GetSession)
	app.Post("/checkout/pay", checkoutHandler.Pay)
}
