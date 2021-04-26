package models

import (
	"zonart/db"
)

// Galeri is class
type Galeri struct {
	IDGaleri   int    `json:"idGaleri"`
	IDKategori int    `json:"idKategori" validate:"required"`
	Kategori   string `json:"kategori"`
	Gambar     string `json:"gambar"`
}

// Galeris is list of galeri
type Galeris struct {
	Galeris []Galeri `json:"galeri"`
}

// GetGaleris is func
func (g Galeri) GetGaleris(idToko string) Galeris {
	con := db.Connect()
	query := "SELECT a.idGaleri, a.idKategori, b.namaProduk, a.gambar FROM galeri a JOIN produk b ON a.idKategori = b.idProduk WHERE a.idToko = ? ORDER BY a.idGaleri DESC"
	rows, _ := con.Query(query, idToko)

	var galeris Galeris

	for rows.Next() {
		rows.Scan(
			&g.IDGaleri, &g.IDKategori, &g.Kategori, &g.Gambar,
		)

		galeris.Galeris = append(galeris.Galeris, g)
	}

	defer con.Close()

	return galeris
}

// CreateGaleri is func
func (g Galeri) CreateGaleri(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO galeri (idToko, idKategori, gambar) VALUES (?,?,?)"
	exec, err := con.Exec(query, idToko, g.IDKategori, g.Gambar)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idGaleri := int(idInt64)

	defer con.Close()

	return idGaleri, err
}

// DeleteGaleri is func
func (g Galeri) DeleteGaleri(idToko string, idGaleri []string) error {
	con := db.Connect()
	var err error
	for _, v := range idGaleri {
		query := "DELETE FROM galeri WHERE idToko = ? AND idGaleri = ?"
		_, err = con.Exec(query, idToko, v)
	}

	defer con.Close()

	return err
}
