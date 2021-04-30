package models

import (
	"encoding/json"
	"strconv"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Order is class
type Order struct {
	idOrder         int
	pemesan         int
	jenisPesanan    string
	tambahanWajah   int
	catatan         string
	pcs             int
	rencanaPakai    string
	waktuPengerjaan string
	contohGambar    string
	tglOrder        string
	invoice         Invoice
	produkOrder     ProdukOrder
	pengiriman      Pengiriman
	penangan        Penangan
	hasilOrder      HasilOrder
	fileOrder       []FileOrder
	biayaTambahan   []BiayaTambahan
	opsiOrder       []OpsiOrder
	revisi          []Revisi
	pembayaran      []Pembayaran
}

func (o *Order) GetIDOrder() int {
	return o.idOrder
}

func (o *Order) SetPemesan(data int) {
	o.pemesan = data
}

func (o *Order) GetPemesan() int {
	return o.pemesan
}

func (o *Order) GetJenisPesanan() string {
	return o.jenisPesanan
}

func (o *Order) GetPcs() int {
	return o.pcs
}

func (o *Order) GetTambahanWajah() int {
	return o.tambahanWajah
}

func (o *Order) SetWaktu(data string) {
	o.waktuPengerjaan = data
}

func (o *Order) GetInvoice() *Invoice {
	return &o.invoice
}

func (o *Order) GetPengiriman() *Pengiriman {
	return &o.pengiriman
}

func (o *Order) GetPenangan() *Penangan {
	return &o.penangan
}

func (o *Order) GetHasilOrder() *HasilOrder {
	return &o.hasilOrder
}

func (o *Order) SetTglOrder(data string) {
	o.tglOrder = data
}

func (o *Order) GetTglOrder() string {
	return o.tglOrder
}

func (o *Order) GetProdukOrder() *ProdukOrder {
	return &o.produkOrder
}

func (o *Order) SetFileOrder(data []FileOrder) {
	o.fileOrder = data
}

func (o *Order) GetFileOrder() []FileOrder {
	return o.fileOrder
}

func (o *Order) GetOpsiOrder() []OpsiOrder {
	return o.opsiOrder
}

func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDOrder         int             `json:"idOrder"`
		Pemesan         int             `json:"pemesan"`
		JenisPesanan    string          `json:"jenisPesanan"`
		TambahanWajah   int             `json:"tambahanWajah"`
		Catatan         string          `json:"catatan"`
		Pcs             int             `json:"pcs"`
		RencanaPakai    string          `json:"rencanaPakai"`
		WaktuPengerjaan string          `json:"waktuPengerjaan"`
		ContohGambar    string          `json:"contohGambar"`
		TglOrder        string          `json:"tglOrder"`
		Invoice         *Invoice        `json:"invoice"`
		ProdukOrder     *ProdukOrder    `json:"produkOrder"`
		Pengiriman      *Pengiriman     `json:"pengiriman"`
		Penangan        *Penangan       `json:"penangan"`
		HasilOrder      *HasilOrder     `json:"hasilOrder"`
		FileOrder       []FileOrder     `json:"fileOrder"`
		BiayaTambahan   []BiayaTambahan `json:"biayaTambahan"`
		OpsiOrder       []OpsiOrder     `json:"opsiOrder"`
		Revisi          []Revisi        `json:"revisi"`
		Pembayaran      []Pembayaran    `json:"pembayaran"`
	}{
		IDOrder:         o.idOrder,
		Pemesan:         o.pemesan,
		JenisPesanan:    o.jenisPesanan,
		TambahanWajah:   o.tambahanWajah,
		Catatan:         o.catatan,
		Pcs:             o.pcs,
		RencanaPakai:    o.rencanaPakai,
		WaktuPengerjaan: o.waktuPengerjaan,
		ContohGambar:    o.contohGambar,
		TglOrder:        o.tglOrder,
		Invoice:         &o.invoice,
		ProdukOrder:     &o.produkOrder,
		Pengiriman:      &o.pengiriman,
		Penangan:        &o.penangan,
		HasilOrder:      &o.hasilOrder,
		FileOrder:       o.fileOrder,
		BiayaTambahan:   o.biayaTambahan,
		OpsiOrder:       o.opsiOrder,
		Revisi:          o.revisi,
		Pembayaran:      o.pembayaran,
	})
}

