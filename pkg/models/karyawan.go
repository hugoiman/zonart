package models

import (
	"time"
	"zonart/db"
)

// Karyawan is class
type Karyawan struct {
	IDKaryawan   int    `json:"idKaryawan"`
	IDToko       int    `json:"idToko"`
	IDCustomer   int    `json:"idCustomer"`
	NamaKaryawan string `json:"namaKaryawan"`
	Email        string `json:"email" validate:"required"`
	Hp           string `json:"hp"`
	Posisi       string `json:"posisi"`
	Status       string `json:"status"`
	Alamat       string `json:"alamat"`
	Bergabung    string `json:"bergabung"`
}

// Karyawans is list of karyawan
type Karyawans struct {
	Karyawans []Karyawan `json:"karyawan"`
}

// AuthKaryawan is func
func (k Karyawan) AuthKaryawan(idToko, idCustomer string) (Karyawan, error) {
	con := db.Connect()
	query := "SELECT idKaryawan, idToko, idCustomer, namaKaryawan, email, hp, posisi, status, alamat, bergabung FROM karyawan WHERE idToko = ? AND idCustomer = ?"

	var bergabung time.Time

	err := con.QueryRow(query, idToko, idCustomer).Scan(
		&k.IDKaryawan, &k.IDToko, &k.IDCustomer, &k.NamaKaryawan, &k.Email, &k.Hp, &k.Posisi, &k.Status, &k.Alamat, &bergabung)

	k.Bergabung = bergabung.Format("02 Jan 2006")

	defer con.Close()

	return k, err
}

// GetKaryawan is func
func (k Karyawan) GetKaryawan(idToko, idKaryawan string) (Karyawan, error) {
	con := db.Connect()
	query := "SELECT idKaryawan, idToko, idCustomer, namaKaryawan, email, hp, posisi, status, alamat, bergabung FROM karyawan WHERE idToko = ? AND idKaryawan = ?"

	var bergabung time.Time

	err := con.QueryRow(query, idToko, idKaryawan).Scan(
		&k.IDKaryawan, &k.IDToko, &k.IDCustomer, &k.NamaKaryawan, &k.Email, &k.Hp, &k.Posisi, &k.Status, &k.Alamat, &bergabung)

	k.Bergabung = bergabung.Format("02 Jan 2006")

	defer con.Close()

	return k, err
}

// GetKaryawans is func
func (k Karyawan) GetKaryawans(idToko string) Karyawans {
	con := db.Connect()
	query := "SELECT idKaryawan, idToko, idCustomer, namaKaryawan, email, hp, posisi, status, alamat, bergabung FROM karyawan WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var bergabung time.Time
	var karyawans Karyawans

	for rows.Next() {
		rows.Scan(
			&k.IDKaryawan, &k.IDToko, &k.IDCustomer, &k.NamaKaryawan, &k.Email, &k.Hp, &k.Posisi, &k.Status, &k.Alamat, &bergabung,
		)

		k.Bergabung = bergabung.Format("02 Jan 2006")
		karyawans.Karyawans = append(karyawans.Karyawans, k)
	}

	defer con.Close()

	return karyawans
}

// UpdateKaryawan is func
func (k Karyawan) UpdateKaryawan(idToko, idKaryawan string) error {
	con := db.Connect()
	query := "UPDATE karyawan SET namaKaryawan = ?, email = ?, hp = ?, posisi = ?, status = ?, alamat = ? WHERE idToko = ? AND idKaryawan = ?"
	_, err := con.Exec(query, k.NamaKaryawan, k.Email, k.Hp, k.Posisi, k.Status, k.Alamat, idToko, idKaryawan)

	defer con.Close()

	return err
}

// CreateKaryawan is func
func (k Karyawan) CreateKaryawan() error {
	con := db.Connect()
	query := "INSERT INTO karyawan (idToko, idCustomer, namaKaryawan, email, hp, posisi, status, alamat, bergabung) VALUES (?,?,?,?,?,?,?,?,?)"
	_, err := con.Exec(query, k.IDToko, k.IDCustomer, k.NamaKaryawan, k.Email, k.Hp, k.Posisi, k.Status, k.Alamat, k.Bergabung)

	defer con.Close()

	return err
}
