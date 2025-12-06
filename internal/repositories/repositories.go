package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/kodra-pay/checkout-service/internal/models"
)

type PaymentLinkRepository struct {
	db *sql.DB
}

func NewPaymentLinkRepository(dsn string) (*PaymentLinkRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &PaymentLinkRepository{db: db}, nil
}

func (r *PaymentLinkRepository) Create(ctx context.Context, pl *models.PaymentLink) error {
	query := `
		INSERT INTO payment_links (merchant_id, mode, amount, currency, description, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		pl.MerchantID, pl.Mode, pl.Amount, pl.Currency, pl.Description, pl.Status,
	).Scan(&pl.ID, &pl.CreatedAt, &pl.UpdatedAt)
}

func (r *PaymentLinkRepository) GetByID(ctx context.Context, id string) (*models.PaymentLink, error) {
	query := `
		SELECT id, merchant_id, mode, amount, currency, description, status, expires_at, created_at, updated_at
		FROM payment_links
		WHERE id = $1
	`
	var pl models.PaymentLink
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pl.ID, &pl.MerchantID, &pl.Mode, &pl.Amount, &pl.Currency,
		&pl.Description, &pl.Status, &pl.ExpiresAt, &pl.CreatedAt, &pl.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &pl, nil
}

func (r *PaymentLinkRepository) ListByMerchant(ctx context.Context, merchantID string, limit int) ([]*models.PaymentLink, error) {
	query := `
		SELECT id, merchant_id, mode, amount, currency, description, status, expires_at, created_at, updated_at
		FROM payment_links
		WHERE merchant_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.PaymentLink
	for rows.Next() {
		var pl models.PaymentLink
		if err := rows.Scan(
			&pl.ID, &pl.MerchantID, &pl.Mode, &pl.Amount, &pl.Currency,
			&pl.Description, &pl.Status, &pl.ExpiresAt, &pl.CreatedAt, &pl.UpdatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, &pl)
	}
	return links, rows.Err()
}
