package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// HasilOrder is class
type HasilOrder struct {
	idHasilOrder int
	hasil        string
	status       string
	createdAt    string
}

func (ho *HasilOrder) SetHasil(data string) {
	ho.hasil = data
}

func (ho *HasilOrder) GetHasil() string {
	return ho.hasil
}

func (ho *HasilOrder) SetStatus(data string) {
	ho.status = data
}

func (ho *HasilOrder) SetCreatedAt(data string) {
	ho.createdAt = data
}

func (ho *HasilOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDHasilOrder int    `json:"idHasilOrder"`
		Hasil        string `json:"hasil"`
		Nominal      int    `json:"nominal"`
		Status       string `json:"status"`
		CreatedAt    string `json:"createdAt"`
	}{
		IDHasilOrder: ho.idHasilOrder,
		Hasil:        ho.hasil,
		Status:       ho.status,
		CreatedAt:    ho.createdAt,
	})
}

func (ho *HasilOrder) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDHasilOrder int    `json:"idHasilOrder"`
		Hasil        string `json:"hasil"`
		Status       string `json:"status"`
		CreatedAt    string `json:"createdAt"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	ho.idHasilOrder = alias.IDHasilOrder
	ho.hasil = alias.Hasil
	ho.status = alias.Status
	ho.createdAt = alias.CreatedAt

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetHasilOrder is func
func (ho HasilOrder) GetHasilOrder(idOrder string) (HasilOrder, error) {
	con := db.Connect()
	query := "SELECT idHasilOrder, hasil, status, createdAt FROM hasilOrder WHERE idOrder = ?"
	var createdAt time.Time

	err := con.QueryRow(query, idOrder).Scan(&ho.idHasilOrder, &ho.hasil, &ho.status, &createdAt)
	ho.createdAt = createdAt.Format("02 Jan 2006")

	defer con.Close()

	return ho, err
}

// AddHasilOrder is func
func (ho HasilOrder) AddHasilOrder(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO `hasilOrder` (idHasilOrder, idOrder, hasil, status, createdAt) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE hasil = ?, status = ?, createdAt = ?"
	_, err := con.Exec(query, ho.idHasilOrder, idOrder, ho.hasil, ho.status, ho.createdAt, ho.hasil, ho.status, ho.createdAt)

	defer con.Close()

	return err
}

// SetujuiHasilOrder is func
func (ho HasilOrder) SetujuiHasilOrder(idOrder string) error {
	con := db.Connect()
	query := "UPDATE `hasilOrder` SET status = 'sudah disetujui' WHERE idOrder = ?"
	_, err := con.Exec(query, idOrder)

	defer con.Close()

	return err
}
