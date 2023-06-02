package tests

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Akhanrok/go_labs/handlers"
	"github.com/Akhanrok/go_labs/repositories"
	"github.com/Akhanrok/go_labs/services"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

var (
	db    *sql.DB
	store *sessions.CookieStore
)

func TestMain(m *testing.M) {
	// Create a test database connection
	var err error
	db, err = repositories.NewDatabase("root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app_test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Generate a unique secret key for session cookie store
	secretKey := uuid.New().String()

	// Configure session store
	store = sessions.NewCookieStore([]byte(secretKey))

	// Run the tests
	exitCode := m.Run()

	// Exit with the test exit code
	os.Exit(exitCode)
}

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.IndexHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	expected := "Index Page"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestLoginHandler_ValidCredentials(t *testing.T) {
	// Initialize the test database with a user
	userRepo := repositories.NewUserRepository(db)
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
		handlers.LoginHandler(w, r, db, store)
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

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	// Create a test login request with invalid credentials
	payload := strings.NewReader("email=test@example.com&password=wrongpassword")
	req, err := http.NewRequest("POST", "/login", payload)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, db, store)
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
	handlers.CreateListHandler(recorder, req, db, store)

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

func TestCheckCredentials_Performance(t *testing.T) {
	start := time.Now()

	// Initialize the database connection
	db, err := sql.Open("mysql", "root:w8-!oY4-taa630-lsKnW0ut@tcp(localhost:3306)/shopping_list_app")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)

	_, err = userRepo.ValidateCredentials("test@example.com", "password123")
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

	// Prepare the test response recorder
	recorder := httptest.NewRecorder()

	// Call the handler function
	rr := httptest.NewRecorder()
	handlers.CreateListHandler(recorder, req, db, store)

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

	// Call the handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, db, store)
	})

	handler.ServeHTTP(rr, req)

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

	// Call the handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, db)
	})

	handler.ServeHTTP(rr, req)

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

	// Call the handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ViewListsHandler(w, r, db, store)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	duration := time.Since(start)
	t.Logf("Execution time: %s", duration)
}

func TestIsValidEmail(t *testing.T) {
	// Valid email address
	validEmail := "test@example.com"
	assert.True(t, services.IsValidEmail(validEmail))

	// Invalid email addresses
	invalidEmails := []string{
		"test@example",
		"test.example.com",
		"test@",
		"@example.com",
	}
	for _, email := range invalidEmails {
		assert.False(t, services.IsValidEmail(email))
	}
}
