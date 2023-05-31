package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Akhanrok/go_labs/database"
	"github.com/Akhanrok/go_labs/models"
)

type UserHandler struct {
	db *database.Database
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{db: db}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "index.html", nil)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		// Check the credentials in the database
		username, err := checkCredentials(db, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if username == "" {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "Wrong credentials",
			}
			renderTemplate(w, "login.html", data)
			return
		}

		data := struct {
			Username string
		}{
			Username: username,
		}

		renderTemplate(w, "login-success.html", data)
		return
	}

	renderTemplate(w, "login.html", nil)
}

func checkCredentials(db *sql.DB, email, password string) (string, error) {
	query := "SELECT name FROM users WHERE email = ? AND password = ?"
	var username string
	err := db.QueryRow(query, email, password).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.PostForm.Get("name")
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		if name == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if !isValidEmail(email) {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		if len(password) < 8 {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "Password should be at least 8 characters long",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		// Check if the email already exists in the database
		emailExists, err := isEmailExists(db, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if emailExists {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "Email already exists",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		// Insert the new user into the database
		insertQuery := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
		_, err = db.Exec(insertQuery, name, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the register success page
		http.Redirect(w, r, "/register-success", http.StatusFound)
		return
	}

	renderTemplate(w, "register.html", nil)
}

func loginSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "login-success.html", nil)
	}
}

func registerSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "register-success.html", nil)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func isValidEmail(email string) bool {
	// Email validation regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func isEmailExists(db *sql.DB, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func validateCredentials(db *sql.DB, email, password string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = ? AND password = ?"
	var count int
	err := db.QueryRow(query, email, password).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

