package payment

import "time"

type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Create(p *Payment) error
	GetByID(id string) (*Payment, error)
	Update(p *Payment) error
	GetByUser(userID string) ([]*Payment, error)
}

type Service interface {
	CreatePayment(orderID, userID string, amount float64, method string) (*Payment, error)
	GetPayment(id string) (*Payment, error)
	ConfirmPayment(id string) (*Payment, error)
	FailPayment(id string) (*Payment, error)
	GetPaymentsByUser(userID string) ([]*Payment, error)
}
