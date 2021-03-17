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

// Undangans is list of undangan
type Undangans struct {
	Undangans []Undangan `json:"undangan"`
}

// GetUndangans is func
func (u Undangan) GetUndangans(idToko string) Undangans {
	con := db.Connect()
	query := "SELECT a.idUndangan, a.idToko, a.idCustomer, a.posisi, a.status, b.namaToko, c.nama, c.email, a.date FROM undangan a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idToko = ? ORDER BY a.idUndangan DESC"
	rows, _ := con.Query(query, idToko)

	var tgl time.Time
	var undangans Undangans

	for rows.Next() {
		rows.Scan(
			&u.IDUndangan, &u.IDToko, &u.IDCustomer, &u.Posisi, &u.Status, &u.NamaToko, &u.NamaCustomer, &u.Email, &tgl,
		)

		u.Date = tgl.Format("02 Jan 2006")
		undangans.Undangans = append(undangans.Undangans, u)
	}

	defer con.Close()

	return undangans
}

// GetUndangan is func
func (u Undangan) GetUndangan(idUndangan string) (Undangan, error) {
	con := db.Connect()
	query := "SELECT a.idUndangan, a.idToko, a.idCustomer, a.posisi, a.status, b.namaToko, c.nama, c.email, a.date FROM undangan a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idUndangan = ?"

	var tgl time.Time

	err := con.QueryRow(query, idUndangan).Scan(
		&u.IDUndangan, &u.IDToko, &u.IDCustomer, &u.Posisi, &u.Status, &u.NamaToko, &u.NamaCustomer, &u.Email, &tgl)

	u.Date = tgl.Format("02 Jan 2006")

	defer con.Close()

	return u, err
}

// UndangKaryawan is func
func (u Undangan) UndangKaryawan(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO undangan (idUndangan, idToko, idCustomer, posisi, status, date) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE posisi = ?, status = ?, date = ?"
	exec, err := con.Exec(query, u.IDUndangan, idToko, u.IDCustomer, u.Posisi, u.Status, u.Date, u.Posisi, u.Status, u.Date)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idUndangan := int(idInt64)

	defer con.Close()

	return idUndangan, err
}

// TolakUndangan is func
func (u Undangan) TolakUndangan(idUndangan string) error {
	con := db.Connect()
	query := "UPDATE undangan SET status = 'ditolak' WHERE idUndangan = ?"
	_, err := con.Exec(query, idUndangan)

	defer con.Close()

	return err
}

// TerimaUndangan is func
func (u Undangan) TerimaUndangan(idUndangan string) error {
	con := db.Connect()
	query := "UPDATE undangan SET status = 'diterima' WHERE idUndangan = ?"
	_, err := con.Exec(query, idUndangan)

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
