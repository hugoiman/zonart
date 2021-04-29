package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Pengiriman is class
type Pengiriman struct {
	idPengiriman int
	penerima     string
	telp         string
	alamat       string
	kota         string
	label        string
	berat        int
	kurir        string
	kodeKurir    string
	service      string
	estimasi     string
	resi         string
	ongkir       int
}

func (p *Pengiriman) SetAlamat(data string) {
	p.alamat = data
}

func (p *Pengiriman) SetKota(data string) {
	p.kota = data
}

func (p *Pengiriman) GetKota() string {
	return p.kota
}

func (p *Pengiriman) SetLabel(data string) {
	p.label = data
}

func (p *Pengiriman) SetBerat(data int) {
	p.berat = data
}

func (p *Pengiriman) GetBerat() int {
	return p.berat
}

func (p *Pengiriman) SetKurir(data string) {
	p.kurir = data
}

func (p *Pengiriman) GetKodeKurir() string {
	return p.kodeKurir
}

func (p *Pengiriman) SetService(data string) {
	p.service = data
}

func (p *Pengiriman) GetService() string {
	return p.service
}

func (p *Pengiriman) SetEstimasi(data string) {
	p.estimasi = data
}

func (p *Pengiriman) SetResi(data string) {
	p.resi = data
}

func (p *Pengiriman) SetOngkir(data int) {
	p.ongkir = data
}

func (p *Pengiriman) GetOngkir() int {
	return p.ongkir
}

func (p *Pengiriman) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDPengiriman int    `json:"idPengiriman"`
		Penerima     string `json:"penerima"`
		Telp         string `json:"telp"`
		Alamat       string `json:"alamat"`
		Kota         string `json:"kota"`
		Label        string `json:"label"`
		Berat        int    `json:"berat"`
		Kurir        string `json:"kurir"`
		KodeKurir    string `json:"kodeKurir"`
		Service      string `json:"service"`
		Estimasi     string `json:"estimasi"`
		Resi         string `json:"resi"`
		Ongkir       int    `json:"ongkir"`
	}{
		IDPengiriman: p.idPengiriman,
		Penerima:     p.penerima,
		Telp:         p.telp,
		Alamat:       p.alamat,
		Kota:         p.kota,
		Label:        p.label,
		Berat:        p.berat,
		Kurir:        p.kurir,
		KodeKurir:    p.kodeKurir,
		Service:      p.service,
		Estimasi:     p.estimasi,
		Resi:         p.resi,
		Ongkir:       p.ongkir,
	})
}

func (p *Pengiriman) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDPengiriman int    `json:"idPengiriman"`
		Penerima     string `json:"penerima" validate:"required"`
		Telp         string `json:"telp" validate:"required"`
		Alamat       string `json:"alamat"`
		Kota         string `json:"kota"`
		Label        string `json:"label"`
		Berat        int    `json:"berat"`
		Kurir        string `json:"kurir"`
		KodeKurir    string `json:"kodeKurir"`
		Service      string `json:"service"`
		Estimasi     string `json:"estimasi"`
		Resi         string `json:"resi"`
		Ongkir       int    `json:"ongkir"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.idPengiriman = alias.IDPengiriman
	p.penerima = alias.Penerima
	p.telp = alias.Telp
	p.alamat = alias.Alamat
	p.kota = alias.Kota
	p.label = alias.Label
	p.berat = alias.Berat
	p.kurir = alias.Kurir
	p.kodeKurir = alias.KodeKurir
	p.service = alias.Service
	p.estimasi = alias.Estimasi
	p.resi = alias.Resi
	p.ongkir = alias.Ongkir

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// CreatePengiriman is func
func (p Pengiriman) CreatePengiriman(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO pengiriman (idOrder, penerima, telp, alamat, kota, label, berat, kurir, kodeKurir, service, estimasi, ongkir) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, p.penerima, p.telp, p.alamat, p.kota, p.label, p.berat, p.kurir, p.kodeKurir, p.service, p.estimasi, p.ongkir)

	defer con.Close()

	return err
}

// GetPengiriman is func
func (p Pengiriman) GetPengiriman(idOrder string) (Pengiriman, error) {
	con := db.Connect()
	query := "SELECT idPengiriman, penerima, telp, alamat, kota, label, berat, kurir, kodeKurir, service, estimasi, resi, ongkir FROM pengiriman WHERE idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&p.idPengiriman, &p.penerima, &p.telp, &p.alamat, &p.kota, &p.label, &p.berat, &p.kurir, &p.kodeKurir, &p.service, &p.estimasi, &p.resi, &p.ongkir)

	defer con.Close()
	return p, err
}

// SetResi is func
func (p Pengiriman) InputResi(idOrder string) error {
	con := db.Connect()
	query := "UPDATE pengiriman SET resi = ? WHERE idOrder = ?"
	_, err := con.Exec(query, p.resi, idOrder)

	defer con.Close()

	return err
}

// UpdateBeratOngkir is func
func (p Pengiriman) UpdateBeratOngkir(idOrder string, berat int, ongkir int) error {
	con := db.Connect()
	query := "UPDATE pengiriman SET berat = ?, ongkir = ? WHERE idOrder = ?"
	_, err := con.Exec(query, berat, ongkir, idOrder)

	defer con.Close()

	return err
}
