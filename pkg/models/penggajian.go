package models

import (
	"time"
	"zonart/db"
)

// Penggajian is class
type Penggajian struct {
	IDPenggajian int    `json:"idPenggajian"`
	NamaKaryawan string `json:"namaKaryawan"`
	Nominal      int    `json:"nominal" validate:"required"`
	TglTransaksi string `json:"tglTransaksi" validate:"required"`
}

// Penggajians is list of penggajian
type Penggajians struct {
	Penggajians []Penggajian `json:"penggajian"`
}

// GetGajis is func
func (p Penggajian) GetGajis(idToko string) Penggajians {
	con := db.Connect()
	query := "SELECT a.idPenggajian, b.namaKaryawan, a.nominal, a.tglTransaksi FROM penggajian a " +
		"JOIN karyawan b ON a.idKaryawan = b.idKaryawan WHERE a.idToko = ?"
	rows, _ := con.Query(query, idToko)

	var tglTransaksi time.Time
	var penggajians Penggajians

	for rows.Next() {
		rows.Scan(
			&p.IDPenggajian, &p.NamaKaryawan, &p.Nominal, &tglTransaksi,
		)

		p.TglTransaksi = tglTransaksi.Format("02 Jan 2006")
		penggajians.Penggajians = append(penggajians.Penggajians, p)
	}

	defer con.Close()

	return penggajians
}

// CreateGaji is func
func (p Penggajian) CreateGaji(idToko, idKaryawan string) error {
	con := db.Connect()
	query := "INSERT INTO penggajian (idToko, idKaryawan, nominal, tglTransaksi) VALUES (?,?,?,?)"
	_, err := con.Exec(query, idToko, idKaryawan, p.Nominal, p.TglTransaksi)

	defer con.Close()

	return err
}

// DeleteGaji is func
func (p Penggajian) DeleteGaji(idToko, idPenggajian string) error {
	con := db.Connect()
	query := "DELETE FROM penggajian WHERE idToko = ? AND idPenggajian = ?"
	_, err := con.Exec(query, idToko, idPenggajian)

	defer con.Close()

	return err
}
