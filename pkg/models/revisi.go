package models

import (
	"encoding/json"
	"time"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Revisi is class
type Revisi struct {
	idRevisi  int
	catatan   string
	createdAt string
}

func (r *Revisi) SetCreatedAt(data string) {
	r.createdAt = data
}

func (r *Revisi) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDRevisi  int    `json:"idRevisi"`
		Catatan   string `json:"catatan"`
		CreatedAt string `json:"createdAt"`
	}{
		IDRevisi:  r.idRevisi,
		Catatan:   r.catatan,
		CreatedAt: r.createdAt,
	})
}

func (r *Revisi) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDRevisi  int    `json:"idRevisi"`
		Catatan   string `json:"catatan" validate:"required"`
		CreatedAt string `json:"createdAt"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	r.idRevisi = alias.IDRevisi
	r.catatan = alias.Catatan
	r.createdAt = alias.CreatedAt

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetRevisi is func
func (r Revisi) GetRevisi(idOrder string) []Revisi {
	con := db.Connect()
	var revisis []Revisi
	var createdAt time.Time

	query := "SELECT idRevisi, catatan, createdAt FROM revisi WHERE idOrder = ?"

	rows, _ := con.Query(query, idOrder)
	for rows.Next() {
		rows.Scan(
			&r.idRevisi, &r.catatan, &createdAt,
		)

		r.createdAt = createdAt.Format("02 Jan 2006")
		revisis = append(revisis, r)
	}

	defer con.Close()

	return revisis
}

// CreateRevisi is func
func (r Revisi) CreateRevisi(idOrder string) error {
	con := db.Connect()
	query := "INSERT INTO revisi (idOrder, catatan, createdAt) VALUES (?,?,?)"
	_, err := con.Exec(query, idOrder, r.catatan, r.createdAt)

	defer con.Close()

	return err
}
