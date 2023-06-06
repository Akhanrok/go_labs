package list_repository

import (
	"database/sql"

	"github.com/Akhanrok/go_labs/repositories/product_repository"
)

type ListData struct {
	ListName string
	Products []product_repository.Product
}

type ListRepository interface {
	IsListExists(userID int, listName string) (bool, error)
	GetListsData(userID int) ([]ListData, error)
}

type listRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) ListRepository {
	return &listRepository{db}
}

func (r *listRepository) IsListExists(userID int, listName string) (bool, error) {
	query := "SELECT COUNT(*) FROM lists WHERE user_id = ? AND name = ?"
	var count int
	err := r.db.QueryRow(query, userID, listName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *listRepository) GetListsData(userID int) ([]ListData, error) {
	query := "SELECT id, name FROM lists WHERE user_id = ?"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []ListData

	for rows.Next() {
		var listID int
		var listName string

		err := rows.Scan(&listID, &listName)
		if err != nil {
			return nil, err
		}

		products, err := product_repository.NewProductRepository(r.db).GetProductsData(listID)
		if err != nil {
			return nil, err
		}

		listData := ListData{
			ListName: listName,
			Products: products,
		}

		lists = append(lists, listData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}
