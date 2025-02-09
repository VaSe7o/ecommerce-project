package order

import (
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateOrder(name, phone, email, address, paymentMethod string, finalTotal float64, items []OrderItem) (*Order, error) {
	o := &Order{
		ID:            uuid.New().String(),
		CustomerName:  name,
		Phone:         phone,
		Email:         email,
		Address:       address,
		PaymentMethod: paymentMethod,
		FinalTotal:    finalTotal,
		Items:         items,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.repo.Create(o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *service) GetOrder(id string) (*Order, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateOrderStatus(id, newStatus string) (*Order, error) {
	ord, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	ord.Status = newStatus
	ord.UpdatedAt = time.Now()
	if e := s.repo.Update(ord); e != nil {
		return nil, e
	}
	return ord, nil
}

func (s *service) DeleteOrder(id string) error {
	return s.repo.Delete(id)
}

func (s *service) GetOrdersByUser(userID string) ([]*Order, error) {
	return s.repo.GetByUser(userID)
}

func (s *service) ListAllOrders() ([]*Order, error) {
	return s.repo.List()
}