func (o *Order) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDOrder         int             `json:"idOrder"`
		Pemesan         int             `json:"pemesan"`
		JenisPesanan    string          `json:"jenisPesanan" validate:"required,eq=cetak|eq=soft copy"`
		TambahanWajah   int             `json:"tambahanWajah"`
		Catatan         string          `json:"catatan"`
		Pcs             int             `json:"pcs" validate:"required,min=1"`
		RencanaPakai    string          `json:"rencanaPakai"`
		WaktuPengerjaan string          `json:"waktuPengerjaan"`
		ContohGambar    string          `json:"contohGambar"`
		TglOrder        string          `json:"tglOrder"`
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
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	o.idOrder = alias.IDOrder
	o.pemesan = alias.Pemesan
	o.jenisPesanan = alias.JenisPesanan
	o.tambahanWajah = alias.TambahanWajah
	o.catatan = alias.Catatan
	o.pcs = alias.Pcs
	o.rencanaPakai = alias.RencanaPakai
	o.waktuPengerjaan = alias.WaktuPengerjaan
	o.contohGambar = alias.ContohGambar
	o.tglOrder = alias.TglOrder
	o.invoice = alias.Invoice
	o.produkOrder = alias.ProdukOrder
	o.pengiriman = alias.Pengiriman
	o.penangan = alias.Penangan
	o.hasilOrder = alias.HasilOrder
	o.fileOrder = alias.FileOrder
	o.biayaTambahan = alias.BiayaTambahan
	o.opsiOrder = alias.OpsiOrder
	o.revisi = alias.Revisi
	o.pembayaran = alias.Pembayaran

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetOrder is func
func (o Order) GetOrder(idOrder string) (Order, error) {
	con := db.Connect()
	query := "SELECT idOrder, idCustomer, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, tglOrder " +
		"FROM `order` WHERE idOrder = ?"

	var tglOrder time.Time

	err := con.QueryRow(query, idOrder).Scan(
		&o.idOrder, &o.pemesan,
		&o.jenisPesanan, &o.tambahanWajah, &o.catatan, &o.pcs, &o.rencanaPakai,
		&o.waktuPengerjaan, &o.contohGambar, &tglOrder,
	)

	o.tglOrder = tglOrder.Format("02 Jan 2006")

	o.hasilOrder, _ = o.hasilOrder.GetHasilOrder(idOrder)
	o.invoice, _ = o.invoice.GetInvoiceByOrder(idOrder)
	o.produkOrder, _ = o.produkOrder.GetProdukOrder(idOrder)
	o.penangan, _ = o.penangan.GetPenangan(idOrder)
	o.pengiriman, _ = o.pengiriman.GetPengiriman(idOrder)

	var opsiOrder OpsiOrder
	o.opsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	var fileOrder FileOrder
	o.fileOrder = fileOrder.GetFileOrder(idOrder)

	var revisi Revisi
	o.revisi = revisi.GetRevisi(idOrder)

	var bt BiayaTambahan
	o.biayaTambahan = bt.GetBiayaTambahans(idOrder)

	var pembayaran Pembayaran
	o.pembayaran = pembayaran.GetPembayarans(idOrder)

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
		&o.idOrder, &o.pemesan,
		&o.jenisPesanan, &o.tambahanWajah, &o.catatan, &o.pcs, &o.rencanaPakai,
		&o.waktuPengerjaan, &o.contohGambar, &tglOrder,
	)

	o.tglOrder = tglOrder.Format("02 Jan 2006")

	o.hasilOrder, _ = o.hasilOrder.GetHasilOrder(idOrder)
	o.invoice, _ = o.invoice.GetInvoiceByOrder(idOrder)
	o.produkOrder, _ = o.produkOrder.GetProdukOrder(idOrder)
	o.penangan, _ = o.penangan.GetPenangan(idOrder)
	o.pengiriman, _ = o.pengiriman.GetPengiriman(idOrder)

	var opsiOrder OpsiOrder
	o.opsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	var fileOrder FileOrder
	o.fileOrder = fileOrder.GetFileOrder(idOrder)

	var revisi Revisi
	o.revisi = revisi.GetRevisi(idOrder)

	var bt BiayaTambahan
	o.biayaTambahan = bt.GetBiayaTambahans(idOrder)

	var pembayaran Pembayaran
	o.pembayaran = pembayaran.GetPembayarans(idOrder)

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
		&o.idOrder, &o.pemesan,
		&o.jenisPesanan, &o.tambahanWajah, &o.catatan, &o.pcs, &o.rencanaPakai,
		&o.waktuPengerjaan, &o.contohGambar, &tglOrder,
	)

	o.tglOrder = tglOrder.Format("02 Jan 2006")

	o.hasilOrder, _ = o.hasilOrder.GetHasilOrder(idOrder)
	o.invoice, _ = o.invoice.GetInvoice(idOrder)
	o.produkOrder, _ = o.produkOrder.GetProdukOrder(idOrder)
	o.penangan, _ = o.penangan.GetPenangan(idOrder)
	o.pengiriman, _ = o.pengiriman.GetPengiriman(idOrder)

	var opsiOrder OpsiOrder
	o.opsiOrder = opsiOrder.GetOpsiOrder(idOrder)

	var fileOrder FileOrder
	o.fileOrder = fileOrder.GetFileOrder(idOrder)

	var revisi Revisi
	o.revisi = revisi.GetRevisi(idOrder)

	var bt BiayaTambahan
	o.biayaTambahan = bt.GetBiayaTambahans(idOrder)

	var pembayaran Pembayaran
	o.pembayaran = pembayaran.GetPembayarans(idOrder)

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
		&o.idOrder, &o.pemesan,
		&o.jenisPesanan, &o.tambahanWajah, &o.catatan, &o.pcs, &o.rencanaPakai,
		&o.waktuPengerjaan, &o.contohGambar, &tglOrder,
	)

	o.tglOrder = tglOrder.Format("02 Jan 2006")

	o.invoice, _ = o.invoice.GetInvoice(idInvoice)
	o.produkOrder, _ = o.produkOrder.GetProdukOrder(strconv.Itoa(o.idOrder))
	o.pengiriman, _ = o.pengiriman.GetPengiriman(strconv.Itoa(o.idOrder))

	var opsiOrder OpsiOrder
	o.opsiOrder = opsiOrder.GetOpsiOrder(strconv.Itoa(o.idOrder))

	var bt BiayaTambahan
	o.biayaTambahan = bt.GetBiayaTambahans(strconv.Itoa(o.idOrder))

	defer con.Close()
	return o, err
}

