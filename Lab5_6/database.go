package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	os.Setenv("CGO_ENABLED", "1")
}

func main() {
	// Open MySQL database connection
	db, err := sql.Open("mysql", "root:root123@tcp(localhost:3306)/user_management")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Create 'users' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			avatar VARCHAR(255) NOT NULL,
			phonenumber VARCHAR(255) NOT NULL,
			dateofbirth DATE NOT NULL,
			address VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	fmt.Println("Table created successfully")
}
