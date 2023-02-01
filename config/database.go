package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// Connect to MySQL database and return a pointer to sql.DB
func Connect() *sql.DB {
	db, err := sql.Open("mysql", "dauqu:7388139606@tcp(localhost:3306)/dauqu")
	if err != nil {
		fmt.Println(err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connected to database")
	}
	return db
}

// Close the database connection
func Close(db *sql.DB) {
	db.Close()
}
