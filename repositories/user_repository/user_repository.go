package user_repository

import (
	"database/sql"
)

type UserRepository interface {
	ValidateCredentials(email, password string) (string, error)
	IsEmailExists(email string) (bool, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) ValidateCredentials(email, password string) (string, error) {
	query := "SELECT name FROM users WHERE email = ? AND password = ?"
	var username string
	err := r.db.QueryRow(query, email, password).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (r *userRepository) IsEmailExists(email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
