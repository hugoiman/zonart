package models

import "zonart/db"

// OpsiOrder is class
type OpsiOrder struct {
	IDOpsiOrder int    `json:"idOpsiOrder"`
	IDOrder     int    `json:"idOrder"`
	IDOpsi      int    `json:"idOpsi"`
	IDGrupOpsi  int    `json:"idGrupOpsi"`
	NamaGrup    string `json:"namaGrup"`
	Opsi        string `json:"opsi"`
	Harga       int    `json:"harga"`
	Berat       int    `json:"berat"`
	PerProduk   bool   `json:"perProduk"`
}

// CreateOpsiOrder is func
func (oo OpsiOrder) CreateOpsiOrder(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO opsiOrder (idOrder, namaGrup, opsi, harga, berat, perProduk) VALUES (?,?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, oo.NamaGrup, oo.Opsi, oo.Harga, oo.Berat, oo.PerProduk)

	defer con.Close()

	return err
}
