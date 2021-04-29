package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Penangan is class
type Penangan struct {
	idPenangan   int
	idKaryawan   int
	namaKaryawan string
}

func (p *Penangan) GetIDPenangan() int {
	return p.idPenangan
}

func (p *Penangan) GetIDKaryawan() int {
	return p.idKaryawan
}

func (p *Penangan) GetNamaKaryawan() string {
	return p.namaKaryawan
}

func (p *Penangan) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDPenangan   int    `json:"idPenangan"`
		IDKaryawan   int    `json:"idKaryawan"`
		NamaKaryawan string `json:"namaKaryawan"`
	}{
		IDPenangan:   p.idPenangan,
		IDKaryawan:   p.idKaryawan,
		NamaKaryawan: p.namaKaryawan,
	})
}

func (p *Penangan) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDPenangan   int    `json:"idPenangan"`
		IDKaryawan   int    `json:"idKaryawan"`
		NamaKaryawan string `json:"namaKaryawan"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.idPenangan = alias.IDPenangan
	p.idKaryawan = alias.IDKaryawan
	p.namaKaryawan = alias.NamaKaryawan

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetPenangan is func
func (p Penangan) GetPenangan(idOrder string) (Penangan, error) {
	con := db.Connect()
	query := "SELECT a.idPenangan, a.idKaryawan, b.namaKaryawan FROM penangan a JOIN karyawan b ON a.idKaryawan = b.idKaryawan WHERE a.idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&p.idPenangan, &p.idKaryawan, &p.namaKaryawan)

	defer con.Close()

	return p, err
}

// SetPenangan is func
func (p Penangan) SetPenangan(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO penangan (idPenangan, idOrder, idKaryawan) VALUES (?,?,?) ON DUPLICATE KEY UPDATE idKaryawan = ?"
	_, err := con.Exec(query, p.idPenangan, idOrder, p.idKaryawan, p.idKaryawan)

	defer con.Close()

	return err
}
