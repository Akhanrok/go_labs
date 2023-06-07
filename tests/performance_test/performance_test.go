package performance_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Akhanrok/go_labs/handlers/list_handlers"
	"github.com/Akhanrok/go_labs/handlers/user_handlers"
	"github.com/Akhanrok/go_labs/repositories/user_repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var (
	db    *sql.DB
	store *sessions.CookieStore
)

func TestCheckCredentialsPerformance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := user_repository.NewUserRepository(db)

	_, err = userRepo.ValidateCredentials("test@example.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestCreateListHandlerPerformance(t *testing.T) {
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

	// Prepare the test response recorder
	recorder := httptest.NewRecorder()

	// Call the handler function
	rr := httptest.NewRecorder()
	list_handlers.CreateListHandler(recorder, req, db, store)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestLoginHandlerPerformance(t *testing.T) {
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

	// Call the handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_handlers.LoginHandler(w, r, db, store)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestRegisterHandlerPerformance(t *testing.T) {
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

	// Call the handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_handlers.RegisterHandler(w, r, db)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestViewListsHandlerPerformance(t *testing.T) {
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

	// Call the handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list_handlers.ViewListsHandler(w, r, db, store)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}
