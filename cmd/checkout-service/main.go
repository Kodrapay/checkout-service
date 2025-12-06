package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/checkout-service/internal/clients" // Import clients package
	"github.com/kodra-pay/checkout-service/internal/config"
	"github.com/kodra-pay/checkout-service/internal/middleware"
	"github.com/kodra-pay/checkout-service/internal/routes"
)

func main() {
	cfg := config.Load("checkout-service", "7005")

	// Initialize HTTP clients for other services
	transactionClient := clients.NewHTTPTransactionClient(cfg.TransactionServiceURL)
	walletLedgerClient := clients.NewHTTPWalletLedgerClient(cfg.WalletLedgerServiceURL)
	feeClient := clients.NewHTTPFeeClient(cfg.FeeServiceURL)

	app := fiber.New()
	app.Use(middleware.RequestID())

	// Pass the clients to the routes registration
	routes.Register(app, cfg.ServiceName, transactionClient, walletLedgerClient, feeClient)

	log.Printf("%s listening on :%s", cfg.ServiceName, cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
