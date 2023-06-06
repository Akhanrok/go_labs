package main

import (
	"log"
	"net/http"

	"github.com/Akhanrok/go_labs/handlers/list_handlers"
	"github.com/Akhanrok/go_labs/handlers/user_handlers"
	"github.com/Akhanrok/go_labs/repositories/database_repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

// generate a unique secret key for session cookie store
func generateSecretKey() string {
	key := uuid.New().String()
	return key
}

func main() {
	// Create a database connection
	var err error
	db, err := database_repository.NewDatabase("root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Generate a unique secret key for session cookie store
	secretKey := generateSecretKey()

	// Configure session store
	store = sessions.NewCookieStore([]byte(secretKey))

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user_handlers.IndexHandler(w, r)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		user_handlers.LoginHandler(w, r, db, store)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		user_handlers.RegisterHandler(w, r, db)
	})

	http.HandleFunc("/login-success", func(w http.ResponseWriter, r *http.Request) {
		user_handlers.LoginSuccessHandler(w, r, db)
	})

	http.HandleFunc("/register-success", func(w http.ResponseWriter, r *http.Request) {
		user_handlers.RegisterSuccessHandler(w, r)
	})

	http.HandleFunc("/create-list", func(w http.ResponseWriter, r *http.Request) {
		list_handlers.CreateListHandler(w, r, db, store)
	})

	http.HandleFunc("/list-success", func(w http.ResponseWriter, r *http.Request) {
		list_handlers.ListSuccessHandler(w, r)
	})

	http.HandleFunc("/view-lists", func(w http.ResponseWriter, r *http.Request) {
		list_handlers.ViewListsHandler(w, r, db, store)
	})

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
