package config

import (
	"os"
	"strings"
)

type Config struct {
	ServiceName            string
	Port                   string
	PostgresDSN            string
	RedisAddr              string
	TransactionServiceURL  string
	WalletLedgerServiceURL string
	FeeServiceURL          string
	FraudServiceURL        string // New field for Fraud Service URL
	FraudServiceAPIKey     string // New field for Fraud Service API Key
}

func Load(serviceName, defaultPort string) Config {
	dsn := getEnv("POSTGRES_URL", "postgres://kodrapay:kodrapay_password@postgres:5432/kodrapay?sslmode=disable")
	if !strings.Contains(strings.ToLower(dsn), "sslmode=") {
		if strings.Contains(dsn, "?") {
			dsn += "&sslmode=disable"
		} else {
			dsn += "?sslmode=disable"
		}
	}

	return Config{
		ServiceName:            serviceName,
		Port:                   getEnv("PORT", defaultPort),
		PostgresDSN:            dsn,
		RedisAddr:              getEnv("REDIS_ADDR", "redis:6379"),
		TransactionServiceURL:  getEnv("TRANSACTION_SERVICE_URL", "http://transaction-service:7004/api/v1"),     // Align with docker-compose port
		WalletLedgerServiceURL: getEnv("WALLET_LEDGER_SERVICE_URL", "http://wallet-ledger-service:7007/api/v1"), // Align with docker-compose port
		FeeServiceURL:          getEnv("FEE_SERVICE_URL", "http://fee-service:7017"),                            // Fee service base
		FraudServiceURL:        getEnv("FRAUD_SERVICE_URL", "http://fraud-service:7012"),                         // Fraud service base
		FraudServiceAPIKey:     getEnv("FRAUD_SERVICE_API_KEY", "my-secret-api-key"),                          // Fraud service API key
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
