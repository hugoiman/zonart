package models

// Penangan is class
type Penangan struct {
	IDPenangan   int    `json:"idPenangan"`
	IDKaryawan   int    `json:"idKaryawan"`
	NamaKaryawan string `json:"namaKaryawan"`
}
