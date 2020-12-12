package models

import (
	"time"
	"zonart/db"
)

// Toko is class
type Toko struct {
	IDToko         int            `json:"idToko"`
	IDOwner        int            `json:"idOwner"`
	NamaToko       string         `json:"namaToko" validate:"required"`
	EmailToko      string         `json:"emailToko" validate:"email"`
	Deskripsi      string         `json:"deskripsi"`
	Alamat         string         `json:"alamat"`
	Kota           string         `json:"kota" validate:"required"`
	Telp           string         `json:"telp"`
	Whatsapp       string         `json:"whatsapp"`
	Instagram      string         `json:"instagram"`
	Website        string         `json:"website"`
	Slug           string         `json:"slug" validate:"required"`
	Foto           string         `json:"foto"`
	SetKaryawan    bool           `json:"setKaryawan"`
	CreatedAt      string         `json:"createdAt"`
	PengirimanToko PengirimanToko `json:"pengirimanToko"`
	Rekening       []Rekening     `json:"rekening"`
}

// GetToko is func
func (t Toko) GetToko(id string) (Toko, error) {
	con := db.Connect()
	query := "SELECT idToko, idOwner, namaToko, emailToko, foto, deskripsi, alamat, kota, telp, whatsapp, instagram, website, slug, setKaryawan, createdAt FROM toko WHERE idToko = ? OR slug = ?"

	var createdAt time.Time

	err := con.QueryRow(query, id, id).Scan(
		&t.IDToko, &t.IDOwner, &t.NamaToko, &t.EmailToko, &t.Foto, &t.Deskripsi, &t.Alamat, &t.Kota,
		&t.Telp, &t.Whatsapp, &t.Instagram, &t.Website, &t.Slug, &t.SetKaryawan, &createdAt)

	t.CreatedAt = createdAt.Format("02 Jan 2006")

	defer con.Close()

	return t, err
}

// CreateToko is func
func (t Toko) CreateToko() (int, error) {
	con := db.Connect()
	query, _ := con.Prepare("INSERT INTO toko (idOwner, namaToko, emailToko, foto, deskripsi, alamat, kota, telp, whatsapp, instagram, website, slug, setKaryawan, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	exec, err := query.Exec(t.IDOwner, t.NamaToko, t.EmailToko, t.Foto, t.Deskripsi, t.Alamat,
		t.Kota, t.Telp, t.Whatsapp, t.Instagram, t.Website, t.Slug, t.SetKaryawan, t.CreatedAt)
	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idToko := int(idInt64)

	defer con.Close()

	return idToko, err
}

// UpdateToko is func
func (t Toko) UpdateToko(id string) error {
	con := db.Connect()
	query := "UPDATE toko SET namaToko = ?, emailToko = ?, deskripsi = ?, alamat = ?, kota = ?, telp = ?, whatsapp = ?, instagram = ?, website = ?, slug = ?, foto = ?, setKaryawan = ? WHERE idToko = ? OR slug = ?"
	_, err := con.Exec(query, t.NamaToko, t.EmailToko, t.Deskripsi, t.Alamat,
		t.Kota, t.Telp, t.Whatsapp, t.Instagram, t.Website, t.Slug, t.Foto, t.SetKaryawan, id, id)

	defer con.Close()

	return err
}
