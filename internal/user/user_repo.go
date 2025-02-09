package user

import (
	"database/sql"
	"errors"
)

type sqliteUserRepo struct {
	db *sql.DB
}

func NewSQLiteUserRepo(db *sql.DB) Repository {
	return &sqliteUserRepo{db: db}
}

func (r *sqliteUserRepo) Create(u *User) error {
	query := `INSERT INTO users (id, first_name, last_name, email, password, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query,
		u.ID,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.CreatedAt,
		u.UpdatedAt,
	)
	return err
}

func (r *sqliteUserRepo) GetByID(id string) (*User, error) {
	var user User
	query := `SELECT id, first_name, last_name, email, password, created_at, updated_at
              FROM users WHERE id = ? LIMIT 1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqliteUserRepo) Update(u *User) error {
	query := `UPDATE users
              SET first_name = ?, last_name = ?, email = ?, password = ?, updated_at = ?
              WHERE id = ?`
	res, err := r.db.Exec(query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.UpdatedAt,
		u.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *sqliteUserRepo) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *sqliteUserRepo) List() ([]*User, error) {
	rows, err := r.db.Query(`SELECT id, first_name, last_name, email, password, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
