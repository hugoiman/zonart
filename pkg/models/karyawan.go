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
	Email        string `json:"email"`
	Hp           string `json:"hp"`
	Posisi       string `json:"posisi"`
	Status       string `json:"status"`
	Alamat       string `json:"alamat"`
	Bergabung    string `json:"bergabung"`
}

// GetKaryawan is func
func (k Karyawan) GetKaryawan(idToko, idCustomer string) (Karyawan, error) {
	con := db.Connect()
	query := "SELECT idKaryawan, idToko, idCustomer, namaKaryawan, email, hp, posisi, status, alamat, bergabung FROM karyawan WHERE idToko = ? AND idCustomer = ?"

	var bergabung time.Time

	err := con.QueryRow(query, idToko, idCustomer).Scan(
		&k.IDKaryawan, &k.IDToko, &k.IDCustomer, &k.NamaKaryawan, &k.Email, &k.Hp, &k.Posisi, &k.Status, &k.Alamat, bergabung)

	k.Bergabung = bergabung.Format("02 Jan 2006")

	defer con.Close()

	return k, err
}
