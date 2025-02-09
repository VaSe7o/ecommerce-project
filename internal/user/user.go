package user

import "time"

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Create(u *User) error
	GetByID(id string) (*User, error)
	Update(u *User) error
	Delete(id string) error
	List() ([]*User, error)
}

type Service interface {
	CreateUser(firstName, lastName, email, password string) (*User, error)
	GetUser(id string) (*User, error)
	UpdateUser(u *User) (*User, error)
	DeleteUser(id string) error
	ListUsers() ([]*User, error)
}
