package models

// Penggajian is class
type Penggajian struct {
	IDPenggajian int    `json:"idPenggajian"`
	IDKaryawan   int    `json:"idKaryawan"`
	NamaKaryawan string `json:"namaKaryawan"`
	Periode      string `json:"periode"`
	Nominal      string `json:"nominal"`
	TglTransaksi string `json:"tglTransaksi"`
}
