package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Karyawan is class
type Karyawan struct {
	idKaryawan   int
	namaKaryawan string
	email        string
	hp           string
	posisi       string
	status       string
	alamat       string
	bergabung    string
}

func (k *Karyawan) GetIDKaryawan() int {
	return k.idKaryawan
}

func (k *Karyawan) SetNamaKaryawan(data string) {
	k.namaKaryawan = data
}

func (k *Karyawan) GetNamaKaryawan() string {
	return k.namaKaryawan
}

func (k *Karyawan) SetEmail(data string) {
	k.email = data
}

func (k *Karyawan) SetHP(data string) {
	k.hp = data
}

func (k *Karyawan) SetPosisi(data string) {
	k.posisi = data
}

func (k *Karyawan) GetPosisi() string {
	return k.posisi
}

func (k *Karyawan) SetStatus(data string) {
	k.status = data
}

func (k *Karyawan) GetStatus() string {
	return k.status
}

func (k *Karyawan) SetAlamat(data string) {
	k.alamat = data
}

func (k *Karyawan) SetBergabung(data string) {
	k.bergabung = data
}

func (k *Karyawan) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDKaryawan   int    `json:"idKaryawan"`
		NamaKaryawan string `json:"namaKaryawan"`
		Email        string `json:"email"`
		Hp           string `json:"hp"`
		Posisi       string `json:"posisi"`
		Status       string `json:"status"`
		Alamat       string `json:"alamat"`
		Bergabung    string `json:"bergabung"`
	}{
		IDKaryawan:   k.idKaryawan,
		NamaKaryawan: k.namaKaryawan,
		Email:        k.email,
		Hp:           k.hp,
		Posisi:       k.posisi,
		Status:       k.status,
		Alamat:       k.alamat,
		Bergabung:    k.bergabung,
	})
}

func (k *Karyawan) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDKaryawan   int    `json:"idKaryawan"`
		NamaKaryawan string `json:"namaKaryawan"`
		Email        string `json:"email" validate:"required,email"`
		Hp           string `json:"hp"`
		Posisi       string `json:"posisi"`
		Status       string `json:"status"`
		Alamat       string `json:"alamat"`
		Bergabung    string `json:"bergabung"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	k.idKaryawan = alias.IDKaryawan
	k.namaKaryawan = alias.NamaKaryawan
	k.email = alias.Email
	k.hp = alias.Hp
	k.posisi = alias.Posisi
	k.status = alias.Status
	k.alamat = alias.Alamat
	k.bergabung = alias.Bergabung

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetKaryawan is func
func (k Karyawan) GetKaryawan(idToko, idKaryawan string) (Karyawan, error) {
	con := db.Connect()
	query := "SELECT idKaryawan, namaKaryawan, email, hp, posisi, status, alamat, bergabung FROM karyawan WHERE idToko = ? AND idKaryawan = ?"

	var bergabung time.Time

	err := con.QueryRow(query, idToko, idKaryawan).Scan(
		&k.idKaryawan, &k.namaKaryawan, &k.email, &k.hp, &k.posisi, &k.status, &k.alamat, &bergabung)

	k.bergabung = bergabung.Format("02 Jan 2006")

	defer con.Close()

	return k, err
}

// GetKaryawanByIDCustomer is func
func (k Karyawan) GetKaryawanByIDCustomer(idToko, idCustomer string) (Karyawan, error) {
	con := db.Connect()
	query := "SELECT idKaryawan, namaKaryawan, email, hp, posisi, status, alamat, bergabung FROM karyawan WHERE idToko = ? AND idCustomer = ?"

	var bergabung time.Time

	err := con.QueryRow(query, idToko, idCustomer).Scan(
		&k.idKaryawan, &k.namaKaryawan, &k.email, &k.hp, &k.posisi, &k.status, &k.alamat, &bergabung)

	k.bergabung = bergabung.Format("02 Jan 2006")

	defer con.Close()

	return k, err
}

// GetKaryawans is func
func (k Karyawan) GetKaryawans(idToko string) []Karyawan {
	con := db.Connect()
	query := "SELECT idKaryawan, namaKaryawan, email, hp, posisi, status, alamat, bergabung FROM karyawan WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var bergabung time.Time
	var karyawans []Karyawan

	for rows.Next() {
		rows.Scan(
			&k.idKaryawan, &k.namaKaryawan, &k.email, &k.hp, &k.posisi, &k.status, &k.alamat, &bergabung,
		)

		k.bergabung = bergabung.Format("02 Jan 2006")
		karyawans = append(karyawans, k)
	}

	defer con.Close()

	return karyawans
}

func (k Karyawan) GetIDCustomerByKaryawan(idKaryawan string) (string, error) {
	var idCustomer string

	con := db.Connect()
	query := "SELECT idCustomer FROM karyawan WHERE idKaryawan = ?"
	err := con.QueryRow(query, idKaryawan).Scan(&idCustomer)

	defer con.Close()

	return idCustomer, err
}

// UpdateKaryawan is func
func (k Karyawan) UpdateKaryawan(idToko, idKaryawan string) error {
	con := db.Connect()
	query := "UPDATE karyawan SET namaKaryawan = ?, email = ?, hp = ?, posisi = ?, status = ?, alamat = ? WHERE idToko = ? AND idKaryawan = ?"
	_, err := con.Exec(query, k.namaKaryawan, k.email, k.hp, k.posisi, k.status, k.alamat, idToko, idKaryawan)

	defer con.Close()

	return err
}

// CreateKaryawan is func
func (k Karyawan) CreateKaryawan(idToko, idCustomer string) error {
	con := db.Connect()
	query := "INSERT INTO karyawan (idToko, idCustomer, namaKaryawan, email, hp, posisi, status, alamat, bergabung) VALUES (?,?,?,?,?,?,?,?,?)"
	_, err := con.Exec(query, idToko, idCustomer, k.namaKaryawan, k.email, k.hp, k.posisi, k.status, k.alamat, k.bergabung)

	defer con.Close()

	return err
}

// GetAdmins is func
func (k Karyawan) GetAdmins(idToko string) []int {
	con := db.Connect()
	query := "SELECT idCustomer FROM karyawan WHERE idToko = ? AND posisi = 'admin' AND status = 'aktif'"
	rows, _ := con.Query(query, idToko)

	var idCustomer int
	admin := []int{}
	for rows.Next() {
		rows.Scan(
			&idCustomer,
		)

		admin = append(admin, idCustomer)
	}

	defer con.Close()

	return admin
}
