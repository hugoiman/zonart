package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Penggajian is class
type Penggajian struct {
	idPenggajian int
	namaKaryawan string
	nominal      int
	tglTransaksi string
}

// Penggajians is list of penggajian
type Penggajians struct {
	Penggajians []Penggajian `json:"penggajian"`
}

func (p *Penggajian) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDPenggajian int    `json:"idPenggajian"`
		NamaKaryawan string `json:"namaKaryawan"`
		Nominal      int    `json:"nominal"`
		TglTransaksi string `json:"tglTransaksi"`
	}{
		IDPenggajian: p.idPenggajian,
		NamaKaryawan: p.namaKaryawan,
		Nominal:      p.nominal,
		TglTransaksi: p.tglTransaksi,
	})
}

func (p *Penggajian) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDPenggajian int    `json:"idPenggajian"`
		NamaKaryawan string `json:"namaKaryawan"`
		Nominal      int    `json:"nominal" validate:"required"`
		TglTransaksi string `json:"tglTransaksi" validate:"required"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.idPenggajian = alias.IDPenggajian
	p.namaKaryawan = alias.NamaKaryawan
	p.nominal = alias.Nominal
	p.tglTransaksi = alias.TglTransaksi

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
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
			&p.idPenggajian, &p.namaKaryawan, &p.nominal, &tglTransaksi,
		)

		p.tglTransaksi = tglTransaksi.Format("02 Jan 2006")
		penggajians.Penggajians = append(penggajians.Penggajians, p)
	}

	defer con.Close()

	return penggajians
}

// CreateGaji is func
func (p Penggajian) CreateGaji(idToko, idKaryawan string) error {
	con := db.Connect()
	query := "INSERT INTO penggajian (idToko, idKaryawan, nominal, tglTransaksi) VALUES (?,?,?,?)"
	_, err := con.Exec(query, idToko, idKaryawan, p.nominal, p.tglTransaksi)

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
