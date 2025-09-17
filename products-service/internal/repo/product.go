package repo

import (
	"context"

	"yray.com/products-service/internal/models"
)

// create product by inserting and filling server generated fields
func (s *Store) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
	INSERT INTO products(sku, name, price, stock)
	VALUES($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(
		ctx,
		query,
		product.SKU,
		product.Name,
		product.Price,
		product.Stock,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// fetch a product by id
func (s *Store) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	product := &models.Product{}

	query := `
	SELECT id, sku, name, price, stock, created_at, updated_at
	FROM products WHERE id = $1
	`

	err := s.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}
