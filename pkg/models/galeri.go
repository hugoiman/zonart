package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Galeri is class
type Galeri struct {
	idGaleri   int
	idKategori int
	kategori   string
	gambar     string
}

// Galeris is list of galeri
type Galeris struct {
	Galeris []Galeri `json:"galeri"`
}

func (g *Galeri) SetGambar(data string) {
	g.gambar = data
}

func (g *Galeri) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDGaleri   int    `json:"idGaleri"`
		IDKategori int    `json:"idKategori"`
		Kategori   string `json:"kategori"`
		Gambar     string `json:"gambar"`
	}{
		IDGaleri:   g.idGaleri,
		IDKategori: g.idKategori,
		Kategori:   g.kategori,
		Gambar:     g.gambar,
	})
}

func (g *Galeri) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDGaleri   int    `json:"idGaleri"`
		IDKategori int    `json:"idKategori" validate:"required"`
		Kategori   string `json:"kategori"`
		Gambar     string `json:"gambar"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	g.idGaleri = alias.IDGaleri
	g.idKategori = alias.IDKategori
	g.kategori = alias.Kategori
	g.gambar = alias.Gambar

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetGaleris is func
func (g Galeri) GetGaleris(idToko string) Galeris {
	con := db.Connect()
	query := "SELECT a.idGaleri, a.idKategori, b.namaProduk, a.gambar FROM galeri a JOIN produk b ON a.idKategori = b.idProduk WHERE a.idToko = ? ORDER BY a.idGaleri DESC"
	rows, _ := con.Query(query, idToko)

	var galeris Galeris

	for rows.Next() {
		rows.Scan(
			&g.idGaleri, &g.idKategori, &g.kategori, &g.gambar,
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
	exec, err := con.Exec(query, idToko, g.idKategori, g.gambar)

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
