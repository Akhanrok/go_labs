package routes

import (
	"net/http"
	"github.com/gorilla/mux"

	"github.com/Akhanrok/go_labs/handlers"
	"github.com/Akhanrok/go_labs/database"
	
)

func RegisterRoutes(db *database.Database) *mux.Router {
	router := mux.NewRouter()

	userHandler := handlers.NewUserHandler(db)
	listHandler := handlers.NewListHandler(db)

	// User routes
	http.HandleFunc("/", userHandler.IndexHandler)
	http.HandleFunc("/login", userHandler.LoginHandler)
	http.HandleFunc("/register", userHandler.RegisterHandler)
	http.HandleFunc("/register-success", userHandler.RegisterSuccessHandler)
	http.HandleFunc("/login-success", userHandler.LoginSuccessHandler)

	// Shopping list routes
	http.HandleFunc("/create-list", listHandler.CreateListHandler)
	http.HandleFunc("/list-success", listHandler.ListSuccessHandler)
	http.HandleFunc("/view-lists", listHandler.ViewListsHandler)

	return router
}
