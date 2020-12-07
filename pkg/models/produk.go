package models

import (
	"zonart/db"
)

// Produk is class
type Produk struct {
	IDProduk      int        `json:"idProduk"`
	IDToko        int        `json:"idToko"`
	NamaProduk    string     `json:"namaProduk" validate:"required"`
	Gambar        string     `json:"gambar" validate:"required"`
	Deskripsi     string     `json:"deskripsi"`
	HargaCetak    string     `json:"hargaCetak"`
	HargaSoftCopy string     `json:"hargaSoftCopy"`
	Status        bool       `json:"status"`
	Catatan       string     `json:"catatan"`
	HargaWajah    int        `json:"hargaWajah" validate:"required"`
	GrupOpsi      []GrupOpsi `json:"grupOpsi"`
}

// Produks is list of produk
type Produks struct {
	Produks []Produk `json:"produk"`
}

// GetProduks is func
func (p Produk) GetProduks(idToko string) Produks {
	con := db.Connect()
	query := "SELECT idProduk, idToko, namaProduk, gambar, deskripsi, hargaCetak, hargaSoftCopy, status, catatan, hargaWajah FROM produk WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var produks Produks

	for rows.Next() {
		rows.Scan(
			&p.IDProduk, &p.IDToko, &p.NamaProduk, &p.Gambar, &p.Deskripsi, &p.HargaCetak, &p.HargaSoftCopy, &p.Status, &p.Catatan, &p.HargaWajah,
		)

		produks.Produks = append(produks.Produks, p)
	}

	defer con.Close()

	return produks
}

// GetProduk is func
func (p Produk) GetProduk(idToko, idProduk string) (Produk, error) {
	con := db.Connect()
	query := "SELECT idProduk, idToko, namaProduk, gambar, deskripsi, hargaCetak, hargaSoftCopy, status, catatan, hargaWajah FROM produk WHERE idToko = ? AND idProduk = ?"

	err := con.QueryRow(query, idToko, idProduk).Scan(
		&p.IDProduk, &p.IDToko, &p.NamaProduk, &p.Gambar, &p.Deskripsi, &p.HargaCetak, &p.HargaSoftCopy, &p.Status, &p.Catatan, &p.HargaWajah)

	defer con.Close()

	return p, err
}
