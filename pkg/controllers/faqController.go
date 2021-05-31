package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"zonart/custerr"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// FaqController is class
type FaqController struct {
}

// GetFaqs is func
func (fc FaqController) GetFaqs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var faq models.Faq

	dataFaq := faq.GetFaqs(idToko)
	message, _ := json.Marshal(&dataFaq)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"faq":` + string(message) + `}`))
}

// GetFaq is func
func (fc FaqController) GetFaq(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idFaq := vars["idFaq"]
	idToko := vars["idToko"]
	var faq models.Faq

	dataFaq, err := faq.GetFaq(idFaq, idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(&dataFaq)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateFaq is func
func (fc FaqController) CreateFaq(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var faq models.Faq

	if err := json.NewDecoder(r.Body).Decode(&faq); err != nil {
		http.Error(w, custerr.CustomError(err).Error(), http.StatusBadRequest)
		return
	}

	faq.SetKategori(strings.Title(strings.ToLower(faq.GetKategori())))

	idFaq, err := faq.CreateFaq(idToko)
	if err != nil {
		http.Error(w, custerr.CustomError(err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"FAQ telah ditambahkan.","idFaq":"` + strconv.Itoa(idFaq) + `"}`))
}

// DeleteFaq is func
func (fc FaqController) DeleteFaq(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idFaq := vars["idFaq"]
	var faq models.Faq

	_ = faq.DeleteFaq(idToko, idFaq)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
