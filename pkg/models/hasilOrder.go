package models

import (
	"time"
	"zonart/db"
)

// HasilOrder is class
type HasilOrder struct {
	IDHasilOrder int    `json:"idHasilOrder"`
	IDOrder      int    `json:"idOrder"`
	Hasil        string `json:"hasil"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
}

// GetHasilOrder is func
func (ho HasilOrder) GetHasilOrder(idOrder string) (HasilOrder, error) {
	con := db.Connect()
	query := "SELECT idHasilOrder, idOrder, hasil, status, createdAt FROM hasilOrder WHERE idOrder = ?"
	var createdAt time.Time

	err := con.QueryRow(query, idOrder).Scan(&ho.IDHasilOrder, &ho.IDOrder, &ho.Hasil, &ho.Status, &createdAt)
	ho.CreatedAt = createdAt.Format("02 Jan 2006")

	defer con.Close()

	return ho, err
}

// AddHasilOrder is func
func (ho HasilOrder) AddHasilOrder(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO `hasilOrder` (idHasilOrder, idOrder, hasil, status, createdAt) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE hasil = ?, status = ?, createdAt = ?"
	_, err := con.Exec(query, ho.IDHasilOrder, idOrder, ho.Hasil, ho.Status, ho.CreatedAt, ho.Hasil, ho.Status, ho.CreatedAt)

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
