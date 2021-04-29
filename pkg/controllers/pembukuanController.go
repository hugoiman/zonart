package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// PembukuanController is class
type PembukuanController struct{}

// GetPembukuans is func
func (pc PembukuanController) GetPembukuans(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var pembukuan models.Pembukuan

	dataPembukuan := pembukuan.GetPembukuans(idToko)
	message, _ := json.Marshal(&dataPembukuan)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreatePembukuan is func
func (pc PembukuanController) CreatePembukuan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var pembukuan models.Pembukuan

	if err := json.NewDecoder(r.Body).Decode(&pembukuan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := pembukuan.CreatePembukuan(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah ditambahkan."}`))
}

// DeletePembukuan is func
func (pc PembukuanController) DeletePembukuan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idPembukuan := vars["idPembukuan"]
	var pembukuan models.Pembukuan

	err := pembukuan.DeletePembukuan(idToko, idPembukuan)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
