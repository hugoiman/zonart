package models

import "zonart/db"

// FileOrder is class
type FileOrder struct {
	IDFileOrder int    `json:"idFileOrder"`
	IDOrder     int    `json:"idOrder"`
	Foto        string `json:"foto"`
}

// GetFileOrder is func
func (fo FileOrder) GetFileOrder(idOrder string) []FileOrder {
	con := db.Connect()
	var fileOrders []FileOrder

	query := "SELECT idFileOrder, idOrder, foto FROM fileOrder WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&fo.IDFileOrder, &fo.IDOrder, &fo.Foto,
		)

		fileOrders = append(fileOrders, fo)
	}

	defer con.Close()

	return fileOrders
}

// CreateFileOrder is func
func (fo FileOrder) CreateFileOrder(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO fileOrder (idOrder, foto) VALUES (?,?)"
	_, err := con.Exec(query, idOrder, fo.Foto)

	defer con.Close()

	return err
}
