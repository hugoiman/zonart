package models

import (
	"zonart/db"
)

// Galeri is class
type Galeri struct {
	IDGaleri  int    `json:"idGaleri"`
	IDToko    int    `json:"idToko"`
	Gambar    string `json:"gambar" validate:"required"`
	Deskripsi string `json:"deskripsi"`
}

// Galeris is list of galeri
type Galeris struct {
	Galeris []Galeri `json:"galeri"`
}

// GetGaleris is func
func (g Galeri) GetGaleris(idToko string) Galeris {
	con := db.Connect()
	query := "SELECT idGaleri, idToko, gambar, deskripsi FROM galeri WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var galeris Galeris

	for rows.Next() {
		rows.Scan(
			&g.IDGaleri, &g.IDToko, &g.Gambar, &g.Deskripsi,
		)

		galeris.Galeris = append(galeris.Galeris, g)
	}

	defer con.Close()

	return galeris
}

// CreateGaleri is func
func (g Galeri) CreateGaleri(idToko string) error {
	con := db.Connect()
	query := "INSERT INTO galeri (idToko, gambar, deskripsi) VALUES (?,?,?)"
	_, err := con.Exec(query, idToko, g.Gambar, g.Deskripsi)

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
