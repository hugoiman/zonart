package models

import "zonart/db"

// JasaPengiriman is class
type JasaPengiriman struct {
	IDJasaPengiriman int    `json:"idJasaPengiriman"`
	Kurir            string `json:"kurir"`
	Kode             string `json:"kode"`
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
