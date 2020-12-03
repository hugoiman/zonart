package models

// Produk is class
type Produk struct {
	IDProduk       int              `json:"idProduk"`
	IDToko         int              `json:"idToko"`
	NamaProduk     string           `json:"namaProduk"`
	GambarProduk   string           `json:"gambarProduk"`
	Deskripsi      string           `json:"deskripsi"`
	HargaCetak     string           `json:"hargaCetak"`
	HargaSoftCopy  string           `json:"hargaSoftCopy"`
	Status         bool             `json:"status"`
	Catatan        string           `json:"catatan"`
	HargaWajah     int              `json:"hargaWajah"`
	GrupOpsiProduk []GrupOpsiProduk `json:"grupOpsiProduk"`
}
