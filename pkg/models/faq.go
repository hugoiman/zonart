package models

import (
	"encoding/json"
	"zonart/db"

	"gopkg.in/go-playground/validator.v9"
)

// Faq is class
type Faq struct {
	idFaq      int
	pertanyaan string
	jawaban    string
	kategori   string
}

// Faqs is list of faq
type Faqs struct {
	Faqs []Faq `json:"faq"`
}

func (f *Faq) SetKategori(data string) {
	f.kategori = data
}

func (f *Faq) GetKategori() string {
	return f.kategori
}

func (f *Faq) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDFaq      int    `json:"idFaq"`
		Pertanyaan string `json:"pertanyaan"`
		Jawaban    string `json:"jawaban"`
		Kategori   string `json:"kategori"`
	}{
		IDFaq:      f.idFaq,
		Pertanyaan: f.pertanyaan,
		Jawaban:    f.jawaban,
		Kategori:   f.kategori,
	})
}

func (f *Faq) UnmarshalJSON(data []byte) error {
	alias := struct {
		IDFaq      int    `json:"idFaq"`
		Pertanyaan string `json:"pertanyaan" validate:"required"`
		Jawaban    string `json:"jawaban" validate:"required"`
		Kategori   string `json:"kategori" validate:"required"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	f.idFaq = alias.IDFaq
	f.pertanyaan = alias.Pertanyaan
	f.jawaban = alias.Jawaban
	f.kategori = alias.Kategori

	if err = validator.New().Struct(alias); err != nil {
		return err
	}
	return nil
}

// GetFaqs is func
func (f Faq) GetFaqs(idToko string) Faqs {
	con := db.Connect()
	query := "SELECT idFaq, pertanyaan, jawaban, kategori FROM faq WHERE idToko = ? ORDER BY idFaq DESC"
	rows, _ := con.Query(query, idToko)

	var faqs Faqs

	for rows.Next() {
		rows.Scan(
			&f.idFaq, &f.pertanyaan, &f.jawaban, &f.kategori,
		)

		faqs.Faqs = append(faqs.Faqs, f)
	}

	defer con.Close()

	return faqs
}

// GetFaq is func
func (f Faq) GetFaq(idFaq, idToko string) (Faq, error) {
	con := db.Connect()
	query := "SELECT idFaq, pertanyaan, jawaban, kategori FROM faq WHERE idFaq = ? AND idToko = ?"

	err := con.QueryRow(query, idFaq, idToko).Scan(
		&f.idFaq, &f.pertanyaan, &f.jawaban, &f.kategori)

	defer con.Close()
	return f, err
}

// CreateFaq is func
func (f Faq) CreateFaq(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO faq (idToko, pertanyaan, jawaban, kategori) VALUES (?,?,?,?)"
	exec, err := con.Exec(query, idToko, f.pertanyaan, f.jawaban, f.kategori)

	if err != nil {
		return 0, err
	}

	idInt64, _ := exec.LastInsertId()
	idProduk := int(idInt64)

	defer con.Close()

	return idProduk, err
}

// DeleteFaq is func
func (f Faq) DeleteFaq(idToko, idFaq string) error {
	con := db.Connect()
	query := "DELETE FROM faq WHERE idToko = ? AND idFaq = ?"
	_, err := con.Exec(query, idToko, idFaq)

	defer con.Close()

	return err
}
