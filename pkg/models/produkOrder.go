package models

import "zonart/db"

// ProdukOrder is class
type ProdukOrder struct {
	IDProdukOrder    int    `json:"idProdukOrder"`
	IDOrder          int    `json:"idOrder"`
	NamaProduk       string `json:"namaProduk"`
	BeratProduk      int    `json:"beratProduk"`
	HargaProduk      int    `json:"hargaProduk"`
	HargaSatuanWajah int    `json:"hargaSatuanWajah"`
}

// GetProdukOrder is func
func (po ProdukOrder) GetProdukOrder(idOrder string) (ProdukOrder, error) {
	con := db.Connect()
	query := "SELECT idProdukOrder, idOrder, namaProduk, beratProduk, hargaProduk, hargaSatuanWajah " +
		"FROM produkOrder WHERE idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&po.IDProdukOrder, &po.IDOrder, &po.NamaProduk, &po.BeratProduk, &po.HargaProduk, &po.HargaSatuanWajah,
	)

	defer con.Close()

	return po, err
}

// CreateProdukOrder is func
func (po ProdukOrder) CreateProdukOrder(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO produkOrder (idOrder, namaProduk, beratProduk, hargaProduk, hargaSatuanWajah) " +
		"VALUES (?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, po.NamaProduk, po.BeratProduk, po.HargaProduk, po.HargaSatuanWajah)

	defer con.Close()

	return err
}
