package models

import (
	"strconv"
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

// Orders is order list
type Orders struct {
	Orders []Order `json:"order"`
}

// GetOrder is func
func (o Order) GetOrder(idOrder string) (Order, error) {
	con := db.Connect()
	query := "SELECT a.idOrder, a.idToko, a.idProduk, a.idCustomer, b.namaToko, b.slug, a.namaProduk, c.nama," +
		"a.jenisPesanan, a.hargaProduk, a.tambahanWajah, a.hargaWajah, a.catatan, a.pcs," +
		"a.statusPesanan, a.statusPembayaran, a.total, a.dibayar, a.tagihan, a.rencanaPakai, a.waktuPengerjaan, a.gambar, a.hasil, a.createdAt FROM `order` a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer " +
		"WHERE a.idOrder = ?"

	var createdAt time.Time

	err := con.QueryRow(query, idOrder).Scan(
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

// GetOrders is func
func (o Order) GetOrders(idCustomer string) Orders {
	con := db.Connect()
	query := "SELECT a.idOrder, a.idToko, a.idProduk, a.idCustomer, b.namaToko, b.slug, a.namaProduk, c.nama," +
		"a.jenisPesanan, a.hargaProduk, a.tambahanWajah, a.hargaWajah, a.catatan, a.pcs," +
		"a.statusPesanan, a.statusPembayaran, a.total, a.dibayar, a.tagihan, a.rencanaPakai, a.waktuPengerjaan, a.gambar, a.hasil, a.createdAt FROM `order` a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer " +
		"WHERE a.idCustomer = ?"
	rows, _ := con.Query(query, idCustomer)

	var orders Orders
	var createdAt time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.NamaToko, &o.SlugToko, &o.NamaProduk, &o.NamaCustomer,
			&o.JenisPesanan, &o.HargaProduk, &o.TambahanWajah, &o.HargaWajah, &o.Catatan, &o.Pcs, &o.StatusPesanan, &o.StatusPembayaran,
			&o.Total, &o.Dibayar, &o.Tagihan, &o.RencanaPakai, &o.WaktuPengerjaan, &o.Gambar, &o.Hasil, &createdAt,
		)

		o.CreatedAt = createdAt.Format("02 Jan 2006")
		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersToko is func
func (o Order) GetOrdersToko(idToko string) Orders {
	con := db.Connect()
	query := "SELECT a.idOrder, a.idToko, a.idProduk, a.idCustomer, b.namaToko, b.slug, a.namaProduk, c.nama," +
		"a.jenisPesanan, a.hargaProduk, a.tambahanWajah, a.hargaWajah, a.catatan, a.pcs," +
		"a.statusPesanan, a.statusPembayaran, a.total, a.dibayar, a.tagihan, a.rencanaPakai, a.waktuPengerjaan, a.gambar, a.hasil, a.createdAt FROM `order` a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer " +
		"WHERE c.idToko = ?"
	rows, _ := con.Query(query, idToko)

	var orders Orders
	var createdAt time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.NamaToko, &o.SlugToko, &o.NamaProduk, &o.NamaCustomer,
			&o.JenisPesanan, &o.HargaProduk, &o.TambahanWajah, &o.HargaWajah, &o.Catatan, &o.Pcs, &o.StatusPesanan, &o.StatusPembayaran,
			&o.Total, &o.Dibayar, &o.Tagihan, &o.RencanaPakai, &o.WaktuPengerjaan, &o.Gambar, &o.Hasil, &createdAt,
		)

		o.CreatedAt = createdAt.Format("02 Jan 2006")
		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersEditor is func
func (o Order) GetOrdersEditor(idToko, idPenangan string) Orders {
	con := db.Connect()
	query := "SELECT a.idOrder, a.idToko, a.idProduk, a.idCustomer, b.namaToko, b.slug, a.namaProduk, c.nama," +
		"a.jenisPesanan, a.hargaProduk, a.tambahanWajah, a.hargaWajah, a.catatan, a.pcs," +
		"a.statusPesanan, a.statusPembayaran, a.total, a.dibayar, a.tagihan, a.rencanaPakai, a.waktuPengerjaan, a.gambar, a.hasil, a.createdAt FROM `order` a " +
		"JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer " +
		"JOIN penangan d ON a.idOrder = d.idOrder" +
		"WHERE a.idOrder = ? AND d.idPenangan = ?"
	rows, _ := con.Query(query, idToko, idPenangan)

	var orders Orders
	var createdAt time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.NamaToko, &o.SlugToko, &o.NamaProduk, &o.NamaCustomer,
			&o.JenisPesanan, &o.HargaProduk, &o.TambahanWajah, &o.HargaWajah, &o.Catatan, &o.Pcs, &o.StatusPesanan, &o.StatusPembayaran,
			&o.Total, &o.Dibayar, &o.Tagihan, &o.RencanaPakai, &o.WaktuPengerjaan, &o.Gambar, &o.Hasil, &createdAt,
		)

		o.CreatedAt = createdAt.Format("02 Jan 2006")
		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
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

	for _, vOpsiOrder := range o.OpsiOrder {
		err = vOpsiOrder.CreateOpsiOrder(strconv.Itoa(idOrder))
		if err != nil {
			return 0, err
		}
	}

	if o.JenisPesanan == "cetak" {
		err = o.Pengiriman.CreatePengiriman(strconv.Itoa(idOrder))
		if err != nil {
			return 0, err
		}
	}

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

// UpdateBiayaOrder is func
func (o Order) UpdateBiayaOrder(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET total = ?, tagihan = ?,  statusPembayaran = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.Total, o.Tagihan, o.StatusPembayaran, idOrder)

	defer con.Close()

	return err
}

// KonfirmasiOrder is func
func (o Order) KonfirmasiOrder(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET tagihan = ?, statusPembayaran = ?, statusPesanan = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.Tagihan, o.StatusPembayaran, o.StatusPesanan, idOrder)

	defer con.Close()

	return err
}

// SetWaktuPengerjaan is func
func (o Order) SetWaktuPengerjaan(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET waktuPengerjaan = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.WaktuPengerjaan, idOrder)

	defer con.Close()

	return err
}

// UploadHasilProduksi is func
func (o Order) UploadHasilProduksi(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET hasil = ?, statusPesanan = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.Hasil, o.StatusPesanan, idOrder)

	defer con.Close()

	return err
}

// SetujuiHasilProduksi is func
func (o Order) SetujuiHasilProduksi(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET statusPesanan = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.StatusPesanan, idOrder)

	defer con.Close()

	return err
}

// CancelOrder is func
func (o Order) CancelOrder(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET statusPesanan = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.StatusPesanan, idOrder)

	defer con.Close()

	return err
}

// FinishOrder is func
func (o Order) FinishOrder(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET statusPesanan = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.StatusPesanan, idOrder)

	defer con.Close()

	return err
}
