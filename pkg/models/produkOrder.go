package models

import "zonart/db"

// ProdukOrder is class
type ProdukOrder struct {
	IDProdukOrder    int    `json:"idProdukOrder"`
	NamaProduk       string `json:"namaProduk"`
	BeratProduk      int    `json:"beratProduk"`
	HargaProduk      int    `json:"hargaProduk"`
	HargaSatuanWajah int    `json:"hargaSatuanWajah"`
	FotoProduk       string `json:"fotoProduk"`
	SlugProduk       string `json:"slugProduk"`
}

// GetProdukOrder is func
func (po ProdukOrder) GetProdukOrder(idOrder string) (ProdukOrder, error) {
	con := db.Connect()
	query := "SELECT a.idProdukOrder, a.namaProduk, a.beratProduk, a.hargaProduk, a.hargaSatuanWajah, c.gambar, c.slug " +
		"FROM produkOrder a JOIN `order` b ON a.idOrder = b.idOrder " +
		"JOIN produk c ON b.idProduk = c.idProduk WHERE a.idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&po.IDProdukOrder, &po.NamaProduk, &po.BeratProduk, &po.HargaProduk, &po.HargaSatuanWajah, &po.FotoProduk, &po.SlugProduk,
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
