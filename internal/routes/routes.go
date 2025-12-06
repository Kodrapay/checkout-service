package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/kodra-pay/checkout-service/internal/clients"
	"github.com/kodra-pay/checkout-service/internal/config"
	"github.com/kodra-pay/checkout-service/internal/handlers"
	"github.com/kodra-pay/checkout-service/internal/repositories"
	"github.com/kodra-pay/checkout-service/internal/services"
)

func Register(app *fiber.App, serviceName string, txClient clients.TransactionClient, wlClient clients.WalletLedgerClient, feeClient clients.FeeClient) {
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	// Note: PaymentLinkRepository and PaymentLinkService setup might need a dedicated DB connection or be refactored
	// to use clients if they interact with other services. For now, assuming they use local DB.
	cfg := config.Load(serviceName, "7005") // Still needed for PaymentLinkRepository's DSN
	repo, err := repositories.NewPaymentLinkRepository(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to initialize PaymentLinkRepository: %v", err) // Use log.Fatalf instead of panic
	}
	plSvc := services.NewPaymentLinkService(repo)
	plHandler := handlers.NewPaymentLinkHandler(plSvc)

	checkoutSvc := services.NewCheckoutService(txClient, wlClient, feeClient, repo) // Pass the clients and payment link repo here
	checkoutHandler := handlers.NewCheckoutHandler(checkoutSvc)

	app.Post("/payment-links", plHandler.Create)
	app.Get("/payment-links", plHandler.List)

	app.Post("/checkout/session", checkoutHandler.CreateSession)
	app.Get("/checkout/session/:id", checkoutHandler.GetSession)
	app.Post("/checkout/pay", checkoutHandler.Pay)
}
