package models

import "zonart/db"

// BiayaTambahan is class
type BiayaTambahan struct {
	IDBiayaTambahan int    `json:"idBiayaTambahan"`
	IDOrder         int    `json:"idOrder"`
	Item            string `json:"item" validate:"required"`
	Berat           string `json:"berat"`
	Total           int    `json:"total" validate:"required"`
}

// GetBiayaTambahan is func
func (bt BiayaTambahan) GetBiayaTambahan(idBiayaTambahan, idOrder string) (BiayaTambahan, error) {
	con := db.Connect()
	query := "SELECT idBiayaTambahan, idOrder, item, berat, total FROM biayaTambahan WHERE idBiayaTambahan = ? AND idOrder = ?"

	err := con.QueryRow(query, idBiayaTambahan, idOrder).Scan(
		&bt.IDBiayaTambahan, &bt.IDOrder, &bt.Item, &bt.Berat, &bt.Total)

	defer con.Close()
	return bt, err
}

// GetBiayaTambahans is func
func (bt BiayaTambahan) GetBiayaTambahans(idOrder string) []BiayaTambahan {
	con := db.Connect()
	var bts []BiayaTambahan

	query := "SELECT idBiayaTambahan, idOrder, item, berat, total FROM biayaTambahan WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&bt.IDBiayaTambahan, &bt.IDOrder, &bt.Item, &bt.Berat, &bt.Total,
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

	_, err := con.Exec(query, idOrder, bt.Item, bt.Berat, bt.Total)
	if err != nil {
		return err
	}

	query = "UPDATE `order` SET total = total + ?, tagihan = tagihan + ?, statusPembayaran = 'belum lunas' WHERE idOrder = ?"
	_, _ = con.Exec(query, bt.Total, bt.Total, idOrder)

	defer con.Close()

	return err
}

// DeleteBiayaTambahan is func
func (bt BiayaTambahan) DeleteBiayaTambahan(idBiayaTambahan, idOrder string) error {
	con := db.Connect()
	query := "DELETE FROM biayaTambahan WHERE idBiayaTambahan = ?"
	_, err := con.Exec(query, idBiayaTambahan)

	query = "UPDATE `order` SET total = total - ?, tagihan = tagihan - ? WHERE idOrder = ?"
	_, _ = con.Exec(query, bt.Total, bt.Total, idOrder)

	defer con.Close()

	return err
}
