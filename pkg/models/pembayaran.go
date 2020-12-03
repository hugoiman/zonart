package models

// Pembayaran is class
type Pembayaran struct {
	IDPembayaran int    `json:"idPembayaran"`
	Bukti        string `json:"bukti"`
	Nominal      int    `json:"nominal"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
}
