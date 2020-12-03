package models

// Pengiriman is class
type Pengiriman struct {
	IDPengiriman int    `json:""`
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
