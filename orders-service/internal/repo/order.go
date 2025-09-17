package repo

import (
	"context"

	"yray.com/orders-service/internal/models"
)

// CreateOrder inserts a new order and its associated items into the database.
//
// The operation is performed inside a transaction to ensure consistency:
//   - The order is first inserted into the "orders" table.
//   - The generated order ID is then used to insert all items into the "order_items" table.
//   - If any insert fails, the transaction is rolled back automatically via defer.
//   - On success, the transaction is committed and the new order ID is returned.
//
// Returns the newly created order ID or an error if the operation fails.
func (s *Store) CreateOrder(ctx context.Context, order *models.Order) (int64, error) {
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	orderQuery := `
	INSERT INTO orders(customer_id, status, total)
	VALUES($1,$2,$3)
	RETURNING id
	`
	if err := tx.QueryRow(
		ctx,
		orderQuery,
		order.CustomerID,
		order.Status,
		order.Total,
	).Scan(
		&id,
	); err != nil {
		return 0, err
	}

	orderItemQuery := `
	INSERT INTO order_items(order_id, product_id, qty, price)
	VALUES($1,$2,$3,$4)
	`
	for _, it := range order.Items {
		if _, err := tx.Exec(
			ctx,
			orderItemQuery,
			id,
			it.ProductID,
			it.Qty,
			it.Price,
		); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}

// GetOrder retrieves an order by its ID, including all associated items.
//
//   - First queries the "orders" table for the main order details.
//   - Then fetches related items from "order_items" and appends them to the result.
//
// Returns the fully populated Order model or an error if not found
func (s *Store) GetOrder(ctx context.Context, id int64) (*models.Order, error) {
	order := &models.Order{}
	orderQuery := `
	SELECT id, customer_id, status, total, created_at, updated_at FROM orders WHERE id=$1
	`
	if err := s.db.QueryRow(
		ctx,
		orderQuery,
		id,
	).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, err
	}

	orderItemQuery := `SELECT product_id, qty, price_cents FROM order_items WHERE order_id=$1`
	rows, err := s.db.Query(
		ctx,
		orderItemQuery,
		id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Qty, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

// UpdateStatus updates the status field of a specific order.
//
// Returns an error if the update fails or the order does not exist.
func (s *Store) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE orders SET status=$1 WHERE id=$2`

	_, err := s.db.Exec(
		ctx,
		query,
		status,
		id,
	)

	return err
}
