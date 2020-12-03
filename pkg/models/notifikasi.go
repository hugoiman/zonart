package models

// Notifikasi is class
type Notifikasi struct {
	IDNotifikasi int    `json:"idNotifikasi"`
	IDPenerima   int    `json:"idPenerima"`
	Pengirim     string `json:"pengirim"`
	Judul        string `json:"judul"`
	Pesan        string `json:"pesan"`
	Link         string `json:"link"`
	Read         bool   `json:"read"`
	CreatedAt    string `json:"createdAt"`
}
