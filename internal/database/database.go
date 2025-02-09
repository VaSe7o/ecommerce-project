package database

import (
	"database/sql"
	"fmt"
)

func OpenDB(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	userTable := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        first_name TEXT,
        last_name TEXT,
        email TEXT,
        password TEXT,
        created_at DATETIME,
        updated_at DATETIME
    );`

	productTable := `CREATE TABLE IF NOT EXISTS products (
        id TEXT PRIMARY KEY,
        name TEXT,
        description TEXT,
        price REAL,
        quantity INTEGER,
        created_at DATETIME,
        updated_at DATETIME
    );`

	orderTable := `CREATE TABLE IF NOT EXISTS orders (
        id TEXT PRIMARY KEY,
        customer_name TEXT,
        phone TEXT,
        email TEXT,
        address TEXT,
        payment_method TEXT,
        final_total REAL,
        status TEXT,
        created_at DATETIME,
        updated_at DATETIME
    );`

	orderItemTable := `CREATE TABLE IF NOT EXISTS order_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        order_id TEXT,
        product_id TEXT,
        quantity INTEGER,
        price REAL
    );`

	paymentTable := `CREATE TABLE IF NOT EXISTS payments (
        id TEXT PRIMARY KEY,
        order_id TEXT,
        user_id TEXT,
        amount REAL,
        method TEXT,
        status TEXT,
        created_at DATETIME,
        updated_at DATETIME
    );`

	stmts := []string{
		userTable,
		productTable,
		orderTable,
		orderItemTable,
		paymentTable,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("failed creating table: %w", err)
		}
	}
	return nil
}
