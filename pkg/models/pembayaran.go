package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Pembayaran is class
type Pembayaran struct {
	idPembayaran int
	bukti        string
	nominal      int
	status       string
	createdAt    string
}

func (p *Pembayaran) SetBukti(data string) {
	p.bukti = data
}

func (p *Pembayaran) GetBukti() string {
	return p.bukti
}

func (p *Pembayaran) GetNominal() int {
	return p.nominal
}

func (p *Pembayaran) SetStatus(data string) {
	p.status = data
}

func (p *Pembayaran) GetStatus() string {
	return p.status
}

func (p *Pembayaran) SetCreatedAt(data string) {
	p.createdAt = data
}

func (p *Pembayaran) GetCreatedAt() string {
	return p.createdAt
}

func (p *Pembayaran) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDPembayaran int    `json:"idPembayaran"`
		Bukti        string `json:"bukti"`
		Nominal      int    `json:"nominal"`
		Status       string `json:"status"`
		CreatedAt    string `json:"createdAt"`
	}{
		IDPembayaran: p.idPembayaran,
		Bukti:        p.bukti,
		Nominal:      p.nominal,
		Status:       p.status,
		CreatedAt:    p.createdAt,
	})
}

func (p *Pembayaran) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDPembayaran int    `json:"idPembayaran"`
		Bukti        string `json:"bukti"`
		Nominal      int    `json:"nominal" validate:"required"`
		Status       string `json:"status"`
		CreatedAt    string `json:"createdAt"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.idPembayaran = alias.IDPembayaran
	p.bukti = alias.Bukti
	p.nominal = alias.Nominal
	p.status = alias.Status
	p.createdAt = alias.CreatedAt

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetPembayaran is func
func (p Pembayaran) GetPembayaran(idPembayaran, idOrder string) (Pembayaran, error) {
	con := db.Connect()
	query := "SELECT idPembayaran, bukti, nominal, status, createdAt FROM pembayaran WHERE idPembayaran = ? AND idOrder = ?"
	var createdAt time.Time

	err := con.QueryRow(query, idPembayaran, idOrder).Scan(
		&p.idPembayaran, &p.bukti, &p.nominal, &p.status, &createdAt)

	defer con.Close()
	return p, err
}

// GetPembayarans is func
func (p Pembayaran) GetPembayarans(idOrder string) []Pembayaran {
	con := db.Connect()
	var ps []Pembayaran

	query := "SELECT idPembayaran, bukti, nominal, status, createdAt FROM pembayaran WHERE idOrder = ?"

	var createdAt time.Time
	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&p.idPembayaran, &p.bukti, &p.nominal, &p.status, &createdAt,
		)

		p.createdAt = createdAt.Format("02 Jan 2006")
		ps = append(ps, p)
	}

	defer con.Close()

	return ps
}

// CreatePembayaran is func
func (p Pembayaran) CreatePembayaran(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO pembayaran (idOrder, bukti, nominal, status, createdAt) VALUES (?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, p.bukti, p.nominal, p.status, p.createdAt)

	defer con.Close()

	return err
}

// UpdatePembayaran is func
func (p Pembayaran) UpdatePembayaran(idPembayaran, idOrder string) error {
	con := db.Connect()
	query := "UPDATE pembayaran SET status = ? WHERE idPembayaran = ? AND idOrder = ?"
	_, err := con.Exec(query, p.status, idPembayaran, idOrder)

	defer con.Close()

	return err
}
