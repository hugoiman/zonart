package models

import (
	"encoding/json"
	"strconv"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// GrupOpsi is class
type GrupOpsi struct {
	idGrupOpsi      int
	namaGrup        string
	deskripsi       string
	required        bool
	min             int
	max             int
	spesificRequest bool
	hardCopy        bool
	softCopy        bool
	opsi            []Opsi
}

func (grupOpsi *GrupOpsi) GetIDGrupOpsi() int {
	return grupOpsi.idGrupOpsi
}

func (grupOpsi *GrupOpsi) GetNamaGrup() string {
	return grupOpsi.namaGrup
}

func (grupOpsi *GrupOpsi) GetRequired() bool {
	return grupOpsi.required
}

func (grupOpsi *GrupOpsi) GetMin() int {
	return grupOpsi.min
}

func (grupOpsi *GrupOpsi) GetMax() int {
	return grupOpsi.max
}

func (grupOpsi *GrupOpsi) GetSpesificRequest() bool {
	return grupOpsi.spesificRequest
}

func (grupOpsi *GrupOpsi) GetHardCopy() bool {
	return grupOpsi.hardCopy
}

func (grupOpsi *GrupOpsi) GetSoftCopy() bool {
	return grupOpsi.softCopy
}

func (grupOpsi *GrupOpsi) GetOpsi() []Opsi {
	return grupOpsi.opsi
}

func (grupOpsi *GrupOpsi) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDGrupOpsi      int    `json:"idGrupOpsi"`
		NamaGrup        string `json:"namaGrup"`
		Deskripsi       string `json:"deskripsi"`
		Required        bool   `json:"required"`
		Min             int    `json:"min"`
		Max             int    `json:"max"`
		SpesificRequest bool   `json:"spesificRequest"`
		HardCopy        bool   `json:"hardcopy"`
		SoftCopy        bool   `json:"softcopy"`
		Opsi            []Opsi `json:"opsi"`
	}{
		IDGrupOpsi:      grupOpsi.idGrupOpsi,
		NamaGrup:        grupOpsi.namaGrup,
		Deskripsi:       grupOpsi.deskripsi,
		Required:        grupOpsi.required,
		Min:             grupOpsi.min,
		Max:             grupOpsi.max,
		SpesificRequest: grupOpsi.spesificRequest,
		HardCopy:        grupOpsi.hardCopy,
		SoftCopy:        grupOpsi.softCopy,
		Opsi:            grupOpsi.opsi,
	})
}

