package models

import (
	"zonart/db"
)

// JasaPengirimanToko is class
type JasaPengirimanToko struct {
	IDJasaPengiriman int    `json:"idJasaPengiriman"`
	IDToko           int    `json:"idToko"`
	Kurir            string `json:"kurir"`
	Kode             string `json:"kode"`
	Status           bool   `json:"status"`
}

// GetJasaPengirimanToko is func
func (jpt JasaPengirimanToko) GetJasaPengirimanToko(idToko int) []JasaPengirimanToko {
	con := db.Connect()
	var jpts []JasaPengirimanToko

	var jp JasaPengiriman
	jasaPengiriman := jp.GetJasaPengiriman()

	query := "SELECT a.idJasaPengiriman, a.idToko, b.kurir, b.kode, a.status FROM jasaPengirimanToko a JOIN jasaPengiriman b ON a.idJasaPengiriman = b.idJasaPengiriman WHERE a.idToko = ? AND a.idJasaPengiriman = ?"

	for _, v := range jasaPengiriman {
		err := con.QueryRow(query, idToko, v.IDJasaPengiriman).Scan(&jpt.IDJasaPengiriman, &jpt.IDToko, &jpt.Kurir, &jpt.Kode, &jpt.Status)
		if err != nil {
			jpt.IDJasaPengiriman = v.IDJasaPengiriman
			jpt.IDToko = idToko
			jpt.Kurir = v.Kurir
			jpt.Kode = v.Kode
			jpt.Status = false
		}
		jpts = append(jpts, jpt)
	}

	defer con.Close()

	return jpts
}

// CreateUpdatePengirimanToko is func
func (jpt JasaPengirimanToko) CreateUpdatePengirimanToko(idToko string) error {
	var isAny bool
	var err error
	con := db.Connect()

	query := "SELECT EXISTS (SELECT 1 FROM jasaPengirimanToko WHERE idToko = ? AND idJasaPengiriman = ?)"

	con.QueryRow(query, idToko, jpt.IDJasaPengiriman).Scan(&isAny)

	if isAny == true {
		queryUpdate := "UPDATE jasaPengirimanToko SET status = ? WHERE idJasaPengiriman = ? AND idToko = ?"
		_, err = con.Exec(queryUpdate, jpt.Status, jpt.IDJasaPengiriman, idToko)
	} else if isAny == false {
		queryInsert := "INSERT INTO jasaPengirimanToko (idJasaPengiriman, idToko, status) VALUES (?,?,?)"
		_, err = con.Exec(queryInsert, jpt.IDJasaPengiriman, idToko, jpt.Status)
	}

	defer con.Close()

	return err
}
