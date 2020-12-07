package models

import "zonart/db"

// GrupOpsi is class
type GrupOpsi struct {
	IDGrupOpsi      int    `json:"idGrupOpsi"`
	IDToko          int    `json:"idToko"`
	NamaGrup        string `json:"namaGrup"`
	Required        bool   `json:"required"`
	Min             int    `json:"min"`
	Max             int    `json:"max"`
	SpesificRequest bool   `json:"spesificRequest"`
	Opsi            []Opsi `json:"opsi"`
}

// GrupOpsis is list of grupOpsi
type GrupOpsis struct {
	GrupOpsis []GrupOpsi `json:"grupOpsi"`
}

// GetGrupOpsis is func
func (grupOpsi GrupOpsi) GetGrupOpsis(idToko string) GrupOpsis {
	con := db.Connect()
	query := "SELECT idGrupOpsi, idToko, namaGrup, required, min, max, spesificRequest FROM grupOpsi WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var grupOpsis GrupOpsis

	for rows.Next() {
		rows.Scan(
			&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest,
		)

		grupOpsis.GrupOpsis = append(grupOpsis.GrupOpsis, grupOpsi)
	}

	defer con.Close()

	return grupOpsis
}

// GetGrupOpsi is func
func (grupOpsi GrupOpsi) GetGrupOpsi(idToko, idGrupOpsi string) (GrupOpsi, error) {
	con := db.Connect()
	query := "SELECT idGrupOpsi, idToko, namaGrup, required, min, max, spesificRequest FROM grupOpsi WHERE idToko = ? AND idGrupOpsi = ?"

	err := con.QueryRow(query, idToko, idGrupOpsi).Scan(
		&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest)

	defer con.Close()

	return grupOpsi, err
}

// GetGrupOpsiProduk is func
func (grupOpsi GrupOpsi) GetGrupOpsiProduk(idToko, idProduk string) []GrupOpsi {
	con := db.Connect()
	query := "SELECT a.idGrupOpsi, a.idToko, a.namaGrup, a.required, a.min, a.max, a.spesificRequest FROM grupOpsi a JOIN grupOpsiProduk b ON a.idGrupOpsi = b.idGrupOpsi WHERE b.idProduk = ? AND a.idToko = ?"
	rows, _ := con.Query(query, idProduk, idToko)

	var gops []GrupOpsi

	for rows.Next() {
		rows.Scan(
			&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest,
		)

		gops = append(gops, grupOpsi)
	}

	defer con.Close()

	return gops
}
