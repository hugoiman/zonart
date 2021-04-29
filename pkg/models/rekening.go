package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Rekening is class
type Rekening struct {
	idRekening int
	bank       string
	norek      string
	pemilik    string
}

func (r *Rekening) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDRekening int    `json:"idRekening"`
		Bank       string `json:"bank"`
		Norek      string `json:"norek"`
		Pemilik    string `json:"pemilik"`
	}{
		IDRekening: r.idRekening,
		Bank:       r.bank,
		Norek:      r.norek,
		Pemilik:    r.pemilik,
	})
}

func (r *Rekening) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDRekening int    `json:"idRekening"`
		Bank       string `json:"bank" validate:"max=50"`
		Norek      string `json:"norek" validate:"max=20"`
		Pemilik    string `json:"pemilik" validate:"max=50"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	r.idRekening = alias.IDRekening
	r.bank = alias.Bank
	r.norek = alias.Norek
	r.pemilik = alias.Pemilik

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetRekening is func
func (r Rekening) GetRekening(id string) []Rekening {
	con := db.Connect()
	var rekenings []Rekening

	query := "SELECT a.idRekening, a.bank, a.norek, a.pemilik FROM rekening a JOIN toko b ON a.idToko = b.idToko WHERE b.idToko = ? OR b.slug = ?"

	rows, _ := con.Query(query, id, id)
	for rows.Next() {
		rows.Scan(
			&r.idRekening, &r.bank, &r.norek, &r.pemilik,
		)

		rekenings = append(rekenings, r)
	}

	defer con.Close()

	return rekenings
}

// CreateUpdateRekening is func
func (r Rekening) CreateUpdateRekening(idToko string) error {
	con := db.Connect()
	query := "INSERT INTO rekening (idRekening, idToko, bank, norek, pemilik) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE bank = ?, norek = ?, pemilik = ?"
	_, err := con.Exec(query, r.idRekening, idToko, r.bank, r.norek, r.pemilik, r.bank, r.norek, r.pemilik)

	defer con.Close()

	return err
}

// DeleteRekening is func
func (r Rekening) DeleteRekening(idToko, idRekening string) error {
	con := db.Connect()
	query := "DELETE FROM rekening WHERE idToko = ? AND idRekening = ?"
	_, err := con.Exec(query, idToko, idRekening)

	defer con.Close()

	return err
}
