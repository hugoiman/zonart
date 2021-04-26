package models

import (
	"strconv"
	"time"
	"zonart/db"
)

// Order is class
type Order struct {
	IDOrder         int             `json:"idOrder"`                                                //y
	Pemesan         int             `json:"pemesan"`                                                //y
	JenisPesanan    string          `json:"jenisPesanan" validate:"required,eq=cetak|eq=soft copy"` //y
	TambahanWajah   int             `json:"tambahanWajah"`                                          //y
	Catatan         string          `json:"catatan"`                                                //y
	Pcs             int             `json:"pcs" validate:"required,min=1"`                          //y
	RencanaPakai    string          `json:"rencanaPakai"`                                           //y
	WaktuPengerjaan string          `json:"waktuPengerjaan"`                                        //y
	ContohGambar    string          `json:"contohGambar"`                                           //y
	TglOrder        string          `json:"tglOrder"`                                               //y
	Invoice         Invoice         `json:"invoice"`
	ProdukOrder     ProdukOrder     `json:"produkOrder"`
	Pengiriman      Pengiriman      `json:"pengiriman" validate:"required,dive"`
	Penangan        Penangan        `json:"penangan"`
	HasilOrder      HasilOrder      `json:"hasilOrder"`
	FileOrder       []FileOrder     `json:"fileOrder"`
	BiayaTambahan   []BiayaTambahan `json:"biayaTambahan"`
	OpsiOrder       []OpsiOrder     `json:"opsiOrder"`
	Revisi          []Revisi        `json:"revisi"`
	Pembayaran      []Pembayaran    `json:"pembayaran"`
}

// Orders is order list
type Orders struct {
	Orders []Order `json:"order"`
}

// GetOrder is func
func (o Order) GetOrder(idOrder string) (Order, error) {
	con := db.Connect()
	query := "SELECT idOrder, idCustomer, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, tglOrder " +
		"FROM `order` WHERE idOrder = ?"

	var tglOrder time.Time

	err := con.QueryRow(query, idOrder).Scan(
		&o.IDOrder, &o.Pemesan,
		&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
		&o.WaktuPengerjaan, &o.ContohGambar, &tglOrder,
	)

	o.TglOrder = tglOrder.Format("02 Jan 2006")

	o.HasilOrder, _ = o.HasilOrder.GetHasilOrder(idOrder)
	o.Invoice, _ = o.Invoice.GetInvoiceByOrder(idOrder)
	o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(idOrder)
	o.Penangan, _ = o.Penangan.GetPenangan(idOrder)
	o.Pengiriman, _ = o.Pengiriman.GetPengiriman(idOrder)

	var opsiOrder OpsiOrder
	o.OpsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	var fileOrder FileOrder
	o.FileOrder = fileOrder.GetFileOrder(idOrder)

	var revisi Revisi
	o.Revisi = revisi.GetRevisi(idOrder)

	var bt BiayaTambahan
	o.BiayaTambahan = bt.GetBiayaTambahans(idOrder)

	var pembayaran Pembayaran
	o.Pembayaran = pembayaran.GetPembayarans(idOrder)

	defer con.Close()
	return o, err
}

// GetOrderToko is func
func (o Order) GetOrderToko(idOrder, idToko string) (Order, error) {
	con := db.Connect()
	query := "SELECT idOrder, idCustomer," +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, tglOrder " +
		"FROM `order` WHERE idOrder = ? AND idToko = ?"

	var tglOrder time.Time
	err := con.QueryRow(query, idOrder, idToko).Scan(
		&o.IDOrder, &o.Pemesan,
		&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
		&o.WaktuPengerjaan, &o.ContohGambar, &tglOrder,
	)

	o.TglOrder = tglOrder.Format("02 Jan 2006")

	o.HasilOrder, _ = o.HasilOrder.GetHasilOrder(idOrder)
	o.Invoice, _ = o.Invoice.GetInvoiceByOrder(idOrder)
	o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(idOrder)
	o.Penangan, _ = o.Penangan.GetPenangan(idOrder)
	o.Pengiriman, _ = o.Pengiriman.GetPengiriman(idOrder)

	var opsiOrder OpsiOrder
	o.OpsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	var fileOrder FileOrder
	o.FileOrder = fileOrder.GetFileOrder(idOrder)

	var revisi Revisi
	o.Revisi = revisi.GetRevisi(idOrder)

	var bt BiayaTambahan
	o.BiayaTambahan = bt.GetBiayaTambahans(idOrder)

	var pembayaran Pembayaran
	o.Pembayaran = pembayaran.GetPembayarans(idOrder)

	defer con.Close()
	return o, err
}

