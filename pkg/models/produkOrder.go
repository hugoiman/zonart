package models

import (
	"encoding/json"
	"zonart/db"
)

// ProdukOrder is class
type ProdukOrder struct {
	idProdukOrder    int
	namaProduk       string
	beratProduk      int
	hargaProduk      int
	hargaSatuanWajah int
	fotoProduk       string
	slugProduk       string
}

func (po *ProdukOrder) SetNamaProduk(data string) {
	po.namaProduk = data
}

func (po *ProdukOrder) GetNamaProduk() string {
	return po.namaProduk
}

func (po *ProdukOrder) SetBeratProduk(data int) {
	po.beratProduk = data
}

func (po *ProdukOrder) SetHargaProduk(data int) {
	po.hargaProduk = data
}

func (po *ProdukOrder) GetHargaProduk() int {
	return po.hargaProduk
}

func (po *ProdukOrder) SetHargaSatuanWajah(data int) {
	po.hargaSatuanWajah = data
}

func (po *ProdukOrder) GetHargaSatuanWajah() int {
	return po.hargaSatuanWajah
}

func (p *ProdukOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDProdukOrder    int    `json:"idProdukOrder"`
		NamaProduk       string `json:"namaProduk"`
		BeratProduk      int    `json:"beratProduk"`
		HargaProduk      int    `json:"hargaProduk"`
		HargaSatuanWajah int    `json:"hargaSatuanWajah"`
		FotoProduk       string `json:"fotoProduk"`
		SlugProduk       string `json:"slugProduk"`
	}{
		IDProdukOrder:    p.idProdukOrder,
		NamaProduk:       p.namaProduk,
		BeratProduk:      p.beratProduk,
		HargaProduk:      p.hargaProduk,
		HargaSatuanWajah: p.hargaSatuanWajah,
		FotoProduk:       p.fotoProduk,
		SlugProduk:       p.slugProduk,
	})
}

func (p *ProdukOrder) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDProdukOrder    int    `json:"idProdukOrder"`
		NamaProduk       string `json:"namaProduk"`
		BeratProduk      int    `json:"beratProduk"`
		HargaProduk      int    `json:"hargaProduk"`
		HargaSatuanWajah int    `json:"hargaSatuanWajah"`
		FotoProduk       string `json:"fotoProduk"`
		SlugProduk       string `json:"slugProduk"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	p.idProdukOrder = alias.IDProdukOrder
	p.namaProduk = alias.NamaProduk
	p.beratProduk = alias.BeratProduk
	p.hargaProduk = alias.HargaProduk
	p.hargaSatuanWajah = alias.HargaSatuanWajah
	p.fotoProduk = alias.FotoProduk
	p.slugProduk = alias.SlugProduk

	return nil
}

// GetProdukOrder is func
func (po ProdukOrder) GetProdukOrder(idOrder string) (ProdukOrder, error) {
	con := db.Connect()
	query := "SELECT a.idProdukOrder, a.namaProduk, a.beratProduk, a.hargaProduk, a.hargaSatuanWajah, c.gambar, c.slug " +
		"FROM produkOrder a JOIN `order` b ON a.idOrder = b.idOrder " +
		"JOIN produk c ON b.idProduk = c.idProduk WHERE a.idOrder = ?"

	err := con.QueryRow(query, idOrder).Scan(
		&po.idProdukOrder, &po.namaProduk, &po.beratProduk, &po.hargaProduk, &po.hargaSatuanWajah, &po.fotoProduk, &po.slugProduk,
	)

	defer con.Close()

	return po, err
}

// CreateProdukOrder is func
func (po ProdukOrder) CreateProdukOrder(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO produkOrder (idOrder, namaProduk, beratProduk, hargaProduk, hargaSatuanWajah) " +
		"VALUES (?,?,?,?,?)"
	_, err := con.Exec(query, idOrder, po.namaProduk, po.beratProduk, po.hargaProduk, po.hargaSatuanWajah)

	defer con.Close()

	return err
}
