package models

import "zonart/db"

// OpsiOrder is class
type OpsiOrder struct {
	IDOpsiOrder int    `json:"idOpsiOrder"`
	NamaGrup    string `json:"namaGrup" validate:"required"`
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

// GetOpsiOrder is func
func (oo OpsiOrder) GetOpsiOrder(idOrder string) []OpsiOrder {
	con := db.Connect()
	var opsiOrders []OpsiOrder

	query := "SELECT idOpsiOrder, namaGrup, opsi, harga, berat, perProduk FROM opsiOrder WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&oo.IDOpsiOrder, &oo.NamaGrup, &oo.Opsi, &oo.Harga, &oo.Berat, &oo.PerProduk,
		)

		opsiOrders = append(opsiOrders, oo)
	}

	defer con.Close()

	return opsiOrders
}
