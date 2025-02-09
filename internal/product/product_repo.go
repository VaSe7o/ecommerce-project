package product

import (
	"database/sql"
	"errors"
)

type sqliteProductRepo struct {
	db *sql.DB
}

func NewSQLiteProductRepo(db *sql.DB) Repository {
	return &sqliteProductRepo{db: db}
}

func (r *sqliteProductRepo) Create(p *Product) error {
	q := `INSERT INTO products (id, name, description, price, quantity, created_at, updated_at)
          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(q,
		p.ID,
		p.Name,
		p.Description,
		p.Price,
		p.Quantity,
		p.CreatedAt,
		p.UpdatedAt,
	)
	return err
}

func (r *sqliteProductRepo) GetByID(id string) (*Product, error) {
	var prod Product
	q := `SELECT id, name, description, price, quantity, created_at, updated_at
          FROM products
          WHERE id = ? LIMIT 1`
	row := r.db.QueryRow(q, id)
	err := row.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Quantity,
		&prod.CreatedAt, &prod.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *sqliteProductRepo) Update(p *Product) error {
	q := `UPDATE products
          SET name = ?, description = ?, price = ?, quantity = ?, updated_at = ?
          WHERE id = ?`
	res, err := r.db.Exec(q,
		p.Name,
		p.Description,
		p.Price,
		p.Quantity,
		p.UpdatedAt,
		p.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *sqliteProductRepo) Delete(id string) error {
	q := `DELETE FROM products WHERE id = ?`
	res, err := r.db.Exec(q, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *sqliteProductRepo) List() ([]*Product, error) {
	q := `SELECT id, name, description, price, quantity, created_at, updated_at
          FROM products`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}
