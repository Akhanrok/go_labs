package database

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/Akhanrok/go_labs/models"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

type UserNotFoundError struct {
	Email string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found: %s", e.Email)
}

type ProductNotFoundError struct {
	ProductID uint
}

func (e ProductNotFoundError) Error() string {
	return fmt.Sprintf("product not found: %d", e.ProductID)
}

var (
	ErrInvalidEmail = errors.New("invalid email address")
)

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("mysql", "config")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) CreateUser(user *models.User) error {
	if !IsValidEmail(user.Email) {
		return ErrInvalidEmail
	}

	statement := `INSERT INTO users (name, email, password, created_at) VALUES (?, ?, ?, ?)`
	_, err := d.db.Exec(statement, user.Name, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetUserByEmail(email string) (*models.User, error) {
	statement := `SELECT id, name, email, password, created_at FROM users WHERE email = ?`
	row := d.db.QueryRow(statement, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (d *Database) CreateShoppingList(list *models.ShoppingList) error {
	statement := `INSERT INTO shopping_lists (user_id, name, store_name, quantity, link, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := d.db.Exec(statement, list.UserID, list.Name, list.StoreName, list.Quantity, list.Link, list.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) AddProductToList(product *models.Product) error {
	statement := `INSERT INTO products (shopping_list_id, name, type, created_at) VALUES (?, ?, ?, ?)`
	_, err := d.db.Exec(statement, product.ShoppingListID, product.Name, product.Type, product.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) RemoveProductFromList(listID, productID uint) error {
	stmt, err := d.db.Prepare("DELETE FROM shopping_list_products WHERE shopping_list_id = ? AND product_id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare product removal statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(listID, productID)
	if err != nil {
		return fmt.Errorf("failed to execute product removal statement: %w", err)
	}

	return nil
}

func (d *Database) SortShoppingListByType(listID uint) error {
	stmt, err := d.db.Prepare("SELECT * FROM shopping_list_products WHERE shopping_list_id = ? ORDER BY type")
	if err != nil {
		return fmt.Errorf("failed to prepare shopping list sorting statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(listID)
	if err != nil {
		return fmt.Errorf("failed to execute shopping list sorting statement: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.ShoppingListID)
		if err != nil {
			return fmt.Errorf("failed to scan product row: %w", err)
		}
	}

	return nil
}

func (d *Database) CreateProduct(product *models.Product) error {
	stmt, err := d.db.Prepare("INSERT INTO products (name, type) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare product creation statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Type)
	if err != nil {
		return fmt.Errorf("failed to execute product creation statement: %w", err)
	}

	return nil
}

func (d *Database) GetProduct(productID uint) (*models.Product, error) {
	stmt, err := d.db.Prepare("SELECT id, name, type FROM products WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare product retrieval statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(productID)

	product := &models.Product{}

	err = row.Scan(&product.ID, &product.Name, &product.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	return product, nil
}

func (d *Database) UpdateProduct(product *models.Product) error {
	stmt, err := d.db.Prepare("UPDATE products SET name = ?, type = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare product update statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Type, product.ID)
	if err != nil {
		return fmt.Errorf("failed to execute product update statement: %w", err)
	}

	return nil
}

func (d *Database) DeleteProduct(productID uint) error {
	stmt, err := d.db.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare product deletion statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(productID)
	if err != nil {
		return fmt.Errorf("failed to execute product deletion statement: %w", err)
	}

	return nil
}
