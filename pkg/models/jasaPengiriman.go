package models

import (
	"encoding/json"
	"zonart/db"
)

// JasaPengiriman is class
type JasaPengiriman struct {
	idJasaPengiriman int
	kurir            string
	kode             string
}

func (jp *JasaPengiriman) GetIDJasaPengiriman() int {
	return jp.idJasaPengiriman
}

func (jp *JasaPengiriman) GetKurir() string {
	return jp.kurir
}

func (jp *JasaPengiriman) GetKode() string {
	return jp.kode
}

func (jp *JasaPengiriman) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDJasaPengiriman int    `json:"idJasaPengiriman"`
		Kurir            string `json:"kurir"`
		Kode             string `json:"kode"`
	}{
		IDJasaPengiriman: jp.idJasaPengiriman,
		Kurir:            jp.kurir,
		Kode:             jp.kode,
	})
}

func (jp *JasaPengiriman) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDJasaPengiriman int    `json:"idJasaPengiriman"`
		Kurir            string `json:"kurir"`
		Kode             string `json:"kode"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	jp.idJasaPengiriman = alias.IDJasaPengiriman
	jp.kurir = alias.Kurir
	jp.kode = alias.Kode

	return nil
}

// GetJasaPengiriman is func
func (jp JasaPengiriman) GetJasaPengirimans() []JasaPengiriman {
	con := db.Connect()
	var jps []JasaPengiriman
	query := "SELECT idJasaPengiriman, kurir, kode FROM jasaPengiriman"
	rows, _ := con.Query(query)
	for rows.Next() {
		rows.Scan(
			&jp.idJasaPengiriman, &jp.kurir, &jp.kode,
		)

		jps = append(jps, jp)
	}

	defer con.Close()

	return jps
}
