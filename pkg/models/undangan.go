package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Undangan is class
type Undangan struct {
	idUndangan   int
	posisi       string
	status       string
	namaToko     string
	namaCustomer string
	email        string
	date         string
}

func (u *Undangan) SetIDUndangan(data int) {
	u.idUndangan = data
}

func (u *Undangan) GetIDUndangan() int {
	return u.idUndangan
}

func (u *Undangan) GetPosisi() string {
	return u.posisi
}

func (u *Undangan) SetStatus(data string) {
	u.status = data
}

func (u *Undangan) GetStatus() string {
	return u.status
}
func (u *Undangan) GetEmail() string {
	return u.email
}

func (u *Undangan) SetDate(data string) {
	u.date = data
}

func (u *Undangan) GetDate() string {
	return u.date
}

func (u *Undangan) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDUndangan   int    `json:"idUndangan"`
		Posisi       string `json:"posisi"`
		Status       string `json:"status"`
		NamaToko     string `json:"namaToko"`
		NamaCustomer string `json:"namaCustomer"`
		Email        string `json:"email"`
		Date         string `json:"date"`
	}{
		IDUndangan:   u.idUndangan,
		Posisi:       u.posisi,
		Status:       u.status,
		NamaToko:     u.namaToko,
		NamaCustomer: u.namaCustomer,
		Email:        u.email,
		Date:         u.date,
	})
}

func (u *Undangan) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDUndangan   int    `json:"idUndangan"`
		Posisi       string `json:"posisi" validate:"required"`
		Status       string `json:"status"`
		NamaToko     string `json:"namaToko"`
		NamaCustomer string `json:"namaCustomer"`
		Email        string `json:"email" validate:"required"`
		Date         string `json:"date"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	u.idUndangan = alias.IDUndangan
	u.posisi = alias.Posisi
	u.status = alias.Status
	u.namaToko = alias.NamaToko
	u.namaCustomer = alias.NamaCustomer
	u.email = alias.Email
	u.date = alias.Date

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetUndangans is func
func (u Undangan) GetUndangans(idToko string) []Undangan {
	con := db.Connect()
	query := "SELECT a.idUndangan, a.posisi, a.status, b.namaToko, c.nama, c.email, a.date FROM undangan a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idToko = ? ORDER BY a.idUndangan DESC"
	rows, _ := con.Query(query, idToko)

	var tgl time.Time
	var undangans []Undangan

	for rows.Next() {
		rows.Scan(
			&u.idUndangan, &u.posisi, &u.status, &u.namaToko, &u.namaCustomer, &u.email, &tgl,
		)

		u.date = tgl.Format("02 Jan 2006")
		undangans = append(undangans, u)
	}

	defer con.Close()

	return undangans
}

// GetUndangan is func
func (u Undangan) GetUndangan(idUndangan string) (Undangan, error) {
	con := db.Connect()
	query := "SELECT a.idUndangan, a.posisi, a.status, b.namaToko, c.nama, c.email, a.date FROM undangan a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idUndangan = ?"

	var tgl time.Time

	err := con.QueryRow(query, idUndangan).Scan(
		&u.idUndangan, &u.posisi, &u.status, &u.namaToko, &u.namaCustomer, &u.email, &tgl)

	u.date = tgl.Format("02 Jan 2006")

	defer con.Close()

	return u, err
}

// GetUndangan is func
func (u Undangan) GetUndanganCustomer(idUndangan, idCustomer string) (Undangan, error) {
	con := db.Connect()
	query := "SELECT a.idUndangan, a.posisi, a.status, b.namaToko, c.nama, c.email, a.date FROM undangan a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idUndangan = ? AND a.idCustomer = ?"

	var tgl time.Time

	err := con.QueryRow(query, idUndangan, idCustomer).Scan(
		&u.idUndangan, &u.posisi, &u.status, &u.namaToko, &u.namaCustomer, &u.email, &tgl)

	u.date = tgl.Format("02 Jan 2006")

	defer con.Close()

	return u, err
}

// UndangKaryawan is func
func (u Undangan) UndangKaryawan(idToko, idCustomer string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO undangan (idUndangan, idToko, idCustomer, posisi, status, date) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE posisi = ?, status = ?, date = ?"
	exec, err := con.Exec(query, u.idUndangan, idToko, idCustomer, u.posisi, u.status, u.date, u.posisi, u.status, u.date)

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
	query := "SELECT idUndangan, posisi, status FROM undangan WHERE idToko = ? AND idCustomer = ?"

	err := con.QueryRow(query, idToko, idCustomer).Scan(
		&u.idUndangan, &u.posisi, &u.status)

	defer con.Close()

	return u, err
}