func (grupOpsi *GrupOpsi) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDGrupOpsi      int    `json:"idGrupOpsi"`
		NamaGrup        string `json:"namaGrup" validate:"required"`
		Deskripsi       string `json:"deskripsi"`
		Required        bool   `json:"required"`
		Min             int    `json:"min"`
		Max             int    `json:"max"`
		SpesificRequest bool   `json:"spesificRequest"`
		HardCopy        bool   `json:"hardcopy"`
		SoftCopy        bool   `json:"softcopy"`
		Opsi            []Opsi `json:"opsi" validate:"dive"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	grupOpsi.idGrupOpsi = alias.IDGrupOpsi
	grupOpsi.namaGrup = alias.NamaGrup
	grupOpsi.deskripsi = alias.Deskripsi
	grupOpsi.required = alias.Required
	grupOpsi.min = alias.Min
	grupOpsi.max = alias.Max
	grupOpsi.spesificRequest = alias.SpesificRequest
	grupOpsi.hardCopy = alias.HardCopy
	grupOpsi.softCopy = alias.SoftCopy
	grupOpsi.opsi = alias.Opsi
	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetGrupOpsi is func
func (grupOpsi GrupOpsi) GetGrupOpsis(idToko string) []GrupOpsi {
	con := db.Connect()
	query := "SELECT idGrupOpsi, namaGrup, deskripsi, required, min, max, spesificRequest, hardcopy, softcopy FROM grupOpsi WHERE idToko = ?"
	rows, _ := con.Query(query, idToko)

	var grupOpsis []GrupOpsi
	var opsi Opsi

	for rows.Next() {
		rows.Scan(
			&grupOpsi.idGrupOpsi, &grupOpsi.namaGrup, &grupOpsi.deskripsi, &grupOpsi.required, &grupOpsi.min, &grupOpsi.max, &grupOpsi.spesificRequest, &grupOpsi.hardCopy, &grupOpsi.softCopy,
		)

		grupOpsi.opsi = opsi.GetOpsis(strconv.Itoa(grupOpsi.idGrupOpsi))
		grupOpsis = append(grupOpsis, grupOpsi)
	}

	defer con.Close()

	return grupOpsis
}

// GetGrupOpsi is func
func (grupOpsi GrupOpsi) GetGrupOpsi(idToko, idGrupOpsi string) (GrupOpsi, error) {
	con := db.Connect()
	query := "SELECT idGrupOpsi, namaGrup, deskripsi, required, min, max, spesificRequest, hardcopy, softcopy FROM grupOpsi WHERE idToko = ? AND idGrupOpsi = ?"

	err := con.QueryRow(query, idToko, idGrupOpsi).Scan(
		&grupOpsi.idGrupOpsi, &grupOpsi.namaGrup, &grupOpsi.deskripsi, &grupOpsi.required, &grupOpsi.min, &grupOpsi.max, &grupOpsi.spesificRequest, &grupOpsi.hardCopy, &grupOpsi.softCopy)

	var opsi Opsi
	grupOpsi.opsi = opsi.GetOpsis(idGrupOpsi)

	defer con.Close()

	return grupOpsi, err
}

// GetGrupOpsiProduk is get all grup opsi in a produk
func (grupOpsi GrupOpsi) GetGrupOpsiProduk(idToko, idProduk string) []GrupOpsi {
	con := db.Connect()
	query := "SELECT a.idGrupOpsi, a.namaGrup, a.deskripsi, a.required, a.min, a.max, a.spesificRequest, a.hardcopy, a.softcopy FROM grupOpsi a JOIN grupOpsiProduk b ON a.idGrupOpsi = b.idGrupOpsi JOIN toko c ON a.idToko = c.idToko WHERE b.idProduk = ? AND (c.idToko = ? OR c.slug = ?)"
	rows, _ := con.Query(query, idProduk, idToko, idToko)

	var gop []GrupOpsi
	var opsi Opsi

	for rows.Next() {
		rows.Scan(
			&grupOpsi.idGrupOpsi, &grupOpsi.namaGrup, &grupOpsi.deskripsi, &grupOpsi.required, &grupOpsi.min, &grupOpsi.max, &grupOpsi.spesificRequest, &grupOpsi.hardCopy, &grupOpsi.softCopy,
		)

		grupOpsi.opsi = opsi.GetOpsis(strconv.Itoa(grupOpsi.idGrupOpsi))
		gop = append(gop, grupOpsi)
	}

	defer con.Close()

	return gop
}

// CreateGrupOpsi is func
func (grupOpsi GrupOpsi) CreateGrupOpsi(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO grupOpsi (idToko, namaGrup, deskripsi, required, min, max, spesificRequest, hardcopy, softcopy) VALUES (?,?,?,?,?,?,?,?,?)"
	exec, err := con.Exec(query, idToko, grupOpsi.namaGrup, grupOpsi.deskripsi, grupOpsi.required, grupOpsi.min, grupOpsi.max, grupOpsi.spesificRequest, grupOpsi.hardCopy, grupOpsi.softCopy)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idGrupOpsi := int(idInt64)

	for _, vOpsi := range grupOpsi.opsi {
		err = vOpsi.CreateUpdateOpsi(strconv.Itoa(idGrupOpsi))
		if err != nil {
			_ = grupOpsi.DeleteGrupOpsi(idToko, strconv.Itoa(idGrupOpsi))
			return idGrupOpsi, err
		}
	}

	defer con.Close()

	return idGrupOpsi, err
}

// UpdateGrupOpsi is func
func (grupOpsi GrupOpsi) UpdateGrupOpsi(idToko, idGrupOpsi string) error {
	con := db.Connect()
	query := "UPDATE grupOpsi SET namaGrup = ?, deskripsi = ?, required = ?, min = ?, max = ?, spesificRequest = ?, hardcopy = ?, softcopy = ? WHERE idToko = ? AND idGrupOpsi = ?"
	_, err := con.Exec(query, grupOpsi.namaGrup, grupOpsi.deskripsi, grupOpsi.required, grupOpsi.min, grupOpsi.max, grupOpsi.spesificRequest, grupOpsi.hardCopy, grupOpsi.softCopy, idToko, idGrupOpsi)

	defer con.Close()

	return err
}

// DeleteGrupOpsi is func
func (grupOpsi GrupOpsi) DeleteGrupOpsi(idToko, idGrupOpsi string) error {
	con := db.Connect()
	query := "DELETE FROM grupOpsi WHERE idToko = ? AND idGrupOpsi = ?"
	_, err := con.Exec(query, idToko, idGrupOpsi)

	defer con.Close()

	return err
}
