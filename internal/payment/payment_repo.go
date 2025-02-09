package payment

import (
	"database/sql"
	"errors"
)

type sqlitePaymentRepo struct {
	db *sql.DB
}

func NewSQLitePaymentRepo(db *sql.DB) Repository {
	return &sqlitePaymentRepo{db: db}
}

func (r *sqlitePaymentRepo) Create(p *Payment) error {
	q := `INSERT INTO payments
          (id, order_id, user_id, amount, method, status, created_at, updated_at)
          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(q, p.ID, p.OrderID, p.UserID, p.Amount, p.Method, p.Status, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *sqlitePaymentRepo) GetByID(id string) (*Payment, error) {
	var pay Payment
	q := `SELECT id, order_id, user_id, amount, method, status, created_at, updated_at
          FROM payments
          WHERE id = ? LIMIT 1`
	row := r.db.QueryRow(q, id)
	err := row.Scan(
		&pay.ID,
		&pay.OrderID,
		&pay.UserID,
		&pay.Amount,
		&pay.Method,
		&pay.Status,
		&pay.CreatedAt,
		&pay.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("payment not found")
	}
	if err != nil {
		return nil, err
	}
	return &pay, nil
}

func (r *sqlitePaymentRepo) Update(p *Payment) error {
	q := `UPDATE payments
          SET order_id = ?, user_id = ?, amount = ?, method = ?, status = ?, updated_at = ?
          WHERE id = ?`
	res, err := r.db.Exec(q,
		p.OrderID,
		p.UserID,
		p.Amount,
		p.Method,
		p.Status,
		p.UpdatedAt,
		p.ID,
	)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("payment not found")
	}
	return nil
}

func (r *sqlitePaymentRepo) GetByUser(userID string) ([]*Payment, error) {
	q := `SELECT id, order_id, user_id, amount, method, status, created_at, updated_at
          FROM payments
          WHERE user_id = ?`
	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pays []*Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.UserID,
			&p.Amount,
			&p.Method,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pays = append(pays, &p)
	}
	return pays, nil
}
