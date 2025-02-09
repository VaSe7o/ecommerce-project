package product

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Repository interface {
	Create(p *Product) error
	GetByID(id string) (*Product, error)
	Update(p *Product) error
	Delete(id string) error
	List() ([]*Product, error)
}

type Service interface {
	CreateProduct(name, desc string, price float64, quantity int) (*Product, error)
	GetProduct(id string) (*Product, error)
	UpdateProduct(p *Product) (*Product, error)
	DeleteProduct(id string) error
	ListProducts() ([]*Product, error)
}
