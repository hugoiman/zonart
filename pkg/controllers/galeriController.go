package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GaleriController is class
type GaleriController struct{}

// GetGaleris is func
func (gc GaleriController) GetGaleris(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var galeri models.Galeri

	dataGaleri := galeri.GetGaleris(idToko)
	message, _ := json.Marshal(dataGaleri)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateGaleri is func
func (gc GaleriController) CreateGaleri(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var galeri models.Galeri

	if err := json.NewDecoder(r.Body).Decode(&galeri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(galeri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := galeri.CreateGaleri(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Gambar telah ditambahkan."}`))
}

// DeleteGaleri is func
func (gc GaleriController) DeleteGaleri(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGaleri := vars["idGaleri"]
	var galeri models.Galeri

	err := galeri.DeleteGaleri(idToko, idGaleri)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
