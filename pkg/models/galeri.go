package models

import (
	"zonart/db"
)

// Galeri is class
type Galeri struct {
	IDGaleri int    `json:"idGaleri"`
	IDToko   int    `json:"idToko"`
	IDProduk int    `json:"idProduk" validate:"required"`
	Kategori string `json:"kategori"`
	Gambar   string `json:"gambar" validate:"required"`
}

// Galeris is list of galeri
type Galeris struct {
	Galeris []Galeri `json:"galeri"`
}

// GetGaleris is func
func (g Galeri) GetGaleris(idToko string) Galeris {
	con := db.Connect()
	query := "SELECT a.idGaleri, a.idToko, a.idProduk, b.namaProduk, a.gambar FROM galeri a JOIN produk b ON a.idProduk = b.idProduk WHERE a.idToko = ? ORDER BY a.idGaleri DESC"
	rows, _ := con.Query(query, idToko)

	var galeris Galeris

	for rows.Next() {
		rows.Scan(
			&g.IDGaleri, &g.IDToko, &g.IDProduk, &g.Kategori, &g.Gambar,
		)

		galeris.Galeris = append(galeris.Galeris, g)
	}

	defer con.Close()

	return galeris
}

// CreateGaleri is func
func (g Galeri) CreateGaleri(idToko string) error {
	con := db.Connect()
	query := "INSERT INTO galeri (idToko, idProduk, gambar) VALUES (?,?,?)"
	_, err := con.Exec(query, idToko, g.IDProduk, g.Gambar)

	defer con.Close()

	return err
}

// DeleteGaleri is func
func (g Galeri) DeleteGaleri(idToko, idGaleri string) error {
	con := db.Connect()
	query := "DELETE FROM galeri WHERE idToko = ? AND idGaleri = ?"
	_, err := con.Exec(query, idToko, idGaleri)

	defer con.Close()

	return err
}
