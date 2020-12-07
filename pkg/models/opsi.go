package models

import "zonart/db"

// Opsi is class
type Opsi struct {
	IDOpsi     int    `json:"idOpsi"`
	IDGrupOpsi int    `json:"idGrupOpsi"`
	Opsi       string `json:"opsi"`
	Harga      int    `json:"harga"`
	Berat      int    `json:"berat"`
}

// GetOpsis is func
func (opsi Opsi) GetOpsis(idGrupOpsi string) []Opsi {
	con := db.Connect()
	query := "SELECT idOpsi, idGrupOpsi, opsi, harga, berat FROM opsi WHERE idGrupOpsi = ?"
	rows, _ := con.Query(query, idGrupOpsi)

	var opsis []Opsi

	for rows.Next() {
		rows.Scan(
			&opsi.IDOpsi, &opsi.IDGrupOpsi, &opsi.Opsi, &opsi.Harga, &opsi.Berat,
		)

		opsis = append(opsis, opsi)
	}

	defer con.Close()

	return opsis
}
