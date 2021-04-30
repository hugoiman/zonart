package models

import (
	"encoding/json"
	"strconv"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Produk is class
type Produk struct {
	idProduk       int
	namaProduk     string
	berat          int
	gambar         string
	deskripsi      string
	status         string
	catatan        string
	hargaWajah     int
	slug           string
	jenisPemesanan []JenisPemesananProduk
	grupOpsi       []GrupOpsi
}

func (p *Produk) GetNamaProduk() string {
	return p.namaProduk
}

func (p *Produk) SetGambar(data string) {
	p.gambar = data
}

func (p *Produk) GetBerat() int {
	return p.berat
}

func (p *Produk) GetGambar() string {
	return p.gambar
}

func (p *Produk) GetStatus() string {
	return p.status
}

func (p *Produk) GetHargaWajah() int {
	return p.hargaWajah
}

func (p *Produk) SetSlug(data string) {
	p.slug = data
}

func (p *Produk) GetSlug() string {
	return p.slug
}

func (p *Produk) GetGrupOpsi() []GrupOpsi {
	return p.grupOpsi
}

func (p *Produk) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDProduk       int                    `json:"idProduk"`
		NamaProduk     string                 `json:"namaProduk"`
		Berat          int                    `json:"berat"`
		Gambar         string                 `json:"gambar"`
		Deskripsi      string                 `json:"deskripsi"`
		Status         string                 `json:"status"`
		Catatan        string                 `json:"catatan"`
		HargaWajah     int                    `json:"hargaWajah"`
		Slug           string                 `json:"slug"`
		JenisPemesanan []JenisPemesananProduk `json:"jenisPemesanan"`
		GrupOpsi       []GrupOpsi             `json:"grupOpsi"`
	}{
		IDProduk:       p.idProduk,
		NamaProduk:     p.namaProduk,
		Berat:          p.berat,
		Gambar:         p.gambar,
		Deskripsi:      p.deskripsi,
		Status:         p.status,
		Catatan:        p.catatan,
		HargaWajah:     p.hargaWajah,
		Slug:           p.slug,
		JenisPemesanan: p.jenisPemesanan,
		GrupOpsi:       p.grupOpsi,
	})
}

func (p *Produk) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDProduk       int                    `json:"idProduk"`
		NamaProduk     string                 `json:"namaProduk" validate:"required"`
		Berat          int                    `json:"berat"`
		Gambar         string                 `json:"gambar"`
		Deskripsi      string                 `json:"deskripsi"`
		Status         string                 `json:"status" validate:"required,eq=aktif|eq=tidak aktif"`
		Catatan        string                 `json:"catatan"`
		HargaWajah     int                    `json:"hargaWajah"`
		Slug           string                 `json:"slug"`
		JenisPemesanan []JenisPemesananProduk `json:"jenisPemesanan" validate:"len=2,dive"`
		GrupOpsi       []GrupOpsi             `json:"grupOpsi"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.idProduk = alias.IDProduk
	p.namaProduk = alias.NamaProduk
	p.berat = alias.Berat
	p.gambar = alias.Gambar
	p.deskripsi = alias.Deskripsi
	p.status = alias.Status
	p.catatan = alias.Catatan
	p.hargaWajah = alias.HargaWajah
	p.slug = alias.Slug
	p.jenisPemesanan = alias.JenisPemesanan
	p.grupOpsi = alias.GrupOpsi
	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetProduks is func
func (p Produk) GetProduks(idToko string) []Produk {
	con := db.Connect()
	query := "SELECT idProduk, namaProduk, gambar, deskripsi, berat, status, catatan, hargaWajah, slug FROM produk WHERE idToko = ? AND status != 'dihapus' ORDER BY idProduk DESC"
	rows, _ := con.Query(query, idToko)

	var produks []Produk

	for rows.Next() {
		rows.Scan(
			&p.idProduk, &p.namaProduk, &p.gambar, &p.deskripsi, &p.berat, &p.status, &p.catatan, &p.hargaWajah, &p.slug,
		)

		var jpp JenisPemesananProduk
		p.jenisPemesanan = jpp.GetJenisPemesananProduk(strconv.Itoa(p.idProduk))

		var grupOpsi GrupOpsi
		p.grupOpsi = grupOpsi.GetGrupOpsiProduk(idToko, strconv.Itoa(p.idProduk))

		produks = append(produks, p)
	}

	defer con.Close()

	return produks
}

// GetProduk is func
func (p Produk) GetProduk(idToko, idProduk string) (Produk, error) {
	con := db.Connect()
	query := "SELECT idProduk, namaProduk, gambar, deskripsi, berat, status, catatan, hargaWajah, slug FROM produk WHERE idToko = ? AND (idProduk = ? OR slug = ?)"

	err := con.QueryRow(query, idToko, idProduk, idProduk).Scan(
		&p.idProduk, &p.namaProduk, &p.gambar, &p.deskripsi, &p.berat, &p.status, &p.catatan, &p.hargaWajah, &p.slug)

	var jpp JenisPemesananProduk
	p.jenisPemesanan = jpp.GetJenisPemesananProduk(strconv.Itoa(p.idProduk))

	var grupOpsi GrupOpsi
	p.grupOpsi = grupOpsi.GetGrupOpsiProduk(idToko, strconv.Itoa(p.idProduk))

	defer con.Close()

	return p, err
}

// CreateProduk is func
func (p Produk) CreateProduk(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO produk (idToko, namaProduk, gambar, deskripsi, berat, status, catatan, hargaWajah, slug) VALUES (?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, p.namaProduk, p.gambar, p.deskripsi, &p.berat, p.status, p.catatan, p.hargaWajah, p.slug)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idProduk := int(idInt64)

	for _, v := range p.jenisPemesanan {
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
	_, err := con.Exec(query, p.namaProduk, p.gambar, p.deskripsi, &p.berat, p.status, p.catatan, p.hargaWajah, p.slug, idToko, idProduk)

	for _, v := range p.jenisPemesanan {
		_ = v.UpdateJenisPemesanan(idProduk, strconv.Itoa(v.idJenisPemesanan))
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
