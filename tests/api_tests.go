package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Akhanrok/go_labs/handlers"
)

func TestRegisterHandler(t *testing.T) {
	// Prepare the form data for the register request
	formData := map[string]string{
		"name":     "John",
		"email":    "john@example.com",
		"password": "password",
	}

	jsonData, err := json.Marshal(formData)
	assert.NoError(t, err, "Failed to marshal form data")

	// Create a request to the "/register" endpoint with the form data
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err, "Failed to create request")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the registerHandler function with the response recorder and request
	registerHandler(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusFound, rr.Code, "Unexpected status code")

	// Check the redirection to "/register-success":
	assert.Equal(t, "/register-success", rr.Header().Get("Location"), "Unexpected redirect location")
}

func TestViewListsHandler(t *testing.T) {
	// Create a request to the "/view-lists" endpoint
	req, err := http.NewRequest("GET", "/view-lists", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the viewListsHandler function with the response recorder and request
	viewListsHandler(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code")
}

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}
