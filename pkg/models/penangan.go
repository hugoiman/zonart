package models

import "zonart/db"

// Penangan is class
type Penangan struct {
	IDPenangan   int    `json:"idPenangan"`
	IDOrder      int    `json:"idOrder"`
	IDKaryawan   int    `json:"idKaryawan"`
	NamaKaryawan string `json:"namaKaryawan"`
}

// GetPenangan is func
func (p Penangan) GetPenangan(idOrder string) Penangan {
	con := db.Connect()
	query := "SELECT idPenangan, idOrder, idKaryawan, namaKaryawan FROM penangan WHERE idOrder = ?"

	_ = con.QueryRow(query, idOrder).Scan(
		&p.IDPenangan, &p.IDOrder, &p.IDKaryawan, &p.NamaKaryawan)

	defer con.Close()

	return p
}

// SetPenangan is func
func (p Penangan) SetPenangan(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO penangan (idPenangan, idOrder, idKaryawan) VALUES (?,?,?) ON DUPLICATE KEY UPDATE idKaryawan = ?"
	_, err := con.Exec(query, p.IDPenangan, idOrder, p.IDKaryawan, p.IDKaryawan)

	defer con.Close()

	return err
}