// GetOrderCustomer is func
func (o Order) GetOrderCustomer(idOrder, idCustomer string) (Order, error) {
	con := db.Connect()
	query := "SELECT idOrder, idCustomer, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, tglOrder " +
		"FROM `order` WHERE idOrder = ? AND idCustomer = ?"

	var tglOrder time.Time
	err := con.QueryRow(query, idOrder, idCustomer).Scan(
		&o.IDOrder, &o.Pemesan,
		&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
		&o.WaktuPengerjaan, &o.ContohGambar, &tglOrder,
	)

	o.TglOrder = tglOrder.Format("02 Jan 2006")

	o.HasilOrder, _ = o.HasilOrder.GetHasilOrder(idOrder)
	o.Invoice, _ = o.Invoice.GetInvoice(idOrder)
	o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(idOrder)
	o.Penangan, _ = o.Penangan.GetPenangan(idOrder)
	o.Pengiriman, _ = o.Pengiriman.GetPengiriman(idOrder)

	var opsiOrder OpsiOrder
	o.OpsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	var fileOrder FileOrder
	o.FileOrder = fileOrder.GetFileOrder(idOrder)

	var revisi Revisi
	o.Revisi = revisi.GetRevisi(idOrder)

	var bt BiayaTambahan
	o.BiayaTambahan = bt.GetBiayaTambahans(idOrder)

	var pembayaran Pembayaran
	o.Pembayaran = pembayaran.GetPembayarans(idOrder)

	defer con.Close()
	return o, err
}

// GetOrderByInvoice is func
func (o Order) GetOrderByInvoice(idInvoice string) (Order, error) {
	con := db.Connect()
	query := "SELECT idOrder, idCustomer, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, tglOrder " +
		"FROM `order` WHERE idInvoice = ?"

	var tglOrder time.Time

	err := con.QueryRow(query, idInvoice).Scan(
		&o.IDOrder, &o.Pemesan,
		&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
		&o.WaktuPengerjaan, &o.ContohGambar, &tglOrder,
	)

	o.TglOrder = tglOrder.Format("02 Jan 2006")

	o.Invoice, _ = o.Invoice.GetInvoice(idInvoice)
	o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(strconv.Itoa(o.IDOrder))
	o.Pengiriman, _ = o.Pengiriman.GetPengiriman(strconv.Itoa(o.IDOrder))

	var opsiOrder OpsiOrder
	o.OpsiOrder = opsiOrder.GetOpsiOrder(strconv.Itoa(o.IDOrder))

	var bt BiayaTambahan
	o.BiayaTambahan = bt.GetBiayaTambahans(strconv.Itoa(o.IDOrder))

	defer con.Close()
	return o, err
}

// GetOrders is func
func (o Order) GetOrders(idCustomer string) Orders {
	con := db.Connect()
	query := "SELECT idOrder, idCustomer, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, tglOrder " +
		"FROM `order` WHERE idCustomer = ?"
	rows, _ := con.Query(query, idCustomer)

	var orders Orders
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &o.Pemesan,
			&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
			&o.WaktuPengerjaan, &o.ContohGambar, &tglOrder,
		)

		o.TglOrder = tglOrder.Format("02 Jan 2006")

		o.HasilOrder, _ = o.HasilOrder.GetHasilOrder(strconv.Itoa(o.IDOrder))
		o.Invoice, _ = o.Invoice.GetInvoiceByOrder(strconv.Itoa(o.IDOrder))
		o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(strconv.Itoa(o.IDOrder))

		var penangan Penangan
		o.Penangan, _ = penangan.GetPenangan(strconv.Itoa(o.IDOrder))

		var pengiriman Pengiriman
		o.Pengiriman, _ = pengiriman.GetPengiriman(strconv.Itoa(o.IDOrder))

		var opsiOrder OpsiOrder
		o.OpsiOrder = opsiOrder.GetOpsiOrder(strconv.Itoa(o.IDOrder))

		var fileOrder FileOrder
		o.FileOrder = fileOrder.GetFileOrder(strconv.Itoa(o.IDOrder))

		var revisi Revisi
		o.Revisi = revisi.GetRevisi(strconv.Itoa(o.IDOrder))

		var bt BiayaTambahan
		o.BiayaTambahan = bt.GetBiayaTambahans(strconv.Itoa(o.IDOrder))

		var pembayaran Pembayaran
		o.Pembayaran = pembayaran.GetPembayarans(strconv.Itoa(o.IDOrder))

		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersToko is func
