package repositories_test

import (
	"database/sql"
	"testing"

	"github.com/Akhanrok/go_labs/repositories/database_repository"
	"github.com/Akhanrok/go_labs/repositories/list_repository"
	"github.com/Akhanrok/go_labs/repositories/product_repository"
	"github.com/Akhanrok/go_labs/repositories/user_repository"
	_ "github.com/go-sql-driver/mysql"
)

func TestNewDatabase(t *testing.T) {
	dataSourceName := "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app"

	// Call the NewDatabase function
	db, err := database_repository.NewDatabase(dataSourceName)

	// Check if there was an error creating the database connection
	if err != nil {
		t.Errorf("failed to create database connection: %v", err)
	}

	// Check if the database connection is not nil
	if db == nil {
		t.Error("database connection is nil")
	}

	db.Close()
}

func TestIsListExists(t *testing.T) {
	dataSourceName := "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app"

	// Create a database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Create an instance of the listRepository
	repo := list_repository.NewListRepository(db)

	// Define the test parameters
	userID := 1
	listName := "example list"

	// Call the IsListExists function
	exists, err := repo.IsListExists(userID, listName)

	// Check if there was an error
	if err != nil {
		t.Errorf("failed to check list existence: %v", err)
	}

	if exists {
		t.Errorf("expected list '%s' for user %d to not exist, but it exists", listName, userID)
	}
}

func TestGetListsData(t *testing.T) {
	dataSourceName := "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app"

	// Create a database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Create an instance of the listRepository
	repo := list_repository.NewListRepository(db)

	// Define the test parameter
	userID := 1

	// Call the GetListsData function
	lists, err := repo.GetListsData(userID)

	if err != nil {
		t.Errorf("failed to get lists data: %v", err)
	}

	if len(lists) == 0 {
		t.Errorf("expected to get lists data for user %d, but received an empty list", userID)
	}
}

func TestValidateCredentials(t *testing.T) {
	dataSourceName := "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app"

	// Create a database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Create an instance of the userRepository
	repo := user_repository.NewUserRepository(db)

	// Define the test parameters
	email := "test@example.com"
	password := "password"

	// Call the ValidateCredentials function
	username, err := repo.ValidateCredentials(email, password)

	if err != nil {
		t.Errorf("failed to validate credentials: %v", err)
	}

	if username != "expectedUsername" {
		t.Errorf("expected username 'expectedUsername', but received '%s'", username)
	}
}

func TestIsEmailExists(t *testing.T) {
	dataSourceName := "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app"

	// Create a database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Create an instance of the userRepository
	repo := user_repository.NewUserRepository(db)

	// Define the test parameter
	email := "test@example.com"

	// Call the IsEmailExists function
	exists, err := repo.IsEmailExists(email)

	if err != nil {
		t.Errorf("failed to check email existence: %v", err)
	}

	if !exists {
		t.Errorf("expected email '%s' to exist, but it does not exist", email)
	}
}

func TestGetProductsData(t *testing.T) {
	dataSourceName := "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app"

	// Create a database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// Create an instance of the productRepository
	repo := product_repository.NewProductRepository(db)

	// Define the test parameter
	listID := 1

	// Call the GetProductsData function
	products, err := repo.GetProductsData(listID)

	if err != nil {
		t.Errorf("failed to get products data: %v", err)
	}

	if len(products) == 0 {
		t.Errorf("expected to get products data for list %d, but received an empty list", listID)
	}
}
