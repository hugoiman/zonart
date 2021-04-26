package models

import "zonart/db"

// Pengiriman is class
type Pengiriman struct {
	IDPengiriman int    `json:"idPengiriman"`
	Penerima     string `json:"penerima" validate:"required"`
	Telp         string `json:"telp" validate:"required"`
	Alamat       string `json:"alamat"`
	Kota         string `json:"kota"`
	Label        string `json:"label"`
	Berat        int    `json:"berat"`
	Kurir        string `json:"kurir"`
	KodeKurir    string `json:"kodeKurir"`
	Service      string `json:"service"`
	Estimasi     string `json:"estimasi"`
	Resi         string `json:"resi"`
	Ongkir       int    `json:"ongkir"`
}

// CreatePengiriman is func
func (p Pengiriman) CreatePengiriman(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO pengiriman (idOrder, penerima, telp, alamat, kota, label, berat, kurir, kodeKurir, service, estimasi, ongkir) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, p.Penerima, p.Telp, p.Alamat, p.Kota, p.Label, p.Berat, p.Kurir, p.KodeKurir, p.Service, p.Estimasi, p.Ongkir)

	defer con.Close()

	return err
}

// GetPengiriman is func
func (p Pengiriman) GetPengiriman(idOrder string) (Pengiriman, error) {
	con := db.Connect()
	query := "SELECT idPengiriman, penerima, telp, alamat, kota, label, berat, kurir, kodeKurir, service, estimasi, resi, ongkir FROM pengiriman WHERE idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&p.IDPengiriman, &p.Penerima, &p.Telp, &p.Alamat, &p.Kota, &p.Label, &p.Berat, &p.Kurir, &p.KodeKurir, &p.Service, &p.Estimasi, &p.Resi, &p.Ongkir)

	defer con.Close()
	return p, err
}

// SetResi is func
func (p Pengiriman) SetResi(idOrder string) error {
	con := db.Connect()
	query := "UPDATE pengiriman SET resi = ? WHERE idOrder = ?"
	_, err := con.Exec(query, p.Resi, idOrder)

	defer con.Close()

	return err
}

// UpdateBeratOngkir is func
func (p Pengiriman) UpdateBeratOngkir(idOrder string, berat int, ongkir int) error {
	con := db.Connect()
	query := "UPDATE pengiriman SET berat = ?, ongkir = ? WHERE idOrder = ?"
	_, err := con.Exec(query, berat, ongkir, idOrder)

	defer con.Close()

	return err
}
