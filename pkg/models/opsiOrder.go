package models

// OpsiOrder is class
type OpsiOrder struct {
	IDOpsiOrder int    `json:"idOpsiOrder"`
	NamaGrup    string `json:"namaGrup"`
	Opsi        string `json:"opsi"`
	Harga       int    `json:"harga"`
	Berat       int    `json:"berat"`
}
