package models

// GrupOpsiProduk is class
type GrupOpsiProduk struct {
	IDGrupOpsi      int          `json:"idGrupOpsi"`
	NamaGrup        string       `json:"namaGrup"`
	Required        bool         `json:"required"`
	Min             int          `json:"min"`
	Max             int          `json:"max"`
	SpesificRequest bool         `json:"spesificRequest"`
	OpsiProduk      []OpsiProduk `json:"opsiProduk"`
}
