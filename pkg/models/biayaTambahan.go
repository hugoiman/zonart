package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// BiayaTambahan is class
type BiayaTambahan struct {
	idBiayaTambahan int
	item            string
	berat           int
	total           int
}

func (bt *BiayaTambahan) GetItem() string {
	return bt.item
}

func (bt *BiayaTambahan) GetBerat() int {
	return bt.berat
}

func (bt *BiayaTambahan) GetTotal() int {
	return bt.total
}

func (bt *BiayaTambahan) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDBiayaTambahan int    `json:"idBiayaTambahan"`
		Item            string `json:"item"`
		Berat           int    `json:"berat"`
		Total           int    `json:"total"`
	}{
		IDBiayaTambahan: bt.idBiayaTambahan,
		Item:            bt.item,
		Berat:           bt.berat,
		Total:           bt.total,
	})
}

func (bt *BiayaTambahan) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDBiayaTambahan int    `json:"idBiayaTambahan"`
		Item            string `json:"item" validate:"required"`
		Berat           int    `json:"berat"`
		Total           int    `json:"total"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	bt.idBiayaTambahan = alias.IDBiayaTambahan
	bt.item = alias.Item
	bt.berat = alias.Berat
	bt.total = alias.Total

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetBiayaTambahan is func
func (bt BiayaTambahan) GetBiayaTambahan(idBiayaTambahan, idOrder string) (BiayaTambahan, error) {
	con := db.Connect()
	query := "SELECT idBiayaTambahan, item, berat, total FROM biayaTambahan WHERE idBiayaTambahan = ? AND idOrder = ?"

	err := con.QueryRow(query, idBiayaTambahan, idOrder).Scan(
		&bt.idBiayaTambahan, &bt.item, &bt.berat, &bt.total)

	defer con.Close()
	return bt, err
}

// GetBiayaTambahans is func
func (bt BiayaTambahan) GetBiayaTambahans(idOrder string) []BiayaTambahan {
	con := db.Connect()
	var bts []BiayaTambahan

	query := "SELECT idBiayaTambahan, item, berat, total FROM biayaTambahan WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&bt.idBiayaTambahan, &bt.item, &bt.berat, &bt.total,
		)

		bts = append(bts, bt)
	}

	defer con.Close()

	return bts
}

// CreateBiayaTambahan is func
func (bt BiayaTambahan) CreateBiayaTambahan(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO biayaTambahan (idOrder, item, berat, total) VALUES (?,?,?,?)"

	_, err := con.Exec(query, idOrder, bt.item, bt.berat, bt.total)
	if err != nil {
		return err
	}

	defer con.Close()

	return err
}

// DeleteBiayaTambahan is func
func (bt BiayaTambahan) DeleteBiayaTambahan(idBiayaTambahan, idOrder string) error {
	con := db.Connect()
	query := "DELETE FROM biayaTambahan WHERE idBiayaTambahan = ? AND idOrder = ?"
	_, err := con.Exec(query, idBiayaTambahan, idOrder)

	defer con.Close()

	return err
}
