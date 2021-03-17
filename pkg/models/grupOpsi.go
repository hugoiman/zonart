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
	Deskripsi       string `json:"deskripsi"`
	Required        bool   `json:"required"`
	Min             int    `json:"min"`
	Max             int    `json:"max"`
	SpesificRequest bool   `json:"spesificRequest"`
	HardCopy        bool   `json:"hardcopy"`
	SoftCopy        bool   `json:"softcopy"`
	Opsi            []Opsi `json:"opsi" validate:"dive"`
}

// GrupOpsis is list of grupOpsi
type GrupOpsis struct {
	GrupOpsis []GrupOpsi `json:"grupopsi"`
}

// GetGrupOpsis is func
func (grupOpsi GrupOpsi) GetGrupOpsis(idToko string) GrupOpsis {
	con := db.Connect()
	query := "SELECT idGrupOpsi, idToko, namaGrup, deskripsi, required, min, max, spesificRequest, hardcopy, softcopy FROM grupOpsi WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var grupOpsis GrupOpsis
	var opsi Opsi

	for rows.Next() {
		rows.Scan(
			&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Deskripsi, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest, &grupOpsi.HardCopy, &grupOpsi.SoftCopy,
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
	query := "SELECT idGrupOpsi, idToko, namaGrup, deskripsi, required, min, max, spesificRequest, hardcopy, softcopy FROM grupOpsi WHERE idToko = ? AND idGrupOpsi = ?"

	err := con.QueryRow(query, idToko, idGrupOpsi).Scan(
		&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Deskripsi, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest, &grupOpsi.HardCopy, &grupOpsi.SoftCopy)

	var opsi Opsi
	grupOpsi.Opsi = opsi.GetOpsis(idGrupOpsi)

	defer con.Close()

	return grupOpsi, err
}

// GetGrupOpsiProduk is get all grup opsi in a produk
func (grupOpsi GrupOpsi) GetGrupOpsiProduk(idToko, idProduk string) []GrupOpsi {
	con := db.Connect()
	query := "SELECT a.idGrupOpsi, a.idToko, a.namaGrup, a.deskripsi, a.required, a.min, a.max, a.spesificRequest, hardcopy, softcopy FROM grupOpsi a JOIN grupOpsiProduk b ON a.idGrupOpsi = b.idGrupOpsi WHERE b.idProduk = ? AND a.idToko = ?"
	rows, _ := con.Query(query, idProduk, idToko)

	var gop []GrupOpsi
	var opsi Opsi

	for rows.Next() {
		rows.Scan(
			&grupOpsi.IDGrupOpsi, &grupOpsi.IDToko, &grupOpsi.NamaGrup, &grupOpsi.Deskripsi, &grupOpsi.Required, &grupOpsi.Min, &grupOpsi.Max, &grupOpsi.SpesificRequest, &grupOpsi.HardCopy, &grupOpsi.SoftCopy,
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
	query := "INSERT INTO grupOpsi (idToko, namaGrup, deskripsi, required, min, max, spesificRequest, hardcopy, softcopy) VALUES (?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, grupOpsi.NamaGrup, grupOpsi.Deskripsi, grupOpsi.Required, grupOpsi.Min, grupOpsi.Max, grupOpsi.SpesificRequest, grupOpsi.HardCopy, grupOpsi.SoftCopy)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idGrupOpsi := int(idInt64)

	for _, vOpsi := range grupOpsi.Opsi {
		err = vOpsi.CreateUpdateOpsi(strconv.Itoa(idGrupOpsi))
		if err != nil {
			_ = grupOpsi.DeleteGrupOpsi(idToko, strconv.Itoa(idGrupOpsi))
			return idGrupOpsi, err
		}
	}

	defer con.Close()

	return idGrupOpsi, err
}

// UpdateGrupOpsi is func
func (grupOpsi GrupOpsi) UpdateGrupOpsi(idToko, idGrupOpsi string) error {
	con := db.Connect()
	query := "UPDATE grupOpsi SET namaGrup = ?, deskripsi = ?, required = ?, min = ?, max = ?, spesificRequest = ?, hardcopy = ?, softcopy = ? WHERE idToko = ? AND idGrupOpsi = ?"
	_, err := con.Exec(query, grupOpsi.NamaGrup, grupOpsi.Deskripsi, grupOpsi.Required, grupOpsi.Min, grupOpsi.Max, grupOpsi.SpesificRequest, grupOpsi.HardCopy, grupOpsi.SoftCopy, idToko, idGrupOpsi)

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
