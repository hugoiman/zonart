package models

import "zonart/db"

// Faq is class
type Faq struct {
	IDFaq      int    `json:"idFaq"`
	IDToko     int    `json:"idToko"`
	Pertanyaan string `json:"pertanyaan" validate:"required"`
	Jawaban    string `json:"jawaban" validate:"required"`
	Kategori   string `json:"kategori" validate:"required"`
}

// Faqs is list of faq
type Faqs struct {
	Faqs []Faq `json:"faq"`
}

// GetFaqs is func
func (f Faq) GetFaqs(idToko string) Faqs {
	con := db.Connect()
	query := "SELECT idFaq, idToko, pertanyaan, jawaban, kategori FROM faq WHERE idToko = ? ORDER BY idFaq DESC"
	rows, _ := con.Query(query, idToko)

	var faqs Faqs

	for rows.Next() {
		rows.Scan(
			&f.IDFaq, &f.IDToko, &f.Pertanyaan, &f.Jawaban, &f.Kategori,
		)

		faqs.Faqs = append(faqs.Faqs, f)
	}

	defer con.Close()

	return faqs
}

// GetFaq is func
func (f Faq) GetFaq(idFaq, idToko string) (Faq, error) {
	con := db.Connect()
	query := "SELECT idFaq, idToko, pertanyaan, jawaban, kategori FROM faq WHERE idFaq = ? AND idToko = ?"

	err := con.QueryRow(query, idFaq, idToko).Scan(
		&f.IDFaq, &f.IDToko, &f.Pertanyaan, &f.Jawaban, &f.Kategori)

	defer con.Close()
	return f, err
}

// CreateFaq is func
func (f Faq) CreateFaq(idToko string) (int, error) {
	con := db.Connect()
	query := "INSERT INTO faq (idToko, pertanyaan, jawaban, kategori) VALUES (?,?,?,?)"
	exec, err := con.Exec(query, idToko, f.Pertanyaan, f.Jawaban, f.Kategori)

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
