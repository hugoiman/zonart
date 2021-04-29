package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// OpsiOrder is class
type OpsiOrder struct {
	idOpsiOrder int
	idGrupOpsi  int
	namaGrup    string
	opsi        string
	harga       int
	berat       int
	perProduk   bool
}

func (oo *OpsiOrder) GetIDGrupOpsi() int {
	return oo.idGrupOpsi
}

func (oo *OpsiOrder) GetIDOpsiOrder() int {
	return oo.idOpsiOrder
}

func (oo *OpsiOrder) SetNamaGrup(data string) {
	oo.namaGrup = data
}

func (oo *OpsiOrder) GetNamaGrup() string {
	return oo.namaGrup
}

func (oo *OpsiOrder) SetOpsi(data string) {
	oo.opsi = data
}

func (oo *OpsiOrder) GetOpsi() string {
	return oo.opsi
}

func (oo *OpsiOrder) SetHarga(data int) {
	oo.harga = data
}

func (oo *OpsiOrder) GetHarga() int {
	return oo.harga
}

func (oo *OpsiOrder) SetBerat(data int) {
	oo.berat = data
}

func (oo *OpsiOrder) GetBerat() int {
	return oo.berat
}

func (oo *OpsiOrder) SetPerProduk(data bool) {
	oo.perProduk = data
}

func (oo *OpsiOrder) GetPerProduk() bool {
	return oo.perProduk
}

func (oo *OpsiOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDOpsiOrder int    `json:"idOpsiOrder"`
		IDGrupOpsi  int    `json:"idGrupOpsi"`
		NamaGrup    string `json:"namaGrup"`
		Opsi        string `json:"opsi"`
		Harga       int    `json:"harga"`
		Berat       int    `json:"berat"`
		PerProduk   bool   `json:"perProduk"`
	}{
		IDOpsiOrder: oo.idOpsiOrder,
		IDGrupOpsi:  oo.idGrupOpsi,
		NamaGrup:    oo.namaGrup,
		Opsi:        oo.opsi,
		Harga:       oo.harga,
		Berat:       oo.berat,
		PerProduk:   oo.perProduk,
	})
}

func (oo *OpsiOrder) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDOpsiOrder int    `json:"idOpsiOrder"`
		IDGrupOpsi  int    `json:"idGrupOpsi"  validate:"required"`
		NamaGrup    string `json:"namaGrup"`
		Opsi        string `json:"opsi"`
		Harga       int    `json:"harga"`
		Berat       int    `json:"berat"`
		PerProduk   bool   `json:"perProduk"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	oo.idOpsiOrder = alias.IDOpsiOrder
	oo.idGrupOpsi = alias.IDGrupOpsi
	oo.namaGrup = alias.NamaGrup
	oo.opsi = alias.Opsi
	oo.harga = alias.Harga
	oo.berat = alias.Berat
	oo.perProduk = alias.PerProduk

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// CreateOpsiOrder is func
func (oo OpsiOrder) CreateOpsiOrder(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO opsiOrder (idOrder, namaGrup, opsi, harga, berat, perProduk) VALUES (?,?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, oo.namaGrup, oo.opsi, oo.harga, oo.berat, oo.perProduk)

	defer con.Close()

	return err
}

// GetOpsiOrder is func
func (oo OpsiOrder) GetOpsiOrder(idOrder string) []OpsiOrder {
	con := db.Connect()
	var opsiOrders []OpsiOrder

	query := "SELECT idOpsiOrder, namaGrup, opsi, harga, berat, perProduk FROM opsiOrder WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&oo.idOpsiOrder, &oo.namaGrup, &oo.opsi, &oo.harga, &oo.berat, &oo.perProduk,
		)

		opsiOrders = append(opsiOrders, oo)
	}

	defer con.Close()

	return opsiOrders
}
