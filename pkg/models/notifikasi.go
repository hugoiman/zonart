package models

import (
	"time"
	"zonart/db"
)

// Notifikasi is class
type Notifikasi struct {
	IDNotifikasi int    `json:"idNotifikasi"`
	IDPenerima   []int  `json:"idPenerima"`
	Pengirim     string `json:"pengirim"`
	Judul        string `json:"judul"`
	Pesan        string `json:"pesan"`
	Link         string `json:"link"`
	Dibaca       bool   `json:"dibaca"`
	CreatedAt    string `json:"createdAt"`
}

// Notifikasis is list of notifikasi
type Notifikasis struct {
	Notifikasis []Notifikasi `json:"notifikasi"`
}

// GetNotifikasis is func
func (n Notifikasi) GetNotifikasis(idCustomer string) Notifikasis {
	con := db.Connect()
	query := "SELECT idNotifikasi, pengirim, judul, pesan, link, dibaca, createdAt FROM notifikasi WHERE idPenerima = ?"
	rows, _ := con.Query(query, idCustomer)

	var notifikasis Notifikasis
	var createdAt time.Time

	for rows.Next() {
		rows.Scan(
			&n.IDNotifikasi, &n.Pengirim, &n.Judul, &n.Pesan, &n.Link, &n.Dibaca, &createdAt,
		)

		n.CreatedAt = createdAt.Format("02 Jan 2006")
		notifikasis.Notifikasis = append(notifikasis.Notifikasis, n)
	}

	defer con.Close()

	return notifikasis
}

// CreateNotifikasi is func
func (n Notifikasi) CreateNotifikasi() {
	con := db.Connect()
	query := "INSERT INTO notifikasi (idPenerima, pengirim, judul, pesan, link, createdAt) VALUES (?,?,?,?,?,?)"

	for _, vIDPenerima := range n.IDPenerima {
		_, _ = con.Exec(query, vIDPenerima, n.Pengirim, n.Judul, n.Pesan, n.Link, n.CreatedAt)
	}

	defer con.Close()
}

// ReadNotifikasi is func
func (n Notifikasi) ReadNotifikasi(idNotifikasi, idPenerima string) error {
	con := db.Connect()
	query := "UPDATE notifikasi SET dibaca = 1 WHERE idNotifikasi = ? AND idPenerima = ?"
	_, err := con.Exec(query, idNotifikasi, idPenerima)

	defer con.Close()

	return err
}
