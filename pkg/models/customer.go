package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Customer is class
type Customer struct {
	idCustomer int
	username   string
	email      string
	nama       string
}

func (c *Customer) GetIDCustomer() int {
	return c.idCustomer
}

func (c *Customer) GetUsername() string {
	return c.username
}

func (c *Customer) GetEmail() string {
	return c.email
}

func (c *Customer) GetNama() string {
	return c.nama
}

func (c *Customer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDCustomer int    `json:"idCustomer"`
		Username   string `json:"username"`
		Email      string `json:"email"`
		Nama       string `json:"nama"`
	}{
		IDCustomer: c.idCustomer,
		Username:   c.username,
		Email:      c.email,
		Nama:       c.nama,
	})
}

func (g *Customer) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDCustomer int    `json:"idCustomer"`
		Username   string `json:"username" validate:"required,min=3,max=20"`
		Email      string `json:"email" validate:"required,email"`
		Nama       string `json:"nama" validate:"required"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	g.idCustomer = alias.IDCustomer
	g.username = alias.Username
	g.email = alias.Email
	g.nama = alias.Nama

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetCustomer is func
func (c Customer) GetCustomer(id string) (Customer, error) {
	con := db.Connect()
	query := "SELECT idCustomer, username, email, nama FROM customer WHERE idCustomer = ? OR email = ? OR username = ?"

	err := con.QueryRow(query, id, id, id).Scan(
		&c.idCustomer, &c.username, &c.email, &c.nama)

	defer con.Close()

	return c, err
}

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
	_, err := con.Exec(query, c.username, c.email, c.nama, idCustomer)

	defer con.Close()

	return err
}
