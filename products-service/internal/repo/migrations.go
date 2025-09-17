package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations applies the required database schema for the application
//
// It creates the following tables if they do not already exist:
//   - products: stores basic product information (id, sku, name, price, stock, timestamps)
//
// It also ensures that:
//   - A trigger function (set_updated_at) exists to automatically update the
//     `updated_at` column whenever a row in "products" is modified.
//   - The trigger (products_updated_at) is created only if it does not already exist.
//
// Returns an error if executing the migration SQL fails.
func RunMigrations(db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS products(
		id BIGSERIAL PRIMARY KEY,
		sku TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		price BIGINT NOT NULL,
		stock BIGINT NOT NULL DEFAULT 0,
		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
	);
	CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
	BEGIN NEW.updated_at = now(); RETURN NEW; END; $$ LANGUAGE plpgsql;
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname='products_updated_at') THEN
		CREATE TRIGGER products_updated_at BEFORE UPDATE ON products
		FOR EACH ROW EXECUTE FUNCTION set_updated_at();
	END IF; END $$;
	`

	_, err := db.Exec(context.Background(), query)

	return err
}
