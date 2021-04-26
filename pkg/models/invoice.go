package models

import (
	"zonart/db"
)

// Invoice is class
type Invoice struct {
	IDInvoice        string `json:"idInvoice"`
	Pembeli          string `json:"pembeli"`
	NamaToko         string `json:"namaToko"`
	SlugToko         string `json:"slugToko"`
	TotalPembelian   int    `json:"totalPembelian"`
	TotalBayar       int    `json:"totalBayar"`
	Tagihan          int    `json:"tagihan"`
	StatusPesanan    string `json:"statusPesanan"`
	StatusPembayaran string `json:"statusPembayaran"`
}

// GetInvoice is func
func (i Invoice) GetInvoice(idInvoice string) (Invoice, error) {
	con := db.Connect()
	query := "SELECT a.idInvoice, c.nama, b.namaToko, b.slug, a.totalPembelian, a.totalBayar, a.tagihan, a.statusPesanan, a.statusPembayaran " +
		"FROM `invoice` a JOIN toko b ON a.idToko = b.idToko " +
		"JOIN customer c ON a.idCustomer = c.idCustomer WHERE a.idInvoice = ?"

	err := con.QueryRow(query, idInvoice).Scan(
		&i.IDInvoice, &i.Pembeli, &i.NamaToko, &i.SlugToko,
		&i.TotalPembelian, &i.TotalBayar, &i.Tagihan, &i.StatusPesanan, &i.StatusPembayaran,
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
		&i.IDInvoice, &i.Pembeli, &i.NamaToko, &i.SlugToko,
		&i.TotalPembelian, &i.TotalBayar, &i.Tagihan, &i.StatusPesanan, &i.StatusPembayaran,
	)

	defer con.Close()

	return i, err
}

// CreateInvoice is func
func (i Invoice) CreateInvoice(idToko, idCustomer string) error {
	con := db.Connect()
	query := "INSERT INTO `invoice` (idInvoice, idCustomer, idToko, totalPembelian, totalbayar, tagihan, statusPesanan, statusPembayaran) " +
		"VALUES (?,?,?,?,?,?,?,?)"
	_, err := con.Exec(query, i.IDInvoice, idCustomer, idToko, i.TotalPembelian, i.TotalBayar, i.Tagihan, i.StatusPesanan, i.StatusPembayaran)

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
	_, err := con.Exec(query, i.TotalBayar, i.Tagihan, i.TotalPembelian, i.StatusPembayaran, i.StatusPesanan, idInvoice)

	defer con.Close()

	return err
}
