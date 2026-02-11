package domain

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"-"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Desc  string  `json:"desc"`
	Stock int     `json:"stock"`
	Price float32 `json:"price"`
}

type ProductImage struct {
	ID       int    `json:"id"`
	ProdID   int    `json:"prod_id"`
	Filepath string `json:"filepath"`
}

type Order struct {
	ID         int       `json:"id"`
	CustomerId int       `json:"customer_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderWithItems struct {
	Order
	Items []OrderItem `json:"items"`
}