// GetOrders is func
func (o Order) GetOrders(idCustomer string) []Order {
	con := db.Connect()
	query := "SELECT idOrder, idCustomer, " +
		"jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, waktuPengerjaan, contohGambar, tglOrder " +
		"FROM `order` WHERE idCustomer = ?"
	rows, _ := con.Query(query, idCustomer)

	var orders []Order
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.idOrder, &o.pemesan,
			&o.jenisPesanan, &o.tambahanWajah, &o.catatan, &o.pcs, &o.rencanaPakai,
			&o.waktuPengerjaan, &o.contohGambar, &tglOrder,
		)

		o.tglOrder = tglOrder.Format("02 Jan 2006")

		o.hasilOrder, _ = o.hasilOrder.GetHasilOrder(strconv.Itoa(o.idOrder))
		o.invoice, _ = o.invoice.GetInvoiceByOrder(strconv.Itoa(o.idOrder))
		o.produkOrder, _ = o.produkOrder.GetProdukOrder(strconv.Itoa(o.idOrder))

		var penangan Penangan
		o.penangan, _ = penangan.GetPenangan(strconv.Itoa(o.idOrder))

		var pengiriman Pengiriman
		o.pengiriman, _ = pengiriman.GetPengiriman(strconv.Itoa(o.idOrder))

		var opsiOrder OpsiOrder
		o.opsiOrder = opsiOrder.GetOpsiOrder(strconv.Itoa(o.idOrder))

		var fileOrder FileOrder
		o.fileOrder = fileOrder.GetFileOrder(strconv.Itoa(o.idOrder))

		var revisi Revisi
		o.revisi = revisi.GetRevisi(strconv.Itoa(o.idOrder))

		var bt BiayaTambahan
		o.biayaTambahan = bt.GetBiayaTambahans(strconv.Itoa(o.idOrder))

		var pembayaran Pembayaran
		o.pembayaran = pembayaran.GetPembayarans(strconv.Itoa(o.idOrder))

		orders = append(orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersToko is func
func (o Order) GetOrdersToko(idToko string) []Order {
	con := db.Connect()
	query := "SELECT idOrder, tglOrder FROM `order` WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var orders []Order
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.idOrder, &tglOrder,
		)

		o.tglOrder = tglOrder.Format("02 Jan 2006")
		o.invoice, _ = o.invoice.GetInvoiceByOrder(strconv.Itoa(o.idOrder))
		o.produkOrder, _ = o.produkOrder.GetProdukOrder(strconv.Itoa(o.idOrder))

		orders = append(orders, o)
	}

	defer con.Close()
	return orders
}