func (o Order) GetOrdersToko(idToko string) Orders {
	con := db.Connect()
	query := "SELECT idOrder, tglOrder FROM `order` WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var orders Orders
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &tglOrder,
		)

		o.TglOrder = tglOrder.Format("02 Jan 2006")
		o.Invoice, _ = o.Invoice.GetInvoiceByOrder(strconv.Itoa(o.IDOrder))
		o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(strconv.Itoa(o.IDOrder))

		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersEditor is func
func (o Order) GetOrdersEditor(idToko, idKaryawan string) Orders {
	con := db.Connect()
	query := "SELECT a.idOrder, a.tglOrder " +
		"FROM `order` a JOIN penangan b ON a.idOrder = b.idOrder " +
		"WHERE a.idToko = ? AND b.idKaryawan = ?"
	rows, _ := con.Query(query, idToko, idKaryawan)

	var orders Orders
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &tglOrder,
		)

		o.TglOrder = tglOrder.Format("02 Jan 2006")
		o.Invoice, _ = o.Invoice.GetInvoiceByOrder(strconv.Itoa(o.IDOrder))
		o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(strconv.Itoa(o.IDOrder))

		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
}

// CreateOrder is func
func (o Order) CreateOrder(idToko, idProduk string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO `order` (idToko, idProduk, idCustomer, idInvoice, jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, contohGambar, tglOrder) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, idProduk, o.Pemesan, o.Invoice.IDInvoice, o.JenisPesanan, o.TambahanWajah, o.Catatan, o.Pcs,
		o.RencanaPakai, o.ContohGambar, o.TglOrder)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idOrder := int(idInt64)

	err = o.Pengiriman.CreatePengiriman(strconv.Itoa(idOrder))
	if err != nil {
		return 0, err
	}

	err = o.ProdukOrder.CreateProdukOrder(strconv.Itoa(idOrder))
	if err != nil {
		return 0, err
	}

	for _, v := range o.FileOrder {
		err = v.CreateFileOrder(strconv.Itoa(idOrder))
		if err != nil {
			return 0, err
		}
	}

	for _, v := range o.OpsiOrder {
		err = v.CreateOpsiOrder(strconv.Itoa(idOrder))
		if err != nil {
			return 0, err
		}
	}

	defer con.Close()

	return idOrder, err
}

// // UpdateStatusOrder is func
// func (o Order) UpdateStatusOrder(idOrder string) error {
// 	con := db.Connect()
// 	query := "UPDATE `order` SET statusPesanan = ?, statusPembayaran = ? WHERE idOrder = ?"
// 	_, err := con.Exec(query, o.StatusPesanan, o.StatusPembayaran, idOrder)

// 	defer con.Close()

// 	return err
// }

// SetWaktuPengerjaan is func
func (o Order) SetWaktuPengerjaan(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET waktuPengerjaan = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.WaktuPengerjaan, idOrder)

	defer con.Close()

	return err
}

// // FinishOrder is func
// func (o Order) FinishOrder(idOrder string) error {
// 	con := db.Connect()
// 	query := "UPDATE `order` SET statusPesanan = ? WHERE idOrder = ?"
// 	_, err := con.Exec(query, o.StatusPesanan, idOrder)

// 	defer con.Close()

// 	return err
// }
