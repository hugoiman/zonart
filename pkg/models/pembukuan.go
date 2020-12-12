package models

import (
	"time"
	"zonart/db"
)

// Pembukuan is class
type Pembukuan struct {
	IDPembukuan  int    `json:"idPembukuan"`
	IDToko       int    `json:"idToko"`
	Jenis        string `json:"jenis"`
	Keterangan   string `json:"keterangan" validate:"required"`
	Nominal      int    `json:"nominal" validate:"required"`
	TglTransaksi string `json:"tglTransaksi" validate:"required"`
}

// Pembukuans is list of pembukuan
type Pembukuans struct {
	Pembukuans []Pembukuan `json:"pembukuan"`
}

// GetPembukuans is func
func (p Pembukuan) GetPembukuans(idToko string) Pembukuans {
	con := db.Connect()
	query := "SELECT idPembukuan, idToko, jenis, keterangan, nominal, tglTransaksi FROM pembukuan WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var tglTransaksi time.Time
	var pembukuans Pembukuans

	for rows.Next() {
		rows.Scan(
			&p.IDPembukuan, &p.IDToko, &p.Jenis, &p.Keterangan, &p.Nominal, &tglTransaksi,
		)
		p.TglTransaksi = tglTransaksi.Format("02 Jan 2006")
		pembukuans.Pembukuans = append(pembukuans.Pembukuans, p)
	}

	defer con.Close()

	return pembukuans
}

// CreatePembukuan is func
func (p Pembukuan) CreatePembukuan(idToko string) error {
	con := db.Connect()
	query := "INSERT INTO pembukuan (idToko, jenis, keterangan, nominal, tglTransaksi) VALUES (?,?,?,?,?)"
	_, err := con.Exec(query, idToko, p.Jenis, p.Keterangan, p.Nominal, p.TglTransaksi)

	defer con.Close()

	return err
}

// DeletePengeluaran is func
func (p Pembukuan) DeletePengeluaran(idToko, idPembukuan string) error {
	con := db.Connect()
	query := "DELETE FROM pembukuan WHERE idToko = ? AND idPembukuan = ? AND jenis = 'pengeluaran'"
	_, err := con.Exec(query, idToko, idPembukuan)

	defer con.Close()

	return err
}
