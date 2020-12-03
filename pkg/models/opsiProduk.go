package models

// OpsiProduk is class
type OpsiProduk struct {
	IDOpsiProduk int    `json:"idOpsiProduk"`
	Opsi         string `json:"opsi"`
	Harga        int    `json:"harga"`
	Berat        int    `json:"berat"`
}
