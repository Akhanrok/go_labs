package user_handlers

import (
	"database/sql"
	"net/http"

	"github.com/Akhanrok/go_labs/repositories/user_repository"
	"github.com/Akhanrok/go_labs/services"
	"github.com/gorilla/sessions"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		services.RenderTemplate(w, "index.html", nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, store sessions.Store) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		// Create instances of the repositories
		userRepo := user_repository.NewUserRepository(db)

		// Check the credentials in the database
		username, err := userRepo.ValidateCredentials(email, password)
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
			services.RenderTemplate(w, "login.html", data)
			return
		}

		// Create a session for the authenticated user
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the username in the session
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Username string
		}{
			Username: username,
		}

		services.RenderTemplate(w, "login-success.html", data)
		return
	}

	services.RenderTemplate(w, "login.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

		if !services.IsValidEmail(email) {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		if len(password) < 8 {
			data := struct {
				ErrorMessage string
			}{
				ErrorMessage: "Password should be at least 8 characters long",
			}
			services.RenderTemplate(w, "register.html", data)
			return
		}

		// Create instances of the repositories
		userRepo := user_repository.NewUserRepository(db)

		// Check if the email already exists in the database
		emailExists, err := userRepo.IsEmailExists(email)
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
			services.RenderTemplate(w, "register.html", data)
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

	services.RenderTemplate(w, "register.html", nil)
}

func LoginSuccessHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		services.RenderTemplate(w, "login-success.html", nil)
	}
}

func RegisterSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		services.RenderTemplate(w, "register-success.html", nil)
	}
}
