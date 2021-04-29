package models

import (
	"encoding/json"
	"time"
	"zonart/db"
)

// Notifikasi is class
type Notifikasi struct {
	idNotifikasi int
	penerima     []int
	pengirim     string
	judul        string
	pesan        string
	link         string
	dibaca       bool
	createdAt    string
}

// Notifikasis is list of notifikasi
type Notifikasis struct {
	Notifikasis []Notifikasi `json:"notifikasi"`
}

func (n *Notifikasi) SetPenerima(data []int) {
	n.penerima = data
}

func (n *Notifikasi) GetPenerima() []int {
	return n.penerima
}

func (n *Notifikasi) SetPengirim(data string) {
	n.pengirim = data
}

func (n *Notifikasi) GetPengirim() string {
	return n.pengirim
}

func (n *Notifikasi) SetJudul(data string) {
	n.judul = data
}

func (n *Notifikasi) SetPesan(data string) {
	n.pesan = data
}

func (n *Notifikasi) SetLink(data string) {
	n.link = data
}

func (n *Notifikasi) SetDibaca(data bool) {
	n.dibaca = data
}

func (n *Notifikasi) SetCreatedAt(data string) {
	n.createdAt = data
}

func (n *Notifikasi) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDNotifikasi int    `json:"idNotifikasi"`
		Penerima     []int  `json:"penerima"`
		Pengirim     string `json:"pengirim"`
		Judul        string `json:"judul"`
		Pesan        string `json:"pesan"`
		Link         string `json:"link"`
		Dibaca       bool   `json:"dibaca"`
		CreatedAt    string `json:"createdAt"`
	}{
		IDNotifikasi: n.idNotifikasi,
		Penerima:     n.penerima,
		Pengirim:     n.pengirim,
		Judul:        n.judul,
		Pesan:        n.pesan,
		Link:         n.link,
		Dibaca:       n.dibaca,
		CreatedAt:    n.createdAt,
	})
}

func (n *Notifikasi) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDNotifikasi int    `json:"idNotifikasi"`
		Penerima     []int  `json:"penerima"`
		Pengirim     string `json:"pengirim"`
		Judul        string `json:"judul"`
		Pesan        string `json:"pesan"`
		Link         string `json:"link"`
		Dibaca       bool   `json:"dibaca"`
		CreatedAt    string `json:"createdAt"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	n.idNotifikasi = alias.IDNotifikasi
	n.penerima = alias.Penerima
	n.pengirim = alias.Pengirim
	n.judul = alias.Judul
	n.pesan = alias.Pesan
	n.link = alias.Link
	n.dibaca = alias.Dibaca
	n.createdAt = alias.CreatedAt

	return nil
}

// GetNotifikasis is func
func (n Notifikasi) GetNotifikasis(idCustomer string) Notifikasis {
	con := db.Connect()
	query := "SELECT idNotifikasi, pengirim, judul, pesan, link, dibaca, createdAt FROM notifikasi WHERE idPenerima = ? ORDER BY idNotifikasi DESC"
	rows, _ := con.Query(query, idCustomer)

	var notifikasis Notifikasis
	var createdAt time.Time

	for rows.Next() {
		rows.Scan(
			&n.idNotifikasi, &n.pengirim, &n.judul, &n.pesan, &n.link, &n.dibaca, &createdAt,
		)

		n.createdAt = createdAt.Format("02 Jan 2006")
		notifikasis.Notifikasis = append(notifikasis.Notifikasis, n)
	}

	defer con.Close()

	return notifikasis
}

// CreateNotifikasi is func
func (n Notifikasi) CreateNotifikasi() error {
	con := db.Connect()
	query := "INSERT INTO notifikasi (idPenerima, pengirim, judul, pesan, link, createdAt) VALUES (?,?,?,?,?,?)"
	var err error
	for _, vIDPenerima := range n.penerima {
		_, err = con.Exec(query, vIDPenerima, n.pengirim, n.judul, n.pesan, n.link, n.createdAt)
	}

	defer con.Close()
	return err
}

// ReadNotifikasi is func
func (n Notifikasi) ReadNotifikasi(idPenerima string) error {
	con := db.Connect()
	query := "UPDATE notifikasi SET dibaca = 1 WHERE idPenerima = ?"
	_, err := con.Exec(query, idPenerima)

	defer con.Close()

	return err
}
