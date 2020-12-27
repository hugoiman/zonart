package models

import (
	"zonart/db"
)

// Customer is class
type Customer struct {
	IDCustomer int    `json:"idCustomer"`
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Nama       string `json:"nama" validate:"required"`
}

// Customers is func
// type Customers struct {
// 	Customers []Customer `json:"customer"`
// }

// GetCustomer is func
func (c Customer) GetCustomer(id string) (Customer, error) {
	con := db.Connect()
	query := "SELECT idCustomer, username, email, nama FROM customer WHERE idCustomer = ? OR email = ? OR username = ?"

	err := con.QueryRow(query, id, id, id).Scan(
		&c.IDCustomer, &c.Username, &c.Email, &c.Nama)

	defer con.Close()

	return c, err
}

// // GetCustomers is func
// func (c Customer) GetCustomers() Customers {
// 	con := db.Connect()
// 	query := "SELECT idCustomer, username, email, nama FROM customer"
// 	rows, _ := con.Query(query)

// 	var customer Customer
// 	var customers Customers

// 	for rows.Next() {
// 		rows.Scan(
// 			&customer.IDCustomer, &customer.Username, &customer.Email, &customer.Nama,
// 		)

// 		customers.Customers = append(customers.Customers, customer)
// 	}

// 	defer con.Close()

// 	return customers
// }

// Register is func
func (c Customer) Register(username, email, nama, password string) error {
	con := db.Connect()
	query := "INSERT INTO customer (username, email, nama, password) VALUES (?,?,?,?)"
	_, err := con.Exec(query, username, email, nama, password)

	defer con.Close()

	return err
}

// UpdatePassword is func
func (c Customer) UpdatePassword(idCustomer int, newPassword string) error {
	con := db.Connect()
	query := "UPDATE customer SET password = ? WHERE idCustomer = ?"
	_, err := con.Exec(query, newPassword, idCustomer)

	defer con.Close()

	return err
}

// UpdateProfil is func
func (c Customer) UpdateProfil(idCustomer int) error {
	con := db.Connect()
	query := "UPDATE customer SET username = ?, email = ?, nama = ? WHERE idCustomer = ?"
	_, err := con.Exec(query, c.Username, c.Email, c.Nama, idCustomer)

	defer con.Close()

	return err
}
