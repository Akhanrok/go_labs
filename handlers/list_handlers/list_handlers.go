package list_handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Akhanrok/go_labs/repositories/list_repository"
	"github.com/Akhanrok/go_labs/services"
	"github.com/gorilla/sessions"
)

func CreateListHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, store sessions.Store) {
	if r.Method == http.MethodPost {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userID := session.Values["userID"].(int) // Get the user ID from the session

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listName := r.PostForm.Get("listName")

		// Create instances of the repositories
		listRepo := list_repository.NewListRepository(db)

		// Check if the list name already exists for the user in the database
		listExists, err := listRepo.IsListExists(userID, listName)
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
			services.RenderTemplate(w, "create-list.html", data)
			return
		}

		// Insert the new list into the database
		insertListQuery := "INSERT INTO lists (user_id, name) VALUES (?, ?)"
		res, err := db.Exec(insertListQuery, userID, listName)
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

	services.RenderTemplate(w, "create-list.html", nil)
}

func ListSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		listName := r.URL.Query().Get("name")
		data := struct {
			ListName string
		}{
			ListName: listName,
		}
		services.RenderTemplate(w, "list-success.html", data)
	}
}

func ViewListsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, store sessions.Store) {
	if r.Method == http.MethodGet {
		// Get the user ID of the currently authenticated user from the session
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userID := session.Values["userID"].(int)

		// Create an instance of the ListRepository
		listRepo := list_repository.NewListRepository(db)

		// Retrieve the list names and items for the user from the database
		lists, err := listRepo.GetListsData(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Lists []list_repository.ListData
		}{
			Lists: lists,
		}

		services.RenderTemplate(w, "view-lists.html", data)
	}
}
