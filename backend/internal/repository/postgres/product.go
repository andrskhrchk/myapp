package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/andrskhrchk/myapp/internal/domain"
)

var (
	ErrProdNotFound = errors.New("user not found")
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (name, desc, stock, price) VALUES ($1, $2, $3, $4)`

	return r.db.QueryRowContext(ctx, query, product.Name, product.Desc, product.Stock, product.Price).Scan(&product.ID)
}

func (r *ProductRepository) GetProdByName(ctx context.Context, name string) (*[]domain.Product, error) {
	query := `SELECT * FROM products WHERE name = $1`

	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Desc, &p.Stock, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &products, nil
}

func (r *ProductRepository) GetProdById(ctx context.Context, id int) (*domain.Product, error) {
	query := `SELECT * FROM products WHERE id = $1`
	var p domain.Product

	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Desc, &p.Stock, &p.Price)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) (*[]domain.Product, error) {
	query := `SELECT * FROM products`
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var p domain.Product

		if err := rows.Scan(&p.ID, &p.Name, &p.Desc, &p.Stock, &p.Price); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &products, err
}

func (r *ProductRepository) UpdateProd(ctx context.Context, id int, product *domain.Product) error {
	query := `
		UPDATE products 
		SET name = $1, desc = $2, stock = $3, price = $4
		WHERE id = $5
	`

	p, err := r.db.ExecContext(ctx, query, product.Name, product.Desc, product.Stock, product.Price, id)
	if err != nil {
		return err
	}

	rows, err := p.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *ProductRepository) DeleteProd(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}
	return err
}
