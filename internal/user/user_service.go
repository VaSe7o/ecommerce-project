package user

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

func (s *service) CreateUser(firstName, lastName, email, password string) (*User, error) {
	u := &User{
		ID:        uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) GetUser(id string) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateUser(u *User) (*User, error) {
	u.UpdatedAt = time.Now()
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

func (s *service) ListUsers() ([]*User, error) {
	return s.repo.List()
}
