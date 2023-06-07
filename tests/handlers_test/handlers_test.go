package handlers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(user_handlers.IndexHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	expected := "Index Page"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestLoginHandlerValidCredentials(t *testing.T) {
	// Initialize the test database with a user
	userRepo := user_repository.NewUserRepository(db)
	_, err := userRepo.ValidateCredentials("test@example.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test login request with valid credentials
	payload := strings.NewReader("email=test@example.com&password=password123")
	req, err := http.NewRequest("POST", "/login", payload)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_handlers.LoginHandler(w, r, db, store)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	expected := "Login Success"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	// Create a test login request with invalid credentials
	payload := strings.NewReader("email=test@example.com&password=wrongpassword")
	req, err := http.NewRequest("POST", "/login", payload)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_handlers.LoginHandler(w, r, db, store)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v, want %v", rr.Code, http.StatusUnauthorized)
	}

	expected := "Invalid Credentials"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestCreateListHandler(t *testing.T) {
	// Prepare the test request
	form := url.Values{}
	form.Set("listName", "Test List")
	form.Set("product[]", "Product 1")
	form.Set("quantity[]", "2")
	form.Set("store[]", "Store 1")

	req, err := http.NewRequest("POST", "/create-list", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Prepare the test response recorder
	recorder := httptest.NewRecorder()

	// Call the handler function
	list_handlers.CreateListHandler(recorder, req, db, store)

	// Check the response status code
	if recorder.Code != http.StatusFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusFound, recorder.Code)
	}

	// Check the response header for the redirect location
	location, err := recorder.Result().Location()
	if err != nil {
		t.Fatal(err)
	}

	expectedLocation := "/list-success?name=Test%20List"
	if location.Path != expectedLocation {
		t.Errorf("Expected redirect location %s, but got %s", expectedLocation, location.Path)
	}
}
