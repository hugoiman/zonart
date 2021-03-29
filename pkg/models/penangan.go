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
func (p Penangan) GetPenangan(idOrder string) (Penangan, error) {
	con := db.Connect()
	query := "SELECT a.idPenangan, a.idOrder, a.idKaryawan, b.namaKaryawan FROM penangan a JOIN karyawan b ON a.idKaryawan = b.idKaryawan WHERE a.idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&p.IDPenangan, &p.IDOrder, &p.IDKaryawan, &p.NamaKaryawan)

	defer con.Close()

	return p, err
}

// SetPenangan is func
func (p Penangan) SetPenangan(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO penangan (idPenangan, idOrder, idKaryawan) VALUES (?,?,?) ON DUPLICATE KEY UPDATE idKaryawan = ?"
	_, err := con.Exec(query, p.IDPenangan, idOrder, p.IDKaryawan, p.IDKaryawan)

	defer con.Close()

	return err
}
