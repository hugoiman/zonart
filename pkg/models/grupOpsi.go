package models

import (
	"strconv"
	"zonart/db"
)

// GrupOpsi is class
type GrupOpsi struct {
	IDGrupOpsi      int    `json:"idGrupOpsi"`
	IDToko          int    `json:"idToko"`
	NamaGrup        string `json:"namaGrup" validate:"required"`
	Required        bool   `json:"required"`
	Min             int    `json:"min"`
	Max             int    `json:"max"`
	SpesificRequest bool   `json:"spesificRequest"`
	Opsi            []Opsi `json:"opsi" validate:"dive"`
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
	var opsi Opsi

	for rows.Next() {
		rows.Scan(
			&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest,
		)

		grupOpsi.Opsi = opsi.GetOpsis(strconv.Itoa(grupOpsi.IDGrupOpsi))
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

	var opsi Opsi
	grupOpsi.Opsi = opsi.GetOpsis(idGrupOpsi)

	defer con.Close()

	return grupOpsi, err
}

// GetGrupOpsiProduk is get all grup opsi in a produk
func (grupOpsi GrupOpsi) GetGrupOpsiProduk(idToko, idProduk string) []GrupOpsi {
	con := db.Connect()
	query := "SELECT a.idGrupOpsi, a.idToko, a.namaGrup, a.required, a.min, a.max, a.spesificRequest FROM grupOpsi a JOIN grupOpsiProduk b ON a.idGrupOpsi = b.idGrupOpsi WHERE b.idProduk = ? AND a.idToko = ?"
	rows, _ := con.Query(query, idProduk, idToko)

	var gop []GrupOpsi
	var opsi Opsi

	for rows.Next() {
		rows.Scan(
			&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest,
		)

		grupOpsi.Opsi = opsi.GetOpsis(strconv.Itoa(grupOpsi.IDGrupOpsi))
		gop = append(gop, grupOpsi)
	}

	defer con.Close()

	return gop
}

// CreateGrupOpsi is func
func (grupOpsi GrupOpsi) CreateGrupOpsi(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO grupOpsi (idToko, namaGrup, required, min, max, spesificRequest) VALUES (?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, grupOpsi.NamaGrup, grupOpsi.Required, grupOpsi.Min, grupOpsi.Max, grupOpsi.SpesificRequest)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idGrupOpsi := int(idInt64)

	defer con.Close()

	return idGrupOpsi, err
}

// UpdateGrupOpsi is func
func (grupOpsi GrupOpsi) UpdateGrupOpsi(idToko, idGrupOpsi string) error {
	con := db.Connect()
	query := "UPDATE grupOpsi SET namaGrup = ?, required = ?, min = ?, max = ?, spesificRequest = ? WHERE idToko = ? AND idGrupOpsi = ?"
	_, err := con.Exec(query, grupOpsi.NamaGrup, grupOpsi.Required, grupOpsi.Min, grupOpsi.Max, grupOpsi.SpesificRequest, idToko, idGrupOpsi)

	defer con.Close()

	return err
}

// DeleteGrupOpsi is func
func (grupOpsi GrupOpsi) DeleteGrupOpsi(idToko, idGrupOpsi string) error {
	con := db.Connect()
	query := "DELETE FROM grupOpsi WHERE idToko = ? AND idGrupOpsi = ?"
	_, err := con.Exec(query, idToko, idGrupOpsi)

	defer con.Close()

	return err
}
