package order

import "time"

type Order struct {
	ID            string      `json:"id"`
	CustomerName  string      `json:"customer_name"`
	Phone         string      `json:"phone"`
	Email         string      `json:"email"`
	Address       string      `json:"address"`
	PaymentMethod string      `json:"payment_method"`
	FinalTotal    float64     `json:"final_total"`
	Items         []OrderItem `json:"items"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Repository interface {
	Create(o *Order) error
	GetByID(id string) (*Order, error)
	Update(o *Order) error
	Delete(id string) error
	GetByUser(userID string) ([]*Order, error)
	List() ([]*Order, error)
}

type Service interface {
	CreateOrder(name, phone, email, address, paymentMethod string, finalTotal float64, items []OrderItem) (*Order, error)
	GetOrder(id string) (*Order, error)
	UpdateOrderStatus(id, newStatus string) (*Order, error)
	DeleteOrder(id string) error
	GetOrdersByUser(userID string) ([]*Order, error)
	ListAllOrders() ([]*Order, error)
}
