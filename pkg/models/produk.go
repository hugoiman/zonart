package models

import (
	"strconv"
	"zonart/db"
)

// Produk is class
type Produk struct {
	IDProduk       int                    `json:"idProduk"`
	NamaProduk     string                 `json:"namaProduk" validate:"required"`
	Berat          int                    `json:"berat"`
	Gambar         string                 `json:"gambar"`
	Deskripsi      string                 `json:"deskripsi"`
	Status         string                 `json:"status" validate:"required,eq=aktif|eq=tidak aktif"`
	Catatan        string                 `json:"catatan"`
	HargaWajah     int                    `json:"hargaWajah"`
	Slug           string                 `json:"slug"`
	JenisPemesanan []JenisPemesananProduk `json:"jenisPemesanan" validate:"unique=IDJenisPemesanan,len=2,dive"`
	GrupOpsi       []GrupOpsi             `json:"grupOpsi"`
}

// Produks is list of produk
type Produks struct {
	Produks []Produk `json:"produk"`
}

// GetProduks is func
func (p Produk) GetProduks(idToko string) Produks {
	con := db.Connect()
	query := "SELECT idProduk, namaProduk, gambar, deskripsi, berat, status, catatan, hargaWajah, slug FROM produk WHERE idToko = ? AND status != 'dihapus' ORDER BY idProduk DESC"
	rows, _ := con.Query(query, idToko)

	var produks Produks

	for rows.Next() {
		rows.Scan(
			&p.IDProduk, &p.NamaProduk, &p.Gambar, &p.Deskripsi, &p.Berat, &p.Status, &p.Catatan, &p.HargaWajah, &p.Slug,
		)

		var jpp JenisPemesananProduk
		p.JenisPemesanan = jpp.GetJenisPemesananProduk(strconv.Itoa(p.IDProduk))

		var grupOpsi GrupOpsi
		p.GrupOpsi = grupOpsi.GetGrupOpsiProduk(idToko, strconv.Itoa(p.IDProduk))

		produks.Produks = append(produks.Produks, p)
	}

	defer con.Close()

	return produks
}

// GetProduk is func
func (p Produk) GetProduk(idToko, idProduk string) (Produk, error) {
	con := db.Connect()
	query := "SELECT idProduk, namaProduk, gambar, deskripsi, berat, status, catatan, hargaWajah, slug FROM produk WHERE idToko = ? AND (idProduk = ? OR slug = ?)"

	err := con.QueryRow(query, idToko, idProduk, idProduk).Scan(
		&p.IDProduk, &p.NamaProduk, &p.Gambar, &p.Deskripsi, &p.Berat, &p.Status, &p.Catatan, &p.HargaWajah, &p.Slug)

	var jpp JenisPemesananProduk
	p.JenisPemesanan = jpp.GetJenisPemesananProduk(strconv.Itoa(p.IDProduk))

	var grupOpsi GrupOpsi
	p.GrupOpsi = grupOpsi.GetGrupOpsiProduk(idToko, strconv.Itoa(p.IDProduk))

	defer con.Close()

	return p, err
}

// CreateProduk is func
func (p Produk) CreateProduk(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO produk (idToko, namaProduk, gambar, deskripsi, berat, status, catatan, hargaWajah, slug) VALUES (?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, p.NamaProduk, p.Gambar, p.Deskripsi, &p.Berat, p.Status, p.Catatan, p.HargaWajah, p.Slug)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idProduk := int(idInt64)

	for _, v := range p.JenisPemesanan {
		err = v.CreateJenisPemesanan(strconv.Itoa(idProduk))
		if err != nil {
			return idProduk, err
		}
	}

	defer con.Close()

	return idProduk, err
}

// func (p Produk) CreateProduk(idToko string) (int, error) {
// 	return 1, nil
// 	return 0, errors.New("Terjadi Error")
// }

// UpdateProduk is func
func (p Produk) UpdateProduk(idToko, idProduk string) error {
	con := db.Connect()
	query := "UPDATE produk SET namaProduk = ?, gambar = ?, deskripsi = ?, berat = ?, status = ?, catatan = ?, hargaWajah = ?, slug = ? WHERE idToko = ? AND idProduk = ?"
	_, err := con.Exec(query, p.NamaProduk, p.Gambar, p.Deskripsi, &p.Berat, p.Status, p.Catatan, p.HargaWajah, p.Slug, idToko, idProduk)

	for _, v := range p.JenisPemesanan {
		_ = v.UpdateJenisPemesanan(idProduk, strconv.Itoa(v.IDJenisPemesanan))
	}
	defer con.Close()

	return err
}

// DeleteProduk is func
func (p Produk) DeleteProduk(idToko, idProduk string) error {
	con := db.Connect()
	query := "UPDATE produk SET status = 'dihapus' WHERE idToko = ? AND idProduk = ?"
	_, err := con.Exec(query, idToko, idProduk)

	defer con.Close()

	return err
}
