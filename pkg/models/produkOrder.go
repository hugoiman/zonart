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
	FotoProduk 		 string `json:"fotoProduk"`
	SlugProduk 		 string `json:"slugProduk"`
}

// GetProdukOrder is func
func (po ProdukOrder) GetProdukOrder(idOrder, idProduk string) (ProdukOrder, error) {
	con := db.Connect()
	query := "SELECT a.idProdukOrder, a.idOrder, a.namaProduk, a.beratProduk, a.hargaProduk, a.hargaSatuanWajah, b.gambar, b.slug " +
		"FROM produkOrder a, produk b WHERE a.idOrder = ? AND b.idProduk = ?"

	err := con.QueryRow(query, idOrder, idProduk).Scan(
		&po.IDProdukOrder, &po.IDOrder, &po.NamaProduk, &po.BeratProduk, &po.HargaProduk, &po.HargaSatuanWajah, &po.FotoProduk, &po.SlugProduk,
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
