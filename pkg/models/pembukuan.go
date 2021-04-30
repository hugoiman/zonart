package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Pembukuan is class
type Pembukuan struct {
	idPembukuan  int
	jenis        string
	keterangan   string
	nominal      int
	tglTransaksi string
}

func (p *Pembukuan) SetJenis(data string) {
	p.jenis = data
}

func (p *Pembukuan) GetJenis() string {
	return p.jenis
}

func (p *Pembukuan) SetKeterangan(data string) {
	p.keterangan = data
}

func (p *Pembukuan) GetKeterangan() string {
	return p.keterangan
}

func (p *Pembukuan) SetNominal(data int) {
	p.nominal = data
}

func (p *Pembukuan) GetNominal() int {
	return p.nominal
}

func (p *Pembukuan) SetTglTransaksi(data string) {
	p.tglTransaksi = data
}

func (p *Pembukuan) GetTglTransaksi() string {
	return p.tglTransaksi
}

func (p *Pembukuan) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDPembukuan  int    `json:"idPembukuan"`
		Jenis        string `json:"jenis"`
		Keterangan   string `json:"keterangan"`
		Nominal      int    `json:"nominal"`
		TglTransaksi string `json:"tglTransaksi"`
	}{
		IDPembukuan:  p.idPembukuan,
		Jenis:        p.jenis,
		Keterangan:   p.keterangan,
		Nominal:      p.nominal,
		TglTransaksi: p.tglTransaksi,
	})
}

func (p *Pembukuan) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDPembukuan  int    `json:"idPembukuan"`
		Jenis        string `json:"jenis" validate:"required"`
		Keterangan   string `json:"keterangan" validate:"required"`
		Nominal      int    `json:"nominal" validate:"required"`
		TglTransaksi string `json:"tglTransaksi" validate:"required"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.idPembukuan = alias.IDPembukuan
	p.jenis = alias.Jenis
	p.keterangan = alias.Keterangan
	p.nominal = alias.Nominal
	p.tglTransaksi = alias.TglTransaksi

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetPembukuans is func
func (p Pembukuan) GetPembukuans(idToko string) []Pembukuan {
	con := db.Connect()
	query := "SELECT idPembukuan, jenis, keterangan, nominal, tglTransaksi FROM pembukuan WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var tglTransaksi time.Time
	var pembukuans []Pembukuan

	for rows.Next() {
		rows.Scan(
			&p.idPembukuan, &p.jenis, &p.keterangan, &p.nominal, &tglTransaksi,
		)
		p.tglTransaksi = tglTransaksi.Format("02 Jan 2006")
		pembukuans = append(pembukuans, p)
	}

	defer con.Close()

	return pembukuans
}

// CreatePembukuan is func
func (p Pembukuan) CreatePembukuan(idToko string) error {
	con := db.Connect()
	query := "INSERT INTO pembukuan (idToko, jenis, keterangan, nominal, tglTransaksi) VALUES (?,?,?,?,?)"
	_, err := con.Exec(query, idToko, p.jenis, p.keterangan, p.nominal, p.tglTransaksi)

	defer con.Close()

	return err
}

// DeletePembukuan is func
func (p Pembukuan) DeletePembukuan(idToko, idPembukuan string) error {
	con := db.Connect()
	query := "DELETE FROM pembukuan WHERE idToko = ? AND idPembukuan = ?"
	_, err := con.Exec(query, idToko, idPembukuan)

	defer con.Close()

	return err
}