// GetOrdersEditor is func
func (o Order) GetOrdersEditor(idToko, idKaryawan string) []Order {
	con := db.Connect()
	query := "SELECT a.idOrder, a.tglOrder " +
		"FROM `order` a JOIN penangan b ON a.idOrder = b.idOrder " +
		"WHERE a.idToko = ? AND b.idKaryawan = ?"
	rows, _ := con.Query(query, idToko, idKaryawan)

	var orders []Order
	var tglOrder time.Time

	for rows.Next() {
		rows.Scan(
			&o.idOrder, &tglOrder,
		)

		o.tglOrder = tglOrder.Format("02 Jan 2006")
		o.invoice, _ = o.invoice.GetInvoiceByOrder(strconv.Itoa(o.idOrder))
		o.produkOrder, _ = o.produkOrder.GetProdukOrder(strconv.Itoa(o.idOrder))

		orders = append(orders, o)
	}

	defer con.Close()
	return orders
}

// CreateOrder is func
func (o Order) CreateOrder(idToko, idProduk string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO `order` (idToko, idProduk, idCustomer, idInvoice, jenisPesanan, tambahanWajah, catatan, pcs, rencanaPakai, contohGambar, tglOrder) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, idProduk, o.pemesan, o.invoice.idInvoice, o.jenisPesanan, o.tambahanWajah, o.catatan, o.pcs,
		o.rencanaPakai, o.contohGambar, o.tglOrder)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idOrder := int(idInt64)

	err = o.pengiriman.CreatePengiriman(strconv.Itoa(idOrder))
	if err != nil {
		return 0, err
	}

	err = o.produkOrder.CreateProdukOrder(strconv.Itoa(idOrder))
	if err != nil {
		return 0, err
	}

	for _, v := range o.fileOrder {
		err = v.CreateFileOrder(strconv.Itoa(idOrder))
		if err != nil {
			return 0, err
		}
	}

	for _, v := range o.opsiOrder {
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
	_, err := con.Exec(query, o.waktuPengerjaan, idOrder)

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
