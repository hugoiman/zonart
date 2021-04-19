package models

import "zonart/db"

// Rekening is class
type Rekening struct {
	IDRekening int    `json:"idRekening"`
	Bank       string `json:"bank" validate:"max=50"`
	Norek      string `json:"norek" validate:"max=20"`
	Pemilik    string `json:"pemilik" validate:"max=50"`
}

// GetRekening is func
func (r Rekening) GetRekening(id string) []Rekening {
	con := db.Connect()
	var rekenings []Rekening

	query := "SELECT a.idRekening, a.bank, a.norek, a.pemilik FROM rekening a JOIN toko b ON a.idToko = b.idToko WHERE b.idToko = ? OR b.slug = ?"

	rows, _ := con.Query(query, id, id)
	for rows.Next() {
		rows.Scan(
			&r.IDRekening, &r.Bank, &r.Norek, &r.Pemilik,
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
	_, err := con.Exec(query, r.IDRekening, idToko, r.Bank, r.Norek, r.Pemilik, r.Bank, r.Norek, r.Pemilik)

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
