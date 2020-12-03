package models

import (
	"zonart/db"
)

// PengirimanToko is class
type PengirimanToko struct {
	IDPengirimanToko int  `json:"idPengirimanToko"`
	Jne              bool `json:"jne"`
	Jnt              bool `json:"jnt"`
	Ninja            bool `json:"ninja"`
	Sicepat          bool `json:"sicepat"`
	Cod              bool `json:"cod"`
}

// GetPengirimanToko is func
func (pt PengirimanToko) GetPengirimanToko(id string) PengirimanToko {
	con := db.Connect()
	query := "SELECT a.idPengirimanToko, a.jne, a.jnt, a.ninja, a.sicepat, a.cod FROM pengirimanToko a JOIN toko b ON a.idToko = b.idToko WHERE b.idToko = ? OR b.slug = ?"

	_ = con.QueryRow(query, id, id).Scan(
		&pt.IDPengirimanToko, &pt.Jne, &pt.Jnt, &pt.Ninja, &pt.Sicepat, &pt.Cod)

	defer con.Close()

	return pt
}

// InitializePengirimanToko is func
func (pt PengirimanToko) InitializePengirimanToko(idToko int) error {
	con := db.Connect()
	query := "INSERT INTO pengirimanToko (idToko) VALUES (?)"
	_, err := con.Exec(query, idToko)

	defer con.Close()

	return err
}

// UpdatePengirimanToko is func
func (pt PengirimanToko) UpdatePengirimanToko(id string) error {
	con := db.Connect()
	query := "UPDATE pengirimanToko SET jne = ?, jnt = ?, ninja = ?, sicepat = ?, cod = ? WHERE idToko = ?"
	_, err := con.Exec(query, pt.Jne, pt.Jnt, pt.Ninja, pt.Sicepat, pt.Cod, id)

	defer con.Close()

	return err
}
