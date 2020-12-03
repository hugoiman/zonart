package models

import (
	"zonart/db"
)

// Auth is class
type Auth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login is func
func (a *Auth) Login() (string, error) {
	var idCustomer string

	con := db.Connect()
	query := "SELECT idCustomer FROM customer WHERE username = ? AND password = ?"
	err := con.QueryRow(query, a.Username, a.Password).Scan(&idCustomer)

	defer con.Close()

	return idCustomer, err
}

// CheckOldPassword is func
func CheckOldPassword(idCustomer int, password string) bool {
	var isAny bool
	con := db.Connect()
	query := "SELECT EXISTS (SELECT 1 FROM customer WHERE idCustomer = ? AND password = ?)"
	con.QueryRow(query, idCustomer, password).Scan(&isAny)

	defer con.Close()

	return isAny
}
