package routes

import (
	"github.com/Akhanrok/go_labs/database"
	"github.com/Akhanrok/go_labs/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(db *database.Database) *mux.Router {
	router := mux.NewRouter()

	userHandler := handlers.NewUserHandler(db)
	productHandler := handlers.NewProductHandler(db)
	shoppingListHandler := handlers.NewShoppingListHandler(db)

	// User routes
	router.HandleFunc("/registerUser", userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/loginUser", userHandler.LoginUser).Methods("POST")

	// Product routes
	router.HandleFunc("/product", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/product/{productID}", productHandler.GetProduct).Methods("GET")
	router.HandleFunc("/product/{productID}", productHandler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/product/{productID}", productHandler.DeleteProduct).Methods("DELETE")

	// Shopping list routes
	router.HandleFunc("/shopping-list", shoppingListHandler.CreateShoppingList).Methods("POST")
	router.HandleFunc("/shopping-list/{listID}/product", shoppingListHandler.AddProductToList).Methods("POST")
	router.HandleFunc("/shopping-list/{listID}/product/{productID}", shoppingListHandler.RemoveProductFromList).Methods("DELETE")
	router.HandleFunc("/shopping-list/{listID}/sort/type", shoppingListHandler.SortShoppingListByType).Methods("GET")
	router.HandleFunc("/shopping-lists/{listID}/sort-by-type", shoppingListHandler.SortShoppingListByType)

	return router
}
