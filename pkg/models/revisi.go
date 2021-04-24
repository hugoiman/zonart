package models

import (
	"time"
	"zonart/db"
)

// Revisi is class
type Revisi struct {
	IDRevisi  int    `json:"idRevisi"`
	IDOrder   int    `json:"idOrder"`
	Catatan   string `json:"catatan" validate:"required"`
	CreatedAt string `json:"createdAt"`
}

// GetRevisi is func
func (r Revisi) GetRevisi(idOrder string) []Revisi {
	con := db.Connect()
	var revisis []Revisi
	var createdAt time.Time

	query := "SELECT idRevisi, idOrder, catatan, createdAt FROM revisi WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&r.IDRevisi, &r.IDOrder, &r.Catatan, &createdAt,
		)

		r.CreatedAt = createdAt.Format("02 Jan 2006")
		revisis = append(revisis, r)
	}

	defer con.Close()

	return revisis
}

// CreateRevisi is func
func (r Revisi) CreateRevisi(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO revisi (idOrder, catatan, createdAt) VALUES (?,?,?)"
	_, err := con.Exec(query, idOrder, r.Catatan, r.CreatedAt)

	defer con.Close()

	return err
}
