package models

import "zonart/db"

// Opsi is class
type Opsi struct {
	IDOpsi     int    `json:"idOpsi"`
	IDGrupOpsi int    `json:"idGrupOpsi"`
	NamaGrup   string `json:"namaGrup"`
	Opsi       string `json:"opsi" validate:"required,max=50"`
	Harga      int    `json:"harga" validate:"max=10"`
	Berat      int    `json:"berat" validate:"max=10"`
	PerProduk  bool   `json:"perProduk"`
	Status     bool   `json:"status"`
}

// GetOpsi is func
func (opsi Opsi) GetOpsi(idGrupOpsi, idOpsi string) (Opsi, error) {
	con := db.Connect()
	query := "SELECT idOpsi, idGrupOpsi, opsi, harga, berat, perProduk, status FROM opsi WHERE idGrupOpsi = ? AND idOpsi = ?"

	err := con.QueryRow(query, idGrupOpsi, idOpsi).Scan(
		&opsi.IDOpsi, &opsi.IDGrupOpsi, &opsi.Opsi, &opsi.Harga, &opsi.Berat, &opsi.PerProduk, &opsi.Status)

	defer con.Close()
	return opsi, err
}

// GetOpsis is func
func (opsi Opsi) GetOpsis(idGrupOpsi string) []Opsi {
	con := db.Connect()
	query := "SELECT idOpsi, idGrupOpsi, opsi, harga, berat, perProduk, status FROM opsi WHERE idGrupOpsi = ?"
	rows, _ := con.Query(query, idGrupOpsi)

	var opsis []Opsi

	for rows.Next() {
		rows.Scan(
			&opsi.IDOpsi, &opsi.IDGrupOpsi, &opsi.Opsi, &opsi.Harga, &opsi.Berat, &opsi.PerProduk, &opsi.Status,
		)

		opsis = append(opsis, opsi)
	}

	defer con.Close()

	return opsis
}

// CreateUpdateOpsi is func
func (opsi Opsi) CreateUpdateOpsi(idGrupOpsi string) error {
	con := db.Connect()
	query := "INSERT INTO opsi (idOpsi, idGrupOpsi, opsi, harga, berat, perProduk, status) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE opsi = ?, harga = ?, berat = ?, perProduk = ?, status = ?"
	_, err := con.Exec(query, opsi.IDOpsi, idGrupOpsi, opsi.Opsi, opsi.Harga, opsi.Berat, opsi.PerProduk, opsi.Status, opsi.Opsi, opsi.Harga, opsi.Berat, opsi.PerProduk, opsi.Status)

	defer con.Close()

	return err
}

// DeleteOpsi is func
func (opsi Opsi) DeleteOpsi(idGrupOpsi, idOpsi string) error {
	con := db.Connect()
	query := "DELETE FROM opsi WHERE idGrupOpsi = ? AND idOpsi = ?"
	_, err := con.Exec(query, idGrupOpsi, idOpsi)

	defer con.Close()

	return err
}
