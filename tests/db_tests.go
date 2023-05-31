package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Akhanrok/go_labs/handlers"
)

func TestLoginHandler(t *testing.T) {
	// Create a test database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app_test")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the form data for the request
	values := make(map[string][]string)
	values["email"] = []string{"test@example.com"}
	values["password"] = []string{"password123"}
	req.PostForm = values

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the loginHandler with the test request and response recorder
	http.HandlerFunc(loginHandler).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestCreateListHandler(t *testing.T) {
	// Create a test database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app_test")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/create-list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the form data for the request
	values := make(map[string][]string)
	values["listName"] = []string{"Test List"}
	values["product[]"] = []string{"Product 1", "Product 2"}
	values["quantity[]"] = []string{"1", "2"}
	values["store[]"] = []string{"Store 1", "Store 2"}
	req.PostForm = values

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the createListHandler with the test request and response recorder
	http.HandlerFunc(createListHandler).ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d", http.StatusFound, rr.Code)
	}
}

