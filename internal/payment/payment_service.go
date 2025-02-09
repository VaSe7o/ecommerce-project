package payment

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

func (s *service) CreatePayment(orderID, userID string, amount float64, method string) (*Payment, error) {
	p := &Payment{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		UserID:    userID,
		Amount:    amount,
		Method:    method,
		Status:    "initiated",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetPayment(id string) (*Payment, error) {
	return s.repo.GetByID(id)
}

func (s *service) ConfirmPayment(id string) (*Payment, error) {
	pay, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	pay.Status = "completed"
	pay.UpdatedAt = time.Now()
	if err := s.repo.Update(pay); err != nil {
		return nil, err
	}
	return pay, nil
}

func (s *service) FailPayment(id string) (*Payment, error) {
	pay, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	pay.Status = "failed"
	pay.UpdatedAt = time.Now()
	if err := s.repo.Update(pay); err != nil {
		return nil, err
	}
	return pay, nil
}

func (s *service) GetPaymentsByUser(userID string) ([]*Payment, error) {
	return s.repo.GetByUser(userID)
}
