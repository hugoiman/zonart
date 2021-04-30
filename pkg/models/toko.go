package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Toko is class
type Toko struct {
	idToko             int
	owner              int
	namaToko           string
	emailToko          string
	deskripsi          string
	alamat             string
	kota               string
	telp               string
	whatsapp           string
	instagram          string
	website            string
	slug               string
	foto               string
	createdAt          string
	jasaPengirimanToko []JasaPengirimanToko
	rekening           []Rekening
}

func (t *Toko) SetOwner(data int) {
	t.owner = data
}

func (t *Toko) GetOwner() int {
	return t.owner
}

func (t *Toko) GetNamaToko() string {
	return t.namaToko
}

func (t *Toko) GetKota() string {
	return t.kota
}

func (t *Toko) GetSlug() string {
	return t.slug
}

func (t *Toko) SetFoto(data string) {
	t.foto = data
}

func (t *Toko) GetFoto() string {
	return t.foto
}

func (t *Toko) SetCreatedAt(data string) {
	t.createdAt = data
}

func (t *Toko) GetJasaPengirimanToko() []JasaPengirimanToko {
	return t.jasaPengirimanToko
}

func (t *Toko) GetRekening() []Rekening {
	return t.rekening
}

func (t *Toko) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDToko             int                  `json:"idToko"`
		Owner              int                  `json:"owner"`
		NamaToko           string               `json:"namaToko"`
		EmailToko          string               `json:"emailToko"`
		Deskripsi          string               `json:"deskripsi"`
		Alamat             string               `json:"alamat"`
		Kota               string               `json:"kota"`
		Telp               string               `json:"telp"`
		Whatsapp           string               `json:"whatsapp"`
		Instagram          string               `json:"instagram"`
		Website            string               `json:"website"`
		Slug               string               `json:"slug"`
		Foto               string               `json:"foto"`
		CreatedAt          string               `json:"createdAt"`
		JasaPengirimanToko []JasaPengirimanToko `json:"jasaPengirimanToko"`
		Rekening           []Rekening           `json:"rekening"`
	}{
		IDToko:             t.idToko,
		Owner:              t.owner,
		NamaToko:           t.namaToko,
		EmailToko:          t.emailToko,
		Deskripsi:          t.deskripsi,
		Alamat:             t.alamat,
		Kota:               t.kota,
		Telp:               t.telp,
		Whatsapp:           t.whatsapp,
		Instagram:          t.instagram,
		Website:            t.website,
		Slug:               t.slug,
		Foto:               t.foto,
		CreatedAt:          t.createdAt,
		JasaPengirimanToko: t.jasaPengirimanToko,
		Rekening:           t.rekening,
	})
}

func (t *Toko) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDToko             int                  `json:"idToko"`
		Owner              int                  `json:"owner"`
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
		Rekening           []Rekening           `json:"rekening" validate:"dive"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	t.idToko = alias.IDToko
	t.owner = alias.Owner
	t.namaToko = alias.NamaToko
	t.emailToko = alias.EmailToko
	t.deskripsi = alias.Deskripsi
	t.alamat = alias.Alamat
	t.kota = alias.Kota
	t.telp = alias.Telp
	t.whatsapp = alias.Whatsapp
	t.instagram = alias.Instagram
	t.website = alias.Website
	t.slug = alias.Slug
	t.foto = alias.Foto
	t.createdAt = alias.CreatedAt
	t.jasaPengirimanToko = alias.JasaPengirimanToko
	t.rekening = alias.Rekening

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetToko is func
func (t Toko) GetToko(id string) (Toko, error) {
	con := db.Connect()
	query := "SELECT idToko, owner, namaToko, emailToko, foto, deskripsi, alamat, kota, telp, whatsapp, instagram, website, slug, createdAt FROM toko WHERE idToko = ? OR slug = ?"

	var createdAt time.Time

	err := con.QueryRow(query, id, id).Scan(
		&t.idToko, &t.owner, &t.namaToko, &t.emailToko, &t.foto, &t.deskripsi, &t.alamat, &t.kota,
		&t.telp, &t.whatsapp, &t.instagram, &t.website, &t.slug, &createdAt)

	t.createdAt = createdAt.Format("02 Jan 2006")

	var jpt JasaPengirimanToko
	dataJasaPengirimanToko := jpt.GetJasaPengirimanToko(t.idToko)
	t.jasaPengirimanToko = dataJasaPengirimanToko

	var r Rekening
	dataRekeningToko := r.GetRekening(id)
	t.rekening = dataRekeningToko

	defer con.Close()

	return t, err
}

// CreateToko is func
func (t Toko) CreateToko() (int, error) {
	con := db.Connect()
	query, _ := con.Prepare("INSERT INTO toko (owner, namaToko, emailToko, foto, deskripsi, alamat, kota, telp, whatsapp, instagram, website, slug, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
	exec, err := query.Exec(t.owner, t.namaToko, t.emailToko, t.foto, t.deskripsi, t.alamat,
		t.kota, t.telp, t.whatsapp, t.instagram, t.website, t.slug, t.createdAt)
	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idToko := int(idInt64)

	defer con.Close()

	return idToko, err
}

func (t Toko) GetIDTokoByOrder(idOrder string) (string, error) {
	var idToko string

	con := db.Connect()
	query := "SELECT idToko FROM `order` WHERE idOrder = ?"
	err := con.QueryRow(query, idOrder).Scan(&idToko)

	defer con.Close()

	return idToko, err
}

func (t Toko) GetIDTokoByUndangan(idUndangan string) (string, error) {
	var idToko string

	con := db.Connect()
	query := "SELECT idToko FROM undangan WHERE idUndangan = ?"
	err := con.QueryRow(query, idUndangan).Scan(&idToko)

	defer con.Close()

	return idToko, err
}

// UpdateToko is func
func (t Toko) UpdateToko(id string) error {
	con := db.Connect()
	query := "UPDATE toko SET namaToko = ?, emailToko = ?, deskripsi = ?, alamat = ?, kota = ?, telp = ?, whatsapp = ?, instagram = ?, website = ?, slug = ?, foto = ? WHERE idToko = ? OR slug = ?"
	_, err := con.Exec(query, t.namaToko, t.emailToko, t.deskripsi, t.alamat,
		t.kota, t.telp, t.whatsapp, t.instagram, t.website, t.slug, t.foto, id, id)

	defer con.Close()

	return err
}

// GetMyToko is func
func (t Toko) GetMyToko(idCustomer string) []Toko {
	con := db.Connect()
	query := "SELECT idToko, namaToko, foto, slug FROM toko WHERE owner = ?"
	rows, _ := con.Query(query, idCustomer)

	var tokos []Toko

	for rows.Next() {
		rows.Scan(
			&t.idToko, &t.namaToko, &t.foto, &t.slug,
		)
		t.deskripsi = "owner"
		tokos = append(tokos, t)
	}

	defer con.Close()

	return tokos
}

// GetTokoByEmploye is func
func (t Toko) GetTokoByEmploye(idCustomer string) []Toko {
	con := db.Connect()
	query := "SELECT a.idToko, a.namaToko, a.foto, a.slug FROM toko a " +
		"JOIN karyawan b ON a.idToko = b.idToko " +
		"WHERE b.idCustomer = ?"
	rows, _ := con.Query(query, idCustomer)

	var tokos []Toko

	for rows.Next() {
		rows.Scan(
			&t.idToko, &t.namaToko, &t.foto, &t.slug,
		)
		t.deskripsi = "employee"
		tokos = append(tokos, t)
	}

	defer con.Close()

	return tokos
}
