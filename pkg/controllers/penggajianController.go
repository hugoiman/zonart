package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// PenggajianController is class
type PenggajianController struct{}

// GetGajis is func
func (pc PenggajianController) GetGajis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var produk models.Penggajian

	dataPenggajian := produk.GetGajis(idToko)
	message, _ := json.Marshal(dataPenggajian)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateGaji is func
func (pc PenggajianController) CreateGaji(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var penggajian models.Penggajian

	if err := json.NewDecoder(r.Body).Decode(&penggajian); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(penggajian); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := penggajian.CreateGaji(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Gaji telah disimpan"}`))
}

// DeleteGaji is func
func (pc PenggajianController) DeleteGaji(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idPenggajian := vars["idPenggajian"]
	var penggajian models.Penggajian

	err := penggajian.DeleteGaji(idToko, idPenggajian)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Data telah dihapus!"}`))

}
