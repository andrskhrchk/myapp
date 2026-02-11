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
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, customerID int,
	items []struct {
		ProductID int
		Quantity  int
	}) (*domain.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var totalPrice float64

	type productInfo struct {
		ProductID int
		Quantity  int
		Price     float64
		Stock     int
	}
	var productsInfo []productInfo

	for _, item := range items {
		var stock int
		var price float64

		query := `SELECT stock, price FROM products WHERE id = $1 FOR UPDATE`
		err := tx.QueryRowContext(ctx, query, item.ProductID).Scan(&stock, &price)
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, errors.New("not enough stock for product id " + string(rune(item.ProductID)))
		}

		totalPrice += price * float64(item.Quantity)
		productsInfo = append(productsInfo, productInfo{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
			Stock:     stock,
		})
	}

	var orderID int
	query := `INSERT INTO orders (customer_id, total_price, status) 
	          VALUES ($1, $2, 'pending') RETURNING id`
	err = tx.QueryRowContext(ctx, query, customerID, totalPrice).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	for _, info := range productsInfo {
		query := `INSERT INTO order_items (order_id, product_id, quantity, price) 
		          VALUES ($1, $2, $3, $4)`
		_, err = tx.ExecContext(ctx, query, orderID, info.ProductID, info.Quantity, info.Price)
		if err != nil {
			return nil, err
		}

		query = `UPDATE products SET stock = stock - $1 WHERE id = $2`
		_, err = tx.ExecContext(ctx, query, info.Quantity, info.ProductID)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:         orderID,
		CustomerId: customerID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, orderID int) (*domain.OrderWithItems, error) {
	order := &domain.Order{}
	query := `SELECT id, customer_id, total_price, status, created_at FROM orders WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID,
		&order.CustomerId,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	query = `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &domain.OrderWithItems{
		Order: *order,
		Items: items,
	}, rows.Err()
}

func (r *OrderRepository) GetByCustomerID(ctx context.Context, customerID int) ([]domain.OrderWithItems, error) {
	query := `SELECT id, customer_id, total_price, status, created_at 
	          FROM orders 
	          WHERE customer_id = $1 
	          ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.OrderWithItems
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.CustomerId, &order.TotalPrice, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}

		itemsQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
		itemRows, err := r.db.QueryContext(ctx, itemsQuery, order.ID)
		if err != nil {
			return nil, err
		}

		var items []domain.OrderItem
		for itemRows.Next() {
			var item domain.OrderItem
			if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, item)
		}
		itemRows.Close()

		orders = append(orders, domain.OrderWithItems{
			Order: order,
			Items: items,
		})
	}

	return orders, rows.Err()
}
