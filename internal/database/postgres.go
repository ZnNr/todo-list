package database

import (
	"database/sql"
	"fmt"
)

func NewPostgresDB() (*sql.DB, error) {
	host := "localhost"
	port := "5432"
	user := "postgres"
	password := "postgres"
	dbname := "postgres"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}
