package models

// Pembukuan is class
type Pembukuan struct {
	IDPembukuan  int    `json:"idPembukuan"`
	IDToko       int    `json:"idToko"`
	Jenis        string `json:"jenis"`
	Keterangan   string `json:"keterangan"`
	Nominal      int    `json:"nominal"`
	TglTransaksi string `json:"tglTransaksi"`
}
