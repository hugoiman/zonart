package models

import "zonart/db"

// JenisPemesananProduk is func
type JenisPemesananProduk struct {
	IDProduk         int    `json:"idProduk"`
	IDJenisPemesanan int    `json:"idJenisPemesanan" validate:"required"`
	Jenis            string `json:"jenis"`
	Harga            int    `json:"harga" validate:"required"`
	Status           bool   `json:"status"`
}

// GetJenisPemesananProduk is func
func (jpp JenisPemesananProduk) GetJenisPemesananProduk(idProduk string) []JenisPemesananProduk {
	con := db.Connect()
	var jenisPemesanan []JenisPemesananProduk

	query := "SELECT a.idProduk, a.idJenisPemesanan, b.jenis, a.harga, a.status FROM jenisPemesananProduk a JOIN jenisPemesanan b ON a.idJenisPemesanan = b.idJenisPemesanan WHERE a.idProduk = ?"

	rows, _ := con.Query(query, idProduk)
	for rows.Next() {
		rows.Scan(
			&jpp.IDProduk, &jpp.IDJenisPemesanan, &jpp.Jenis, &jpp.Harga, &jpp.Status,
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
	_, err := con.Exec(query, idProduk, jpp.IDJenisPemesanan, jpp.Harga, jpp.Status)

	defer con.Close()

	return err
}

// UpdateJenisPemesanan is func
func (jpp JenisPemesananProduk) UpdateJenisPemesanan(idProduk, idJenisPemesanan string) error {
	con := db.Connect()
	query := "UPDATE jenisPemesananProduk SET harga = ?, status = ? WHERE idProduk = ? AND idJenisPemesanan = ?"
	_, err := con.Exec(query, jpp.Harga, jpp.Status, idProduk, idJenisPemesanan)

	defer con.Close()

	return err
}
