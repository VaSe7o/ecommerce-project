package order

import (
	"database/sql"
	"errors"
)

type sqliteOrderRepo struct {
	db *sql.DB
}

func NewSQLiteOrderRepo(db *sql.DB) Repository {
	return &sqliteOrderRepo{db: db}
}

func (r *sqliteOrderRepo) Create(o *Order) error {

	orderQ := `INSERT INTO orders (
        id,
        customer_name,
        phone,
        email,
        address,
        payment_method,
        final_total,
        status,
        created_at,
        updated_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(orderQ,
		o.ID,
		o.CustomerName,
		o.Phone,
		o.Email,
		o.Address,
		o.PaymentMethod,
		o.FinalTotal,
		o.Status,
		o.CreatedAt,
		o.UpdatedAt,
	)
	if err != nil {
		return err
	}
	itemQ := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
	for _, it := range o.Items {
		_, e := r.db.Exec(itemQ, o.ID, it.ProductID, it.Quantity, it.Price)
		if e != nil {
			return e
		}
	}
	return nil
}

func (r *sqliteOrderRepo) GetByID(id string) (*Order, error) {
	var o Order
	q := `SELECT
          id,
          customer_name,
          phone,
          email,
          address,
          payment_method,
          final_total,
          status,
          created_at,
          updated_at
          FROM orders
          WHERE id = ? LIMIT 1`
	row := r.db.QueryRow(q, id)
	err := row.Scan(
		&o.ID,
		&o.CustomerName,
		&o.Phone,
		&o.Email,
		&o.Address,
		&o.PaymentMethod,
		&o.FinalTotal,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}
	itemsQ := `SELECT id, order_id, product_id, quantity, price
               FROM order_items WHERE order_id = ?`
	rows, e := r.db.Query(itemsQ, o.ID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var items []OrderItem
	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.Price); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	o.Items = items
	return &o, nil
}

func (r *sqliteOrderRepo) Update(o *Order) error {
	q := `UPDATE orders
          SET status = ?, updated_at = ?
          WHERE id = ?`
	res, err := r.db.Exec(q, o.Status, o.UpdatedAt, o.ID)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *sqliteOrderRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM order_items WHERE order_id = ?`, id)
	if err != nil {
		return err
	}
	res, e := r.db.Exec(`DELETE FROM orders WHERE id = ?`, id)
	if e != nil {
		return e
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *sqliteOrderRepo) GetByUser(userID string) ([]*Order, error) {

	q := `SELECT
          id,
          customer_name,
          phone,
          email,
          address,
          payment_method,
          final_total,
          status,
          created_at,
          updated_at
          FROM orders
          WHERE customer_name = ?`
	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var all []*Order
	for rows.Next() {
		var o Order
		e := rows.Scan(
			&o.ID,
			&o.CustomerName,
			&o.Phone,
			&o.Email,
			&o.Address,
			&o.PaymentMethod,
			&o.FinalTotal,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if e != nil {
			return nil, e
		}
		itsQ := `SELECT id, order_id, product_id, quantity, price
                 FROM order_items
                 WHERE order_id = ?`
		itsRows, e2 := r.db.Query(itsQ, o.ID)
		if e2 != nil {
			return nil, e2
		}
		var items []OrderItem
		for itsRows.Next() {
			var it OrderItem
			if err := itsRows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.Price); err != nil {
				itsRows.Close()
				return nil, err
			}
			items = append(items, it)
		}
		itsRows.Close()
		o.Items = items
		all = append(all, &o)
	}
	return all, nil
}

func (r *sqliteOrderRepo) List() ([]*Order, error) {
	q := `SELECT
          id,
          customer_name,
          phone,
          email,
          address,
          payment_method,
          final_total,
          status,
          created_at,
          updated_at
          FROM orders`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var all []*Order
	for rows.Next() {
		var o Order
		e := rows.Scan(
			&o.ID,
			&o.CustomerName,
			&o.Phone,
			&o.Email,
			&o.Address,
			&o.PaymentMethod,
			&o.FinalTotal,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if e != nil {
			return nil, e
		}
		itsQ := `SELECT id, order_id, product_id, quantity, price
                 FROM order_items
                 WHERE order_id = ?`
		itsRows, e2 := r.db.Query(itsQ, o.ID)
		if e2 != nil {
			return nil, e2
		}
		var items []OrderItem
		for itsRows.Next() {
			var it OrderItem
			if err := itsRows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.Price); err != nil {
				itsRows.Close()
				return nil, err
			}
			items = append(items, it)
		}
		itsRows.Close()
		o.Items = items
		all = append(all, &o)
	}
	return all, nil
}
