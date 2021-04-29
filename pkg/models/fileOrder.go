package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// FileOrder is class
type FileOrder struct {
	idFileOrder int
	foto        string
}

func (fo *FileOrder) SetFoto(data string) {
	fo.foto = data
}

func (fo *FileOrder) GetFoto() string {
	return fo.foto
}

func (fo *FileOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDFileOrder int    `json:"idFileOrder"`
		Foto        string `json:"foto"`
	}{
		IDFileOrder: fo.idFileOrder,
		Foto:        fo.foto,
	})
}

func (fo *FileOrder) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDFileOrder int    `json:"idFileOrder"`
		Foto        string `json:"foto"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	fo.idFileOrder = alias.IDFileOrder
	fo.foto = alias.Foto

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetFileOrder is func
func (fo FileOrder) GetFileOrder(idOrder string) []FileOrder {
	con := db.Connect()
	var fileOrders []FileOrder

	query := "SELECT idFileOrder, foto FROM fileOrder WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&fo.idFileOrder, &fo.foto,
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
	_, err := con.Exec(query, idOrder, fo.foto)

	defer con.Close()

	return err
}
