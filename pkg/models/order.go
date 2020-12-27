package models

import (
	"time"
	"zonart/db"
)

// Order is class
type Order struct {
	IDOrder            int             `json:"idOrder"`    //y
	IDToko             int             `json:"idToko"`     //y
	IDProduk           int             `json:"idProduk"`   //y
	IDCustomer         int             `json:"idCustomer"` //y
	NamaToko           string          `json:"namaToko"`
	SlugToko           string          `json:"slugToko"`
	NamaProduk         string          `json:"namaProduk"` //y
	NamaCustomer       string          `json:"namaCustomer"`
	JenisPesanan       string          `json:"jenisPesanan" validate:"required"` //y
	HargaProduk        int             `json:"hargaProduk"`                      // getProduk then get harga produk cetak/softcopy //y
	TambahanWajah      int             `json:"tambahanWajah"`                    //y
	HargaWajah         int             `json:"hargaWajah"`                       //y
	TotalHargaWajah    int             `json:"totalHargaWajah"`                  // get harga per wajah dari produk * tambahan wajah
	Catatan            string          `json:"catatan"`                          //y
	Pcs                int             `json:"pcs" validate:"required,min=1"`    //y
	StatusPesanan      string          `json:"statusPesanan"`                    //y
	StatusPembayaran   string          `json:"statusPembayaran"`                 //y
	TotalHargaOpsi     int             `json:"totalHargaOpsi"`
	TotalTambahanBiaya int             `json:"totalTambayanBiaya"`
	Total              int             `json:"total"`                      // (harga produk * pcs) + total harga wajah + total harga opsi + tambahan biaya //y
	Dibayar            int             `json:"dibayar"`                    //y
	Tagihan            int             `json:"tagihan"`                    // total - dibayar //y
	RencanaPakai       string          `json:"rencanaPakai"`               //y
	WaktuPengerjaan    string          `json:"waktuPengerjaan"`            //y
	Gambar             string          `json:"gambar" validate:"required"` //y
	Hasil              string          `json:"hasil"`                      //y
	TotalBeratOpsi     int             `json:"totalBeratOpsi"`
	CreatedAt          string          `json:"createdAt"` //y
	Pengiriman         Pengiriman      `json:"pengiriman"`
	Penangan           Penangan        `json:"penangan"`
	BiayaTambahan      []BiayaTambahan `json:"biayaTambahan"`
	Pembayaran         []Pembayaran    `json:"pembayaran"`
	OpsiOrder          []OpsiOrder     `json:"opsiOrder"`
	Revisi             []Revisi        `json:"revisi"`
}

// GetOrder is func
func (o Order) GetOrder(idOrder, idCustomer string) (Order, error) {
	con := db.Connect()
	query := "SELECT a.idOrder, a.idToko, a.idProduk, a.idCustomer, b.namaToko, b.slug, a.namaProduk, c.nama," +
		"a.jenisPesanan, a.hargaProduk, a.tambahanWajah, a.hargaWajah, a.catatan, a.pcs," +
		"a.statusPesanan, a.statusPembayaran, a.total, a.dibayar, a.tagihan, a.rencanaPakai, a.waktuPengerjaan, a.gambar, a.hasil, a.createdAt FROM `order` a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer " +
		"WHERE a.idOrder = ? AND a.idCustomer = ?"

	var createdAt time.Time

	err := con.QueryRow(query, idOrder, idCustomer).Scan(
		&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.NamaToko, &o.SlugToko, &o.NamaProduk, &o.NamaCustomer,
		&o.JenisPesanan, &o.HargaProduk, &o.TambahanWajah, &o.HargaWajah, &o.Catatan, &o.Pcs, &o.StatusPesanan, &o.StatusPembayaran,
		&o.Total, &o.Dibayar, &o.Tagihan, &o.RencanaPakai, &o.WaktuPengerjaan, &o.Gambar, &o.Hasil, &createdAt,
	)

	o.CreatedAt = createdAt.Format("02 Jan 2006")

	dataPengiriman, _ := o.Pengiriman.GetPengiriman(idOrder)
	o.Pengiriman = dataPengiriman

	o.Penangan = o.Penangan.GetPenangan(idOrder)

	var opsiOrder OpsiOrder
	o.OpsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	defer con.Close()
	return o, err
}

// GetOrderToko is func
func (o Order) GetOrderToko(idOrder, idToko string) (Order, error) {
	con := db.Connect()
	query := "SELECT a.idOrder, a.idToko, a.idProduk, a.idCustomer, b.namaToko, b.slug, a.namaProduk, c.nama," +
		"a.jenisPesanan, a.hargaProduk, a.tambahanWajah, a.hargaWajah, a.catatan, a.pcs," +
		"a.statusPesanan, a.statusPembayaran, a.total, a.dibayar, a.tagihan, a.rencanaPakai, a.waktuPengerjaan, a.gambar, a.hasil, a.createdAt FROM `order` a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer " +
		"WHERE a.idOrder = ? AND a.idToko = ?"

	var createdAt time.Time

	err := con.QueryRow(query, idOrder, idToko).Scan(
		&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.NamaToko, &o.SlugToko, &o.NamaProduk, &o.NamaCustomer,
		&o.JenisPesanan, &o.HargaProduk, &o.TambahanWajah, &o.HargaWajah, &o.Catatan, &o.Pcs, &o.StatusPesanan, &o.StatusPembayaran,
		&o.Total, &o.Dibayar, &o.Tagihan, &o.RencanaPakai, &o.WaktuPengerjaan, &o.Gambar, &o.Hasil, &createdAt,
	)

	o.CreatedAt = createdAt.Format("02 Jan 2006")

	dataPengiriman, _ := o.Pengiriman.GetPengiriman(idOrder)
	o.Pengiriman = dataPengiriman

	o.Penangan = o.Penangan.GetPenangan(idOrder)

	var opsiOrder OpsiOrder
	o.OpsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	var bt BiayaTambahan
	o.BiayaTambahan = bt.GetBiayaTambahans(idOrder)

	defer con.Close()
	return o, err
}

// CreateOrder is func
func (o Order) CreateOrder(idToko, idProduk string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO `order` (idToko, idProduk, idCustomer, namaProduk, jenisPesanan, hargaProduk, tambahanWajah, hargaWajah, catatan, pcs, statusPesanan, statusPembayaran, total, dibayar, tagihan, rencanaPakai, gambar, createdAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, idProduk, o.IDCustomer, o.NamaProduk, o.JenisPesanan, o.HargaProduk, o.TambahanWajah, o.HargaWajah, o.Catatan, o.Pcs, o.StatusPesanan, o.StatusPembayaran, o.Total, o.Dibayar, o.Tagihan, o.RencanaPakai, o.Gambar, o.CreatedAt)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idOrder := int(idInt64)

	defer con.Close()

	return idOrder, err
}

// DeleteOrder is func
func (o Order) DeleteOrder(idOrder string) error {
	con := db.Connect()
	query := "DELETE FROM `order` WHERE idOrder = ?"
	_, err := con.Exec(query, idOrder)

	defer con.Close()

	return err
}

// UpdateStatusOrder is func
func (o Order) UpdateStatusOrder(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET statusPesanan = ?, statusPembayaran = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.StatusPesanan, o.StatusPembayaran, idOrder)

	defer con.Close()

	return err
}
