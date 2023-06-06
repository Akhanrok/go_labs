package database_repository

import (
	"database/sql"
)

var db *sql.DB

// Create a new database connection
func NewDatabase(dataSourceName string) (*sql.DB, error) {
	// Initialize the database connection
	database, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	db = database
	return db, nil
}
