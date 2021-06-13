package db

import (
	"database/sql"
	"fmt"
	"os"

	// _ is import
	_ "github.com/go-sql-driver/mysql"
)

// Connect is func
func Connect() *sql.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+dbname+"?parseTime=true")
	// db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/zonart?parseTime=true")

	if err != nil {
		fmt.Println("db is not connected")
		panic(err.Error())
	}
	// else {
	// 	fmt.Println("db is connected")
	// }
	return db
}
