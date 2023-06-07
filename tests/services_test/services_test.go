package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Akhanrok/go_labs/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

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

func TestRenderTemplate(t *testing.T) {
	// Create HTTP response writer
	w := httptest.NewRecorder()

	// Define the template name and data
	tmpl := "example.html"
	data := struct {
		Title string
	}{
		Title: "Example Title",
	}

	// Call the function
	services.RenderTemplate(w, tmpl, data)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
