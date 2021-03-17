package models

import "zonart/db"

// JasaPengiriman is class
type JasaPengiriman struct {
	IDJasaPengiriman int    `json:"idJasaPengiriman"`
	Kurir            string `json:"kurir"`
	Kode             string `json:"kode"`
}

// JasaPengirimans is list of jasaPengiriman
type JasaPengirimans struct {
	JasaPengirimans []JasaPengiriman `json:"jasaPengiriman"`
}

// GetJasaPengirimans is func
func (jp JasaPengiriman) GetJasaPengirimans() JasaPengirimans {
	con := db.Connect()
	query := "SELECT idJasaPengiriman, kurir, kode FROM jasaPengiriman"
	rows, _ := con.Query(query)

	var jasaPengirimans JasaPengirimans

	for rows.Next() {
		rows.Scan(
			&jp.IDJasaPengiriman, &jp.Kurir, &jp.Kode,
		)

		jasaPengirimans.JasaPengirimans = append(jasaPengirimans.JasaPengirimans, jp)
	}

	defer con.Close()

	return jasaPengirimans
}

// GetJasaPengiriman is func
func (jp JasaPengiriman) GetJasaPengiriman() []JasaPengiriman {
	con := db.Connect()
	var jps []JasaPengiriman
	query := "SELECT idJasaPengiriman, kurir, kode FROM jasaPengiriman"
	rows, _ := con.Query(query)
	for rows.Next() {
		rows.Scan(
			&jp.IDJasaPengiriman, &jp.Kurir, &jp.Kode,
		)

		jps = append(jps, jp)
	}

	defer con.Close()

	return jps
}
