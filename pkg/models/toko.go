package models

import (
	"time"
	"zonart/db"
)

// Toko is class
type Toko struct {
	IDToko             int                  `json:"idToko"`
	IDOwner            int                  `json:"idOwner"`
	NamaToko           string               `json:"namaToko" validate:"required"`
	EmailToko          string               `json:"emailToko" validate:"email"`
	Deskripsi          string               `json:"deskripsi"`
	Alamat             string               `json:"alamat"`
	Kota               string               `json:"kota" validate:"required"`
	Telp               string               `json:"telp"`
	Whatsapp           string               `json:"whatsapp"`
	Instagram          string               `json:"instagram"`
	Website            string               `json:"website"`
	Slug               string               `json:"slug" validate:"required"`
	Foto               string               `json:"foto"`
	CreatedAt          string               `json:"createdAt"`
	JasaPengirimanToko []JasaPengirimanToko `json:"jasaPengirimanToko"`
	Rekening           []Rekening           `json:"rekening"`
}

// Tokos is object
type Tokos struct {
	Tokos []Toko `json:"toko"`
}

// GetToko is func
func (t Toko) GetToko(id string) (Toko, error) {
	con := db.Connect()
	query := "SELECT idToko, idOwner, namaToko, emailToko, foto, deskripsi, alamat, kota, telp, whatsapp, instagram, website, slug, createdAt FROM toko WHERE idToko = ? OR slug = ?"

	var createdAt time.Time

	err := con.QueryRow(query, id, id).Scan(
		&t.IDToko, &t.IDOwner, &t.NamaToko, &t.EmailToko, &t.Foto, &t.Deskripsi, &t.Alamat, &t.Kota,
		&t.Telp, &t.Whatsapp, &t.Instagram, &t.Website, &t.Slug, &createdAt)

	t.CreatedAt = createdAt.Format("02 Jan 2006")

	var jpt JasaPengirimanToko
	dataJasaPengirimanToko := jpt.GetJasaPengirimanToko(t.IDToko)
	t.JasaPengirimanToko = dataJasaPengirimanToko

	var r Rekening
	dataRekeningToko := r.GetRekening(id)
	t.Rekening = dataRekeningToko

	defer con.Close()

	return t, err
}

// CreateToko is func
func (t Toko) CreateToko() (int, error) {
	con := db.Connect()
	query, _ := con.Prepare("INSERT INTO toko (idOwner, namaToko, emailToko, foto, deskripsi, alamat, kota, telp, whatsapp, instagram, website, slug, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	exec, err := query.Exec(t.IDOwner, t.NamaToko, t.EmailToko, t.Foto, t.Deskripsi, t.Alamat,
		t.Kota, t.Telp, t.Whatsapp, t.Instagram, t.Website, t.Slug, t.CreatedAt)
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
	query := "UPDATE toko SET namaToko = ?, emailToko = ?, deskripsi = ?, alamat = ?, kota = ?, telp = ?, whatsapp = ?, instagram = ?, website = ?, slug = ?, foto = ? WHERE idToko = ? OR slug = ?"
	_, err := con.Exec(query, t.NamaToko, t.EmailToko, t.Deskripsi, t.Alamat,
		t.Kota, t.Telp, t.Whatsapp, t.Instagram, t.Website, t.Slug, t.Foto, id, id)

	defer con.Close()

	return err
}

// GetMyToko is func
func (t Toko) GetMyToko(idCustomer string) Tokos {
	con := db.Connect()
	query := "SELECT idToko, namaToko, foto, slug FROM toko WHERE idOwner = ?"
	rows, _ := con.Query(query, idCustomer)

	var tokos Tokos

	for rows.Next() {
		rows.Scan(
			&t.IDToko, &t.NamaToko, &t.Foto, &t.Slug,
		)
		t.Deskripsi = "owner"
		tokos.Tokos = append(tokos.Tokos, t)
	}

	defer con.Close()

	return tokos
}

// GetTokoByEmploye is func
func (t Toko) GetTokoByEmploye(idCustomer string) Tokos {
	con := db.Connect()
	query := "SELECT a.idToko, a.namaToko, a.foto, a.slug FROM toko a " +
		"JOIN karyawan b ON a.idToko = b.idToko " +
		"WHERE b.idCustomer = ?"
	rows, _ := con.Query(query, idCustomer)

	var tokos Tokos

	for rows.Next() {
		rows.Scan(
			&t.IDToko, &t.NamaToko, &t.Foto, &t.Slug,
		)
		t.Deskripsi = "employee"
		tokos.Tokos = append(tokos.Tokos, t)
	}

	defer con.Close()

	return tokos
}
