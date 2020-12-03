package models

// Order is class
type Order struct {
	IDOrder          int             `json:"idOrder"`
	IDToko           int             `json:"idToko"`
	IDProduk         int             `json:"idProduk"`
	IDCustomer       int             `json:"idCustomer"`
	NamaToko         string          `json:"namaToko"`
	SlugToko         string          `json:"slugToko"`
	NamaProduk       string          `json:"namaProduk"`
	NamaCustomer     string          `json:"namaCustomer"`
	JenisPesanan     string          `json:"jenisPesanan"`
	Catatan          string          `json:"catatan"`
	TambahanWajah    int             `json:"tambahanWajah"`
	TotalHargaWajah  int             `json:"totalHargaWajah"`
	Pcs              string          `json:"pcs"`
	StatusPesanan    string          `json:"statusPesanan"`
	StatusPembayaran string          `json:"statusPembayaran"`
	Total            int             `json:"total"`
	Dibayar          int             `json:"dibayar"`
	Tagihan          int             `json:"tagihan"`
	RencanaPakai     string          `json:"rencanaPakai"`
	Gambar           string          `json:"gambar"`
	Hasil            string          `json:"hasil"`
	CreatedAt        string          `json:"createdAt"`
	Pengiriman       Pengiriman      `json:"pengiriman"`
	Penangan         Penangan        `json:"penangan"`
	TambahanBiaya    []TambahanBiaya `json:"tambahanBiaya"`
	Pembayaran       []Pembayaran    `json:"pembayaran"`
	OpsiOrder        []OpsiOrder     `json:"opsiOrder"`
	Revisi           []Revisi        `json:"revisi"`
}
