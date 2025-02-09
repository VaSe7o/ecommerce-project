package product

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

func (s *service) CreateProduct(name, desc string, price float64, quantity int) (*Product, error) {
	p := &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: desc,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetProduct(id string) (*Product, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateProduct(p *Product) (*Product, error) {
	p.UpdatedAt = time.Now()
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) DeleteProduct(id string) error {
	return s.repo.Delete(id)
}

func (s *service) ListProducts() ([]*Product, error) {
	return s.repo.List()
}
