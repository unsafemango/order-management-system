package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations applies the required database schema for the application.
//
// It creates the following tables if they do not already exist:
//   - orders: stores basic order information (id, customer, status, total, timestamps)
//   - order_items: stores line items belonging to orders (product, qty, price)
//
// It also ensures that:
//   - A trigger function (set_updated_at) exists to automatically update the
//     `updated_at` timestamp on the orders table whenever a row is updated.
//   - The trigger (orders_updated_at) is created only if it does not already exist.
//
// Returns an error if executing the migration SQL fails.
func RunMigrations(db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders(
		id BIGSERIAL PRIMARY KEY,
		customer_id TEXT NOT NULL,
		status TEXT NOT NULL,
		total BIGINT NOT NULL DEFAULT 0,
		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS order_items(
		order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		product_id BIGINT NOT NULL,
		qty BIGINT NOT NULL,
		price BIGINT NOT NULL
	);

	CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
	BEGIN NEW.updated_at = now(); RETURN NEW; END; $$ LANGUAGE plpgsql;
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname='orders_updated_at') THEN
		CREATE TRIGGER orders_updated_at BEFORE UPDATE ON orders
		FOR EACH ROW EXECUTE FUNCTION set_updated_at();
	END IF; END $$;
	`

	_, err := db.Exec(context.Background(), query)

	return err
}
