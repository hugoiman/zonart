package models

import (
	"time"
	"zonart/db"
)

// Penggajian is class
type Penggajian struct {
	IDPenggajian int    `json:"idPenggajian"`
	IDToko       int    `json:"idToko"`
	IDKaryawan   int    `json:"idKaryawan"`
	NamaKaryawan string `json:"namaKaryawan" validate:"required"`
	Periode      string `json:"periode" validate:"required"`
	Nominal      string `json:"nominal" validate:"required"`
	TglTransaksi string `json:"tglTransaksi" validate:"required"`
}

// Penggajians is list of penggajian
type Penggajians struct {
	Penggajians []Penggajian `json:"penggajian"`
}

// GetPenggajians is func
func (p Penggajian) GetPenggajians(idToko string) Penggajians {
	con := db.Connect()
	query := "SELECT idPenggajian, idToko, idKaryawan, periode, nominal, tglTransaksi FROM penggajian WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var tglTransaksi time.Time
	var penggajians Penggajians

	for rows.Next() {
		rows.Scan(
			&p.IDPenggajian, &p.IDToko, &p.IDKaryawan, &p.Periode, &p.Nominal, &tglTransaksi,
		)

		p.TglTransaksi = tglTransaksi.Format("02 Jan 2006")
		penggajians.Penggajians = append(penggajians.Penggajians, p)
	}

	defer con.Close()

	return penggajians
}

// CreatePenggajian is func
func (p Penggajian) CreatePenggajian(idToko string) error {
	con := db.Connect()
	query := "INSERT INTO penggajian (idToko, idKaryawan, periode, nominal, tglTransaksi) VALUES (?,?,?,?,?)"
	_, err := con.Exec(query, idToko, p.IDToko, p.IDKaryawan, p.Periode, p.Nominal, p.TglTransaksi)

	defer con.Close()

	return err
}

// DeletePenggajian is func
func (p Penggajian) DeletePenggajian(idToko, idPenggajian string) error {
	con := db.Connect()
	query := "DELETE FROM penggajian WHERE idToko = ? AND idPenggajian = ?"
	_, err := con.Exec(query, idToko, idPenggajian)

	defer con.Close()

	return err
}
