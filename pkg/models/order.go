package models

import (
	"strconv"
	"time"
	"zonart/db"
)

// Order is class
type Order struct {
	IDOrder         int             `json:"idOrder"`    //y
	IDToko          int             `json:"idToko"`     //y
	IDProduk        int             `json:"idProduk"`   //y
	IDCustomer      int             `json:"idCustomer"` //y
	IDInvoice       int             `json:"idInvoice"`
	JenisPesanan    string          `json:"jenisPesanan" validate:"required"` //y
	TambahanWajah   int             `json:"tambahanWajah"`                    //y
	Catatan         string          `json:"catatan"`                          //y
	Pcs             int             `json:"pcs" validate:"required,min=1"`    //y
	RencanaPakai    string          `json:"rencanaPakai"`                     //y
	WaktuPengerjaan string          `json:"waktuPengerjaan"`                  //y
	ContohGambar    string          `json:"contohGambar" validate:"required"` //y
	HasilOrder      string          `json:"hasilOrder"`                       //y
	TglOrder        string          `json:"tglOrder"`                         //y
	Invoice         Invoice         `json:"invoice"`
	ProdukOrder     ProdukOrder     `json:"produkOrder"`
	Pengiriman      Pengiriman      `json:"pengiriman"`
	Penangan        Penangan        `json:"penangan"`
	FileOrder       []FileOrder     `json:"fileOrder" validate:"required,dive"`
	BiayaTambahan   []BiayaTambahan `json:"biayaTambahan"`
	OpsiOrder       []OpsiOrder     `json:"opsiOrder"`
	Revisi          []Revisi        `json:"revisi"`
	Pembayaran		[]Pembayaran	`json:"pembayaran"`
}

// Orders is order list
type Orders struct {
	Orders []Order `json:"order"`
}

// GetOrder is func
func (o Order) GetOrder(idOrder string) (Order, error) {
	con := db.Connect()
	query := "SELECT idOrder, idToko, idProduk, idCustomer, idInvoice, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, hasilOrder, tglOrder " +
		"FROM `order` WHERE idOrder = ?"

	var tglOrder time.Time

	err := con.QueryRow(query, idOrder).Scan(
		&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.IDInvoice,
		&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
		&o.WaktuPengerjaan, &o.ContohGambar, &o.HasilOrder, &tglOrder,
	)

	o.TglOrder = tglOrder.Format("02 Jan 2006")

	o.Invoice, _ = o.Invoice.GetInvoice(strconv.Itoa(o.IDInvoice))
	o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(idOrder)
	o.Penangan = o.Penangan.GetPenangan(idOrder)
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
	query := "SELECT idOrder, idToko, idProduk, idCustomer, idInvoice, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, hasilOrder, tglOrder " +
		"FROM `order` WHERE idInvoice = ?"

	var tglOrder time.Time

	err := con.QueryRow(query, idInvoice).Scan(
		&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.IDInvoice,
		&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
		&o.WaktuPengerjaan, &o.ContohGambar, &o.HasilOrder, &tglOrder,
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
	query := "SELECT idOrder, idToko, idProduk, idCustomer, idInvoice, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, hasilOrder, tglOrder " +
		"FROM `order` WHERE a.idCustomer = ?"
	rows, _ := con.Query(query, idCustomer)

	var orders Orders
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.IDInvoice,
			&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
			&o.WaktuPengerjaan, &o.ContohGambar, &o.HasilOrder, &tglOrder,
		)

		o.TglOrder = tglOrder.Format("02 Jan 2006")
		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersToko is func
func (o Order) GetOrdersToko(idToko string) Orders {
	con := db.Connect()
	query := "SELECT idOrder, idInvoice, tglOrder FROM `order` WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var orders Orders
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &o.IDInvoice, &tglOrder,
		)

		o.TglOrder = tglOrder.Format("02 Jan 2006")
		o.Invoice, _ = o.Invoice.GetInvoice(strconv.Itoa(o.IDInvoice))
		o.ProdukOrder, _ = o.ProdukOrder.GetProdukOrder(strconv.Itoa(o.IDOrder))

		orders.Orders = append(orders.Orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersEditor is func
func (o Order) GetOrdersEditor(idToko, idPenangan string) Orders {
	con := db.Connect()
	query := "SELECT a.idOrder, a.idToko, a.idProduk, a.idCustomer, a.idInvoice, " +
		"a.jenisPesanan, a.tambahanWajah, a.catatan, a.pcs, a.rencanaPakai, a.waktuPengerjaan, a.contohGambar, a.hasilOrder, a.tglOrder " +
		"FROM `order` a JOIN penangan b ON a.idOrder = b.idOrder" +
		"WHERE a.idOrder = ? AND b.idPenangan = ?"
	rows, _ := con.Query(query, idToko, idPenangan)

	var orders Orders
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.IDOrder, &o.IDToko, &o.IDProduk, &o.IDCustomer, &o.IDInvoice,
			&o.JenisPesanan, &o.TambahanWajah, &o.Catatan, &o.Pcs, &o.RencanaPakai,
			&o.WaktuPengerjaan, &o.ContohGambar, &o.HasilOrder, &tglOrder,
		)

		o.TglOrder = tglOrder.Format("02 Jan 2006")
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
	exec, err := con.Exec(query, idToko, idProduk, o.IDCustomer, o.IDInvoice, o.JenisPesanan, o.TambahanWajah, o.Catatan, o.Pcs,
		o.RencanaPakai, o.ContohGambar, o.TglOrder)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idOrder := int(idInt64)

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

	if o.JenisPesanan == "cetak" {
		err = o.Pengiriman.CreatePengiriman(strconv.Itoa(idOrder))
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

// // UpdateBiayaOrder is func
// func (o Order) UpdateBiayaOrder(idOrder string) error {
// 	con := db.Connect()
// 	query := "UPDATE `order` SET total = ?, tagihan = ?,  statusPembayaran = ? WHERE idOrder = ?"
// 	_, err := con.Exec(query, o.Total, o.Tagihan, o.StatusPembayaran, idOrder)

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

// UploadHasilProduksi is func
func (o Order) UploadHasilProduksi(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `order` SET hasilOrder = ? WHERE idOrder = ?"
	_, err := con.Exec(query, o.HasilOrder, idOrder)

	defer con.Close()

	return err
}

// // SetujuiHasilProduksi is func
// func (o Order) SetujuiHasilProduksi(idOrder string) error {
// 	con := db.Connect()
// 	query := "UPDATE `order` SET statusPesanan = ? WHERE idOrder = ?"
// 	_, err := con.Exec(query, o.StatusPesanan, idOrder)

// 	defer con.Close()

// 	return err
// }

// // CancelOrder is func
// func (o Order) CancelOrder(idOrder string) error {
// 	con := db.Connect()
// 	query := "UPDATE `order` SET statusPesanan = ? WHERE idOrder = ?"
// 	_, err := con.Exec(query, o.StatusPesanan, idOrder)

// 	defer con.Close()

// 	return err
// }

// // FinishOrder is func
// func (o Order) FinishOrder(idOrder string) error {
// 	con := db.Connect()
// 	query := "UPDATE `order` SET statusPesanan = ? WHERE idOrder = ?"
// 	_, err := con.Exec(query, o.StatusPesanan, idOrder)

// 	defer con.Close()

// 	return err
// }
