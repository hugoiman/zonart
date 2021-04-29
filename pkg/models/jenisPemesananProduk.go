package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// JenisPemesananProduk is func
type JenisPemesananProduk struct {
	idJenisPemesanan int
	jenis            string
	harga            int
	status           bool
}

func (jpp *JenisPemesananProduk) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDJenisPemesanan int    `json:"idJenisPemesanan"`
		Jenis            string `json:"jenis"`
		Harga            int    `json:"harga"`
		Status           bool   `json:"status"`
	}{
		IDJenisPemesanan: jpp.idJenisPemesanan,
		Jenis:            jpp.jenis,
		Harga:            jpp.harga,
		Status:           jpp.status,
	})
}

func (jpp *JenisPemesananProduk) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDJenisPemesanan int    `json:"idJenisPemesanan" validate:"required,eq=1|eq=2"`
		Jenis            string `json:"jenis"`
		Harga            int    `json:"harga"`
		Status           bool   `json:"status"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	jpp.idJenisPemesanan = alias.IDJenisPemesanan
	jpp.jenis = alias.Jenis
	jpp.harga = alias.Harga
	jpp.status = alias.Status

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetJenisPemesananProduk is func
func (jpp JenisPemesananProduk) GetJenisPemesananProduk(idProduk string) []JenisPemesananProduk {
	con := db.Connect()
	var jenisPemesanan []JenisPemesananProduk

	query := "SELECT a.idJenisPemesanan, b.jenis, a.harga, a.status FROM jenisPemesananProduk a JOIN jenisPemesanan b ON a.idJenisPemesanan = b.idJenisPemesanan WHERE a.idProduk = ?"

	rows, _ := con.Query(query, idProduk)
	for rows.Next() {
		rows.Scan(
			&jpp.idJenisPemesanan, &jpp.jenis, &jpp.harga, &jpp.status,
		)

		jenisPemesanan = append(jenisPemesanan, jpp)
	}

	defer con.Close()

	return jenisPemesanan
}

// CreateJenisPemesanan is func
func (jpp JenisPemesananProduk) CreateJenisPemesanan(idProduk string) error {
	con := db.Connect()
	query := "INSERT INTO jenisPemesananProduk (idProduk, idJenisPemesanan, harga, status) VALUES (?,?,?,?)"
	_, err := con.Exec(query, idProduk, jpp.idJenisPemesanan, jpp.harga, jpp.status)

	defer con.Close()

	return err
}

// UpdateJenisPemesanan is func
func (jpp JenisPemesananProduk) UpdateJenisPemesanan(idProduk, idJenisPemesanan string) error {
	con := db.Connect()
	query := "UPDATE jenisPemesananProduk SET harga = ?, status = ? WHERE idProduk = ? AND idJenisPemesanan = ?"
	_, err := con.Exec(query, jpp.harga, jpp.status, idProduk, idJenisPemesanan)

	defer con.Close()

	return err
}
