package tests

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Akhanrok/go_labs/handlers"
)

func TestCheckCredentials(t *testing.T) {
	// Create a test case with inputs and expected output
	testCases := []struct {
		email    string
		password string
		expected string
	}{
		{email: "test@example.com", password: "password", expected: "Test User"},
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		// Call the function to be tested
		result := checkCredentials(tc.email, tc.password)

		// Compare the result with the expected output
		if result != tc.expected {
			t.Errorf("For email=%s, password=%s, expected=%s, but got %s", tc.email, tc.password, tc.expected, result)
		}
	}
}

func TestIsEmailExists(t *testing.T) {
	// Create a test database connection
	db, err := sql.Open("mysql", "test-connection-string")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Call the isEmailExists function with a test email
	email := "test@example.com"
	exists, err := isEmailExists(db, email)
	if err != nil {
		t.Fatalf("Failed to check email existence: %v", err)
	}

	// Check the result of isEmailExists
	if exists {
		t.Errorf("Expected email %s to not exist, but it exists", email)
	}
}

func TestIsListExists(t *testing.T) {
	// Create a test database connection
	db, err := sql.Open("mysql", "test-connection-string")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Call the isListExists function with a test list name
	listName := "Test List"
	exists, err := isListExists(db, listName)
	if err != nil {
		t.Fatalf("Failed to check list existence: %v", err)
	}

	// Check the result of isListExists
	if exists {
		t.Errorf("Expected list %s to not exist, but it exists", listName)
	}
}
