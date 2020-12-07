package models

import (
	"time"
	"zonart/db"
)

// Undangan is class
type Undangan struct {
	IDUndangan   int    `json:"idUndangan"`
	IDToko       int    `json:"idToko"`
	IDCustomer   int    `json:"idCustomer"`
	Posisi       string `json:"posisi" validate:"required"`
	Status       string `json:"status"`
	NamaToko     string `json:"namaToko"`
	NamaCustomer string `json:"namaCustomer"`
	Email        string `json:"email" validate:"required"`
	Date         string `json:"date"`
}

// GetUndangan is func
func (u Undangan) GetUndangan(idUndangan, idToko, idCustomer string) (Undangan, error) {
	con := db.Connect()
	query := "SELECT a.idUndangan, a.idToko, a.idCustomer, a.posisi, a.status, b.namaToko, c.nama, c.email, a.date FROM undangan a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idToko = ? AND a.idUndangan = ? AND a.idCustomer = ?"

	var tgl time.Time

	err := con.QueryRow(query, idToko, idUndangan, idCustomer).Scan(
		&u.IDUndangan, &u.IDToko, &u.IDCustomer, &u.Posisi, &u.Status, &u.NamaToko, &u.NamaCustomer, &u.Email, &tgl)

	u.Date = tgl.Format("02 Jan 2006")

	defer con.Close()

	return u, err
}

// UndangKaryawan is func
func (u Undangan) UndangKaryawan() error {
	con := db.Connect()
	query := "INSERT INTO undangan (idUndangan, idToko, idCustomer, posisi, status, date) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE posisi = ?, status = ?, date = ?"
	_, err := con.Exec(query, u.IDUndangan, u.IDToko, u.IDCustomer, u.Posisi, u.Status, u.Date, u.Posisi, u.Status, u.Date)

	defer con.Close()

	return err
}

// TolakUndangan is func
func (u Undangan) TolakUndangan(idUndangan, idToko, idCustomer string) error {
	con := db.Connect()
	query := "UPDATE undangan SET status = 'ditolak' WHERE idUndangan = ? AND idToko = ? AND idCustomer = ?"
	_, err := con.Exec(query, idUndangan, idToko, idCustomer)

	defer con.Close()

	return err
}

// TerimaUndangan is func
func (u Undangan) TerimaUndangan(idUndangan, idToko, idCustomer string) error {
	con := db.Connect()
	query := "UPDATE undangan SET status = 'disetujui' WHERE idUndangan = ? AND idToko = ? AND idCustomer = ?"
	_, err := con.Exec(query, idUndangan, idToko, idCustomer)

	defer con.Close()

	return err
}

// BatalkanUndangan is func
func (u Undangan) BatalkanUndangan(idUndangan, idToko string) error {
	con := db.Connect()
	query := "UPDATE undangan SET status = 'dibatalkan' WHERE idUndangan = ? AND idToko = ?"
	_, err := con.Exec(query, idUndangan, idToko)

	defer con.Close()

	return err
}

// CheckUndangan is func
func (u Undangan) CheckUndangan(idToko, idCustomer string) (Undangan, error) {
	con := db.Connect()
	query := "SELECT idUndangan, idToko, idCustomer, posisi, status FROM undangan WHERE idToko = ? AND idCustomer = ?"

	err := con.QueryRow(query, idToko, idCustomer).Scan(
		&u.IDUndangan, &u.IDToko, &u.IDCustomer, &u.Posisi, &u.Status)

	defer con.Close()

	return u, err
}
