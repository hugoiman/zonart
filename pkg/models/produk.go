package models

import (
	"zonart/db"
)

// Produk is class
type Produk struct {
	IDProduk      int        `json:"idProduk"`
	IDToko        int        `json:"idToko"`
	NamaProduk    string     `json:"namaProduk" validate:"required"`
	Cetak         bool       `json:"cetak"`
	SoftCopy      bool       `json:"softCopy"`
	Gambar        string     `json:"gambar" validate:"required"`
	Deskripsi     string     `json:"deskripsi"`
	HargaCetak    int        `json:"hargaCetak"`
	HargaSoftCopy int        `json:"hargaSoftCopy"`
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
	query := "SELECT idProduk, idToko, namaProduk, gambar, deskripsi, cetak, softCopy, hargaCetak, hargaSoftCopy, status, catatan, hargaWajah FROM produk WHERE idToko = ?"
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
	query := "SELECT idProduk, idToko, namaProduk, gambar, deskripsi, cetak, softCopy, hargaCetak, hargaSoftCopy, status, catatan, hargaWajah FROM produk WHERE idToko = ? AND idProduk = ?"

	err := con.QueryRow(query, idToko, idProduk).Scan(
		&p.IDProduk, &p.IDToko, &p.NamaProduk, &p.Gambar, &p.Deskripsi, &p.Cetak, &p.SoftCopy, &p.HargaCetak, &p.HargaSoftCopy, &p.Status, &p.Catatan, &p.HargaWajah)

	var grupOpsi GrupOpsi
	p.GrupOpsi = grupOpsi.GetGrupOpsiProduk(idToko, idProduk)

	defer con.Close()

	return p, err
}

// CreateProduk is func
func (p Produk) CreateProduk(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO produk (idToko, namaProduk, gambar, deskripsi, cetak, softCopy, hargaCetak, hargaSoftCopy, status, catatan, hargaWajah) VALUES (?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, p.NamaProduk, p.Gambar, p.Deskripsi, &p.Cetak, &p.SoftCopy, p.HargaCetak, p.HargaSoftCopy, p.Status, p.Catatan, p.HargaWajah)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idProduk := int(idInt64)

	defer con.Close()

	return idProduk, err
}

// UpdateProduk is func
func (p Produk) UpdateProduk(idToko, idProduk string) error {
	con := db.Connect()
	query := "UPDATE produk SET namaProduk = ?, gambar = ?, deskripsi = ?, cetak = ?, softCopy = ?, hargaCetak = ?, hargaSoftCopy = ?, status = ?, catatan = ?, hargaWajah = ? WHERE idToko = ? AND idProduk = ?"
	_, err := con.Exec(query, p.NamaProduk, p.Gambar, p.Deskripsi, &p.Cetak, &p.SoftCopy, p.HargaCetak, p.HargaSoftCopy, p.Status, p.Catatan, p.HargaWajah, idToko, idProduk)

	defer con.Close()

	return err
}

// DeleteProduk is func
func (p Produk) DeleteProduk(idToko, idProduk string) error {
	con := db.Connect()
	query := "DELETE FROM produk WHERE idToko = ? AND idProduk = ?"
	_, err := con.Exec(query, idToko, idProduk)

	defer con.Close()

	return err
}
