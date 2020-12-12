package models

import "zonart/db"

// Pengiriman is class
type Pengiriman struct {
	IDPengiriman int    `json:""`
	IDOrder      int    `json:""`
	Penerima     string `json:""`
	Telp         string `json:""`
	Alamat       string `json:""`
	Kota         string `json:""`
	Label        string `json:""`
	Berat        int    `json:""`
	Kurir        string `json:""`
	Resi         string `json:""`
	Ongkir       int    `json:""`
}

// CreatePengiriman is func
func (p Pengiriman) CreatePengiriman(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO order (idOrder, penerima, telp, alamat, kota, label, berat, kurir, ongkir) VALUES (?,?,?,?,?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, p.Penerima, p.Telp, p.Alamat, p.Kota, p.Label, p.Berat, p.Kurir, p.Ongkir)

	defer con.Close()

	return err
}
