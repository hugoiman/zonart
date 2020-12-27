package models

import (
	"time"
	"zonart/db"
)

// Pembayaran is class
type Pembayaran struct {
	IDPembayaran int    `json:"idPembayaran"`
	IDOrder      int    `json:"idOrder"`
	Bukti        string `json:"bukti"`
	Nominal      int    `json:"nominal"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
}

// GetPembayaran is func
func (p Pembayaran) GetPembayaran(idPembayaran, idOrder string) (Pembayaran, error) {
	con := db.Connect()
	query := "SELECT idPembayaran, idOrder, bukti, nominal, status, createdAt FROM pembayaran WHERE idPembayaran = ? AND idOrder = ?"
	var createdAt time.Time

	err := con.QueryRow(query, idPembayaran, idOrder).Scan(
		&p.IDPembayaran, &p.IDOrder, &p.Bukti, &p.Nominal, &p.Nominal, &createdAt)

	defer con.Close()
	return p, err
}

// GetPembayarans is func
func (p Pembayaran) GetPembayarans(idOrder string) []Pembayaran {
	con := db.Connect()
	var ps []Pembayaran

	query := "SELECT idPembayaran, idOrder, bukti, nominal, status, createdAt WHERE idOrder = ?"

	var createdAt time.Time
	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&p.IDPembayaran, &p.IDOrder, &p.Bukti, &p.Nominal, &p.Nominal, &createdAt,
		)

		p.CreatedAt = createdAt.Format("02 Jan 2006")
		ps = append(ps, p)
	}

	defer con.Close()

	return ps
}

// CreatePembayaran is func
func (p Pembayaran) CreatePembayaran(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO tambahanBiaya (idOrder, bukti, nominal, status, createdAt) VALUES (?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, p.Bukti, p.Nominal, p.Status, p.CreatedAt)

	defer con.Close()

	return err
}

// KonfirmasiPembayaran is func
func (p Pembayaran) KonfirmasiPembayaran(idPembayaran, idOrder string) error {
	con := db.Connect()
	query := "UPDATE pembayaran a JOIN order b ON a.idOrder = b.idOrder SET" +
		" a.status = ?, b.dibayar = b.dibayar + ?, b.tagihan = b.tagihan - ? WHERE a.idPembayaran = ? AND a.idOrder = ?"
	_, err := con.Exec(query, p.Status, p.Nominal, p.Nominal, idPembayaran, idOrder)

	defer con.Close()

	return err
}
