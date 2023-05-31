package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Akhanrok/go_labs/database"
	"github.com/Akhanrok/go_labs/models"
)

type ShoppingListHandler struct {
	db *database.Database
}

func NewShoppingListHandler(db *database.Database) *ShoppingListHandler {
	return &ShoppingListHandler{db: db}
}

func createListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listName := r.PostForm.Get("listName")

		// Check if the list name already exists in the database
		listExists, err := isListExists(db, listName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if listExists {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "The list with such name already exists",
			}
			renderTemplate(w, "create-list.html", data)
			return
		}

		// Insert the new list into the database
		insertListQuery := "INSERT INTO lists (user_id, name) VALUES (?, ?)"
		res, err := db.Exec(insertListQuery, 1, listName) // Replace 1 with the appropriate user ID
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listID, _ := res.LastInsertId()

		// Get the product, quantity, and store values from the form
		products := r.PostForm["product[]"]
		quantities := r.PostForm["quantity[]"]
		stores := r.PostForm["store[]"]

		// Insert each product into the database
		insertProductQuery := "INSERT INTO products (list_id, name, quantity, store) VALUES (?, ?, ?, ?)"
		for i := range products {
			_, err = db.Exec(insertProductQuery, listID, products[i], quantities[i], stores[i])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Redirect to the list success page
		http.Redirect(w, r, fmt.Sprintf("/list-success?name=%s", listName), http.StatusFound)
		return
	}

	renderTemplate(w, "create-list.html", nil)
}

func isListExists(db *sql.DB, listName string) (bool, error) {
	query := "SELECT COUNT(*) FROM lists WHERE name = ?"
	var count int
	err := db.QueryRow(query, listName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func listSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		listName := r.URL.Query().Get("name")
		data := struct {
			ListName string
		}{
			ListName: listName,
		}
		renderTemplate(w, "list-success.html", data)
	}
}

func viewListsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Retrieve the list names and products from the database
		lists, err := getListsData(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Lists []ListData
		}{
			Lists: lists,
		}

		renderTemplate(w, "view-lists.html", data)
	}
}

func getListsData(db *sql.DB) ([]ListData, error) {
	query := "SELECT id, name FROM lists"
	rows, err := db.Query(query)
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

		products, err := getProductsData(db, listID)
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

func getProductsData(db *sql.DB, listID int) ([]Product, error) {
	query := "SELECT name, quantity, store FROM products WHERE list_id = ?"
	rows, err := db.Query(query, listID)
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
