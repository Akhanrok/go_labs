package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Akhanrok/go_labs/handlers"
)

func TestCheckCredentials_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Call the function to be tested
	_, err = checkCredentials(db, "test@example.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestCreateListHandler_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a sample HTTP request for testing
	req, err := http.NewRequest("POST", "/create-list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the handler function to be tested
	rr := httptest.NewRecorder()
	http.HandlerFunc(createListHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestLoginHandler_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a sample HTTP request for testing
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add form data to the request
	formData := url.Values{}
	formData.Set("email", "test@example.com")
	formData.Set("password", "password123")
	req.PostForm = formData

	// Create a response recorder for capturing the response
	rr := httptest.NewRecorder()

	// Call the handler function to be tested
	http.HandlerFunc(loginHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestRegisterHandler_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a sample HTTP request for testing
	req, err := http.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add form data to the request
	formData := url.Values{}
	formData.Set("name", "John Doe")
	formData.Set("email", "johndoe@example.com")
	formData.Set("password", "password123")
	req.PostForm = formData

	// Create a response recorder for capturing the response
	rr := httptest.NewRecorder()

	// Call the handler function to be tested
	http.HandlerFunc(registerHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestViewListsHandler_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a sample HTTP request for testing
	req, err := http.NewRequest("GET", "/view-lists", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder for capturing the response
	rr := httptest.NewRecorder()

	// Call the handler function to be tested
	http.HandlerFunc(viewListsHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestUpdateListHandler_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a sample HTTP request for testing
	req, err := http.NewRequest("POST", "/update-list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add form data to the request
	formData := url.Values{}
	formData.Set("list_id", "1")
	formData.Set("name", "Updated List Name")
	req.PostForm = formData

	// Create a response recorder for capturing the response
	rr := httptest.NewRecorder()

	// Call the handler function to be tested
	http.HandlerFunc(updateListHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestDeleteListHandler_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a sample HTTP request for testing
	req, err := http.NewRequest("POST", "/delete-list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add form data to the request
	formData := url.Values{}
	formData.Set("list_id", "1")
	req.PostForm = formData

	// Create a response recorder for capturing the response
	rr := httptest.NewRecorder()

	// Call the handler function to be tested
	http.HandlerFunc(deleteListHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}
