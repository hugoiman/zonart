package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetFaqs is func
func GetFaqs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var faq models.Faq

	dataFaq := faq.GetFaqs(idToko)
	message, _ := json.Marshal(dataFaq)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetFaq is func
func GetFaq(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idFaq := vars["idFaq"]
	idToko := vars["idToko"]
	var faq models.Faq

	dataFaq, err := faq.GetFaq(idFaq, idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataFaq)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateFaq is func
func CreateFaq(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var faq models.Faq

	if err := json.NewDecoder(r.Body).Decode(&faq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(faq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := faq.CreateFaq(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! FAQ telah ditambahkan."}`))
}

// DeleteFaq is func
func DeleteFaq(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idFaq := vars["idFaq"]
	var faq models.Faq

	err := faq.DeleteFaq(idToko, idFaq)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Data telah dihapus!"}`))

}
