package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// PembukuanController is class
type PembukuanController struct{}

// GetPembukuans is func
func (pc PembukuanController) GetPembukuans(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var pembukuan models.Pembukuan

	dataPembukuan := pembukuan.GetPembukuans(idToko)
	message, _ := json.Marshal(dataPembukuan)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreatePengeluaran is func
func (pc PembukuanController) CreatePengeluaran(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var pembukuan models.Pembukuan

	if err := json.NewDecoder(r.Body).Decode(&pembukuan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(pembukuan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pembukuan.Jenis = "pengeluaran"

	err := pembukuan.CreatePembukuan(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Pengeluaran telah ditambahkan."}`))
}

// DeletePengeluaran is func
func (pc PembukuanController) DeletePengeluaran(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idPembukuan := vars["idPembukuan"]
	var pembukuan models.Pembukuan

	err := pembukuan.DeletePengeluaran(idToko, idPembukuan)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Data telah dihapus!"}`))

}
