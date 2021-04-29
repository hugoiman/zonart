package models

import (
	"encoding/json"
	"zonart/db"
)

// GrupOpsiProduk is class
type GrupOpsiProduk struct {
	idProduk   int
	idGrupOpsi int
	namaProduk string
	namaGrup   string
}

// GrupOpsiProduks is list of GrupOpsiProduk
type GrupOpsiProduks struct {
	GrupOpsiProduks []GrupOpsiProduk `json:"grupOpsiProduk"`
}

func (gop *GrupOpsiProduk) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDProduk   int    `json:"idProduk"`
		IDGrupOpsi int    `json:"idGrupOpsi"`
		NamaProduk string `json:"namaProduk"`
		NamaGrup   string `json:"namaGrup"`
	}{
		IDProduk:   gop.idProduk,
		IDGrupOpsi: gop.idGrupOpsi,
		NamaProduk: gop.namaProduk,
		NamaGrup:   gop.namaGrup,
	})
}

func (gop *GrupOpsiProduk) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDProduk   int    `json:"idProduk"`
		IDGrupOpsi int    `json:"idGrupOpsi"`
		NamaProduk string `json:"namaProduk"`
		NamaGrup   string `json:"namaGrup"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	gop.idProduk = alias.IDProduk
	gop.idGrupOpsi = alias.IDGrupOpsi
	gop.namaProduk = alias.NamaProduk
	gop.namaGrup = alias.NamaGrup

	return nil
}

// GetGrupOpsiProduks is get all produk in a grup opsi
func (gop GrupOpsiProduk) GetGrupOpsiProduks(idToko, idGrupOpsi string) GrupOpsiProduks {
	con := db.Connect()
	query := "SELECT a.idProduk, a.idGrupOpsi, b.namaGrup, c.namaProduk FROM grupOpsiProduk a " +
		"JOIN grupOpsi b ON a.idGrupOpsi = b.idGrupOpsi " +
		"JOIN produk c ON a.idProduk = c.idProduk WHERE b.idToko = ? AND a.idGrupOpsi = ? ORDER BY a.idGrupOpsi ASC"
	rows, _ := con.Query(query, idToko, idGrupOpsi)

	var gops GrupOpsiProduks

	for rows.Next() {
		rows.Scan(
			&gop.idProduk, &gop.idGrupOpsi, &gop.namaGrup, &gop.namaProduk,
		)

		gops.GrupOpsiProduks = append(gops.GrupOpsiProduks, gop)
	}

	defer con.Close()

	return gops
}

// GetGrupOpsiProduksByProduk is get all produk in a grup opsi by id produk
func (gop GrupOpsiProduk) GetGrupOpsiProduksByProduk(idToko, idProduk string) GrupOpsiProduks {
	con := db.Connect()
	query := "SELECT a.idProduk, a.idGrupOpsi, b.namaGrup, c.namaProduk FROM grupOpsiProduk a " +
		"JOIN grupOpsi b ON a.idGrupOpsi = b.idGrupOpsi " +
		"JOIN produk c ON a.idProduk = c.idProduk WHERE b.idToko = ? AND a.idProduk = ? ORDER BY a.idGrupOpsi ASC"
	rows, _ := con.Query(query, idToko, idProduk)

	var gops GrupOpsiProduks

	for rows.Next() {
		rows.Scan(
			&gop.idProduk, &gop.idGrupOpsi, &gop.namaGrup, &gop.namaProduk,
		)

		gops.GrupOpsiProduks = append(gops.GrupOpsiProduks, gop)
	}

	defer con.Close()

	return gops
}

// SambungGrupOpsikeProduk is func
func (gop GrupOpsiProduk) SambungGrupOpsikeProduk(idProduk, idGrupOpsi string) error {
	con := db.Connect()
	query := "INSERT INTO grupOpsiProduk (idProduk, idGrupOpsi) VALUES (?,?)"
	_, err := con.Exec(query, idProduk, idGrupOpsi)

	defer con.Close()

	return err
}

// PutusGrupOpsidiProduk is func
func (gop GrupOpsiProduk) PutusGrupOpsidiProduk(idProduk, idGrupOpsi string) error {
	con := db.Connect()
	query := "DELETE FROM grupOpsiProduk WHERE idProduk = ? AND idGrupOpsi = ?"
	_, err := con.Exec(query, idProduk, idGrupOpsi)

	defer con.Close()

	return err
}

// CheckSambunganGrupOpsi is func
func (gop GrupOpsiProduk) CheckSambunganGrupOpsi(idProduk, idGrupOpsi string) bool {
	var isAny bool
	con := db.Connect()
	query := "SELECT EXISTS (SELECT 1 FROM grupOpsiProduk WHERE idProduk = ? AND idGrupOpsi = ?)"
	con.QueryRow(query, idProduk, idGrupOpsi).Scan(&isAny)

	defer con.Close()

	return isAny
}
