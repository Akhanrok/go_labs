package product_repository

import (
	"database/sql"
)

type Product struct {
	Product  string
	Quantity int
	Store    string
}

type ProductRepository interface {
	GetProductsData(listID int) ([]Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) GetProductsData(listID int) ([]Product, error) {
	query := "SELECT name, quantity, store FROM products WHERE list_id = ?"
	rows, err := r.db.Query(query, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var name string
		var quantity int
		var store string

		err := rows.Scan(&name, &quantity, &store)
		if err != nil {
			return nil, err
		}

		product := Product{
			Product:  name,
			Quantity: quantity,
			Store:    store,
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
