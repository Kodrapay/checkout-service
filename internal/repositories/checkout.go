package repositories

import "log"

type CheckoutRepository struct {
    dsn string
}

func NewCheckoutRepository(dsn string) *CheckoutRepository {
    log.Printf("CheckoutRepository using DSN: %s", dsn)
    return &CheckoutRepository{dsn: dsn}
}

// TODO: implement persistence for checkout sessions and payments.
