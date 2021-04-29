package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Opsi is class
type Opsi struct {
	idOpsi    int
	namaGrup  string
	opsi      string
	harga     int
	berat     int
	perProduk bool
	status    bool
}

func (o *Opsi) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDOpsi    int    `json:"idOpsi"`
		NamaGrup  string `json:"namaGrup"`
		Opsi      string `json:"opsi"`
		Harga     int    `json:"harga"`
		Berat     int    `json:"berat"`
		PerProduk bool   `json:"perProduk"`
		Status    bool   `json:"status"`
	}{
		IDOpsi:    o.idOpsi,
		NamaGrup:  o.namaGrup,
		Opsi:      o.opsi,
		Harga:     o.harga,
		Berat:     o.berat,
		PerProduk: o.perProduk,
		Status:    o.status,
	})
}

func (o *Opsi) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDOpsi    int    `json:"idOpsi"`
		NamaGrup  string `json:"namaGrup"`
		Opsi      string `json:"opsi" validate:"required"`
		Harga     int    `json:"harga"`
		Berat     int    `json:"berat"`
		PerProduk bool   `json:"perProduk"`
		Status    bool   `json:"status"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	o.idOpsi = alias.IDOpsi
	o.namaGrup = alias.NamaGrup
	o.opsi = alias.Opsi
	o.harga = alias.Harga
	o.berat = alias.Berat
	o.perProduk = alias.PerProduk
	o.status = alias.Status

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetOpsi is func
func (opsi Opsi) GetOpsi(idGrupOpsi, idOpsi string) (Opsi, error) {
	con := db.Connect()
	query := "SELECT idOpsi, opsi, harga, berat, perProduk, status FROM opsi WHERE idGrupOpsi = ? AND idOpsi = ?"

	err := con.QueryRow(query, idGrupOpsi, idOpsi).Scan(
		&opsi.idOpsi, &opsi.opsi, &opsi.harga, &opsi.berat, &opsi.perProduk, &opsi.status)

	defer con.Close()
	return opsi, err
}

// GetOpsis is func
func (opsi Opsi) GetOpsis(idGrupOpsi string) []Opsi {
	con := db.Connect()
	query := "SELECT idOpsi, opsi, harga, berat, perProduk, status FROM opsi WHERE idGrupOpsi = ?"
	rows, _ := con.Query(query, idGrupOpsi)

	var opsis []Opsi

	for rows.Next() {
		rows.Scan(
			&opsi.idOpsi, &opsi.opsi, &opsi.harga, &opsi.berat, &opsi.perProduk, &opsi.status,
		)

		opsis = append(opsis, opsi)
	}

	defer con.Close()

	return opsis
}

// CreateUpdateOpsi is func
func (opsi Opsi) CreateUpdateOpsi(idGrupOpsi string) error {
	con := db.Connect()
	query := "INSERT INTO opsi (idOpsi, idGrupOpsi, opsi, harga, berat, perProduk, status) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE opsi = ?, harga = ?, berat = ?, perProduk = ?, status = ?"
	_, err := con.Exec(query, opsi.idOpsi, idGrupOpsi, opsi.opsi, opsi.harga, opsi.berat, opsi.perProduk, opsi.status, opsi.opsi, opsi.harga, opsi.berat, opsi.perProduk, opsi.status)

	defer con.Close()

	return err
}

// DeleteOpsi is func
func (opsi Opsi) DeleteOpsi(idGrupOpsi, idOpsi string) error {
	con := db.Connect()
	query := "DELETE FROM opsi WHERE idGrupOpsi = ? AND idOpsi = ?"
	_, err := con.Exec(query, idGrupOpsi, idOpsi)

	defer con.Close()

	return err
}
