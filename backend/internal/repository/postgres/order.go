package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/andrskhrchk/myapp/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, customerID, productID, qty int) (*domain.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var stock int
	var price float32
	query := `SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, productID).Scan(&stock, &price)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	if stock < qty {
		return nil, errors.New("not enough stock")
	}

	totalPrice := price * float32(qty)

	order := &domain.Order{
		CustomerId: customerID,
		ProdID:     productID,
		Qty:        qty,
		TotPrice:   totalPrice,
	}

	query = `INSERT INTO orders (customer_id, prod_id, qty, tot_price) 
	         VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRowContext(ctx, query, order.CustomerId, order.ProdID, order.Qty, order.TotPrice).Scan(&order.ID)
	if err != nil {
		return nil, err
	}

	query = `UPDATE products SET stock = stock - $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, query, qty, productID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}
