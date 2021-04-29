package models

import (
	"encoding/json"
	"zonart/db"
)

// Invoice is class
type Invoice struct {
	idInvoice        string
	pembeli          string
	namaToko         string
	slugToko         string
	totalPembelian   int
	totalBayar       int
	tagihan          int
	statusPesanan    string
	statusPembayaran string
}

func (i *Invoice) SetIDInvoice(data string) {
	i.idInvoice = data
}

func (i *Invoice) GetIDInvoice() string {
	return i.idInvoice
}

func (i *Invoice) SetTotalPembelian(data int) {
	i.totalPembelian = data
}

func (i *Invoice) GetTotalPembelian() int {
	return i.totalPembelian
}

func (i *Invoice) GetNamaToko() string {
	return i.namaToko
}

func (i *Invoice) SetTotalBayar(data int) {
	i.totalBayar = data
}

func (i *Invoice) GetTotalBayar() int {
	return i.totalBayar
}

func (i *Invoice) SetTagihan(data int) {
	i.tagihan = data
}

func (i *Invoice) GetTagihan() int {
	return i.tagihan
}

func (i *Invoice) SetStatusPembayaran(data string) {
	i.statusPembayaran = data
}

func (i *Invoice) GetStatusPembayaran() string {
	return i.statusPembayaran
}

func (i *Invoice) SetStatusPesanan(data string) {
	i.statusPesanan = data
}

func (i *Invoice) GetStatusPesanan() string {
	return i.statusPesanan
}

func (i *Invoice) GetSlugToko() string {
	return i.slugToko
}

func (i *Invoice) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDInvoice        string `json:"idInvoice"`
		Pembeli          string `json:"pembeli"`
		NamaToko         string `json:"namaToko"`
		SlugToko         string `json:"slugToko"`
		TotalPembelian   int    `json:"totalPembelian"`
		TotalBayar       int    `json:"totalBayar"`
		Tagihan          int    `json:"tagihan"`
		StatusPesanan    string `json:"statusPesanan"`
		StatusPembayaran string `json:"statusPembayaran"`
	}{
		IDInvoice:        i.idInvoice,
		Pembeli:          i.pembeli,
		NamaToko:         i.namaToko,
		SlugToko:         i.slugToko,
		TotalPembelian:   i.totalPembelian,
		TotalBayar:       i.totalBayar,
		Tagihan:          i.tagihan,
		StatusPesanan:    i.statusPesanan,
		StatusPembayaran: i.statusPembayaran,
	})
}

func (i *Invoice) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDInvoice        string `json:"idInvoice"`
		Pembeli          string `json:"pembeli"`
		NamaToko         string `json:"namaToko"`
		SlugToko         string `json:"slugToko"`
		TotalPembelian   int    `json:"totalPembelian"`
		TotalBayar       int    `json:"totalBayar"`
		Tagihan          int    `json:"tagihan"`
		StatusPesanan    string `json:"statusPesanan"`
		StatusPembayaran string `json:"statusPembayaran"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	i.idInvoice = alias.IDInvoice
	i.pembeli = alias.Pembeli
	i.namaToko = alias.NamaToko
	i.slugToko = alias.SlugToko
	i.totalPembelian = alias.TotalPembelian
	i.totalBayar = alias.TotalBayar
	i.tagihan = alias.Tagihan
	i.statusPesanan = alias.StatusPesanan
	i.statusPembayaran = alias.StatusPembayaran

	return nil
}

// GetInvoice is func
func (i Invoice) GetInvoice(idInvoice string) (Invoice, error) {
	con := db.Connect()
	query := "SELECT a.idInvoice, c.nama, b.namaToko, b.slug, a.totalPembelian, a.totalBayar, a.tagihan, a.statusPesanan, a.statusPembayaran " +
		"FROM `invoice` a JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idInvoice = ?"

	err := con.QueryRow(query, idInvoice).Scan(
		&i.idInvoice, &i.pembeli, &i.namaToko, &i.slugToko,
		&i.totalPembelian, &i.totalBayar, &i.tagihan, &i.statusPesanan, &i.statusPembayaran,
	)

	defer con.Close()

	return i, err
}

// GetInvoiceByOrder is func
func (i Invoice) GetInvoiceByOrder(idOrder string) (Invoice, error) {
	con := db.Connect()
	query := "SELECT a.idInvoice, c.nama, b.namaToko, b.slug, a.totalPembelian, a.totalBayar, a.tagihan, a.statusPesanan, a.statusPembayaran " +
		"FROM `invoice` a JOIN toko b ON a.idToko = b.idToko " +
		"JOIN `order` d ON a.idInvoice = d.idInvoice " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE d.idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&i.idInvoice, &i.pembeli, &i.namaToko, &i.slugToko,
		&i.totalPembelian, &i.totalBayar, &i.tagihan, &i.statusPesanan, &i.statusPembayaran,
	)

	defer con.Close()

	return i, err
}

// CreateInvoice is func
func (i Invoice) CreateInvoice(idToko, idCustomer string) error {
	con := db.Connect()
	query := "INSERT INTO `invoice` (idInvoice, idCustomer, idToko, totalPembelian, totalbayar, tagihan, statusPesanan, statusPembayaran) " +
		"VALUES (?,?,?,?,?,?,?,?)"
	_, err := con.Exec(query, i.idInvoice, idCustomer, idToko, i.totalPembelian, i.totalBayar, i.tagihan, i.statusPesanan, i.statusPembayaran)

	if err != nil {
		return err
	}

	defer con.Close()

	return err
}

// DeleteInvoice is func
func (i Invoice) DeleteInvoice(idInvoice string) error {
	con := db.Connect()
	query := "DELETE FROM `invoice` WHERE idInvoice = ?"
	_, err := con.Exec(query, idInvoice)

	defer con.Close()

	return err
}

// UpdateInvoice is func
func (i Invoice) UpdateInvoice(idInvoice string) error {
	con := db.Connect()
	query := "UPDATE invoice SET totalBayar = ?, tagihan = ?, totalPembelian = ?, statusPembayaran = ?, statusPesanan = ? WHERE idInvoice = ?"
	_, err := con.Exec(query, i.totalBayar, i.tagihan, i.totalPembelian, i.statusPembayaran, i.statusPesanan, idInvoice)

	defer con.Close()

	return err
}
