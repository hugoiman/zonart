package db

import (
	"database/sql"
	"fmt"

	// _ is import
	_ "github.com/go-sql-driver/mysql"
)

// Connect is func
func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/zonart?parseTime=true")

	if err != nil {
		fmt.Println("db is not connected")
		panic(err.Error())
	}
	// else {
	// 	fmt.Println("db is connected")
	// }
	return db
}
