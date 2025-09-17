package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"yray.com/products-service/internal/models"
)

// performs a stock check/update
func (s *Store) Reserve(ctx context.Context, productID, qty int64) (*models.ReserveResult, error) {
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var newStock int64

	query := `
  UPDATE products
  SET stock = stock - $1
  WHERE id = $2 AND stock >= $1
  RETURNING stock
`

	err = tx.QueryRow(
		ctx,
		query,
		qty,
		productID,
	).Scan(&newStock)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// either not found or insufficient stock
			return &models.ReserveResult{
				OK:      false,
				Message: "insufficient stock or product not found",
			}, nil
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.ReserveResult{
		OK:        true,
		Available: newStock,
		Message:   "reserved",
	}, nil
}

// adds stock back (used on cancellation)
func (s *Store) Release(ctx context.Context, productID, qty int64) error {
	query := `UPDATE products SET stock = stock + $1 WHERE id = $2`

	tag, err := s.db.Exec(
		ctx,
		query,
		qty,
		productID,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("product was not found")
	}

	return nil
}
