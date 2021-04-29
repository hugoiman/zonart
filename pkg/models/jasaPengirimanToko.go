package models

import (
	"encoding/json"
	"zonart/db"
)

// JasaPengirimanToko is class
type JasaPengirimanToko struct {
	idJasaPengiriman int
	kurir            string
	kode             string
	status           bool
}

func (jpt *JasaPengirimanToko) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDJasaPengiriman int    `json:"idJasaPengiriman"`
		Kurir            string `json:"kurir"`
		Kode             string `json:"kode"`
		Status           bool   `json:"status"`
	}{
		IDJasaPengiriman: jpt.idJasaPengiriman,
		Kurir:            jpt.kurir,
		Kode:             jpt.kode,
		Status:           jpt.status,
	})
}

func (jpt *JasaPengirimanToko) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDJasaPengiriman int    `json:"idJasaPengiriman"`
		Kurir            string `json:"kurir"`
		Kode             string `json:"kode"`
		Status           bool   `json:"status"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	jpt.idJasaPengiriman = alias.IDJasaPengiriman
	jpt.kurir = alias.Kurir
	jpt.kode = alias.Kode
	jpt.status = alias.Status

	return nil
}

// GetJasaPengirimanToko is func
func (jpt JasaPengirimanToko) GetJasaPengirimanToko(idToko int) []JasaPengirimanToko {
	con := db.Connect()
	var jpts []JasaPengirimanToko

	var jp JasaPengiriman
	jasaPengiriman := jp.GetJasaPengiriman()

	query := "SELECT a.idJasaPengiriman, b.kurir, b.kode, a.status FROM jasaPengirimanToko a JOIN jasaPengiriman b ON a.idJasaPengiriman = b.idJasaPengiriman WHERE a.idToko = ? AND a.idJasaPengiriman = ?"

	for _, v := range jasaPengiriman {
		err := con.QueryRow(query, idToko, v.idJasaPengiriman).Scan(&jpt.idJasaPengiriman, &jpt.kurir, &jpt.kode, &jpt.status)
		if err != nil {
			jpt.idJasaPengiriman = v.idJasaPengiriman
			jpt.kurir = v.kurir
			jpt.kode = v.kode
			jpt.status = false
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

	con.QueryRow(query, idToko, jpt.idJasaPengiriman).Scan(&isAny)

	if isAny == true {
		queryUpdate := "UPDATE jasaPengirimanToko SET status = ? WHERE idJasaPengiriman = ? AND idToko = ?"
		_, err = con.Exec(queryUpdate, jpt.status, jpt.idJasaPengiriman, idToko)
	} else if isAny == false {
		queryInsert := "INSERT INTO jasaPengirimanToko (idJasaPengiriman, idToko, status) VALUES (?,?,?)"
		_, err = con.Exec(queryInsert, jpt.idJasaPengiriman, idToko, jpt.status)
	}

	defer con.Close()

	return err
}
