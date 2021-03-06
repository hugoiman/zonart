package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/custerr"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// PenggajianController is class
type PenggajianController struct{}

// GetGajis is func
func (pc PenggajianController) GetGajis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var produk models.Penggajian

	dataPenggajian := produk.GetGajis(idToko)
	message, _ := json.Marshal(&dataPenggajian)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"penggajian":` + string(message) + `}`))
}

// CreateGaji is func
func (pc PenggajianController) CreateGaji(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idKaryawan := vars["idKaryawan"]

	var penggajian models.Penggajian

	if err := json.NewDecoder(r.Body).Decode(&penggajian); err != nil {
		http.Error(w, custerr.CustomError(err).Error(), http.StatusBadRequest)
		return
	}

	var karyawan models.Karyawan
	_, err := karyawan.GetKaryawan(idToko, idKaryawan)
	if err != nil {
		http.Error(w, "Karyawan tidak ditemukan", http.StatusBadRequest)
		return
	}

	err = penggajian.CreateGaji(idToko, idKaryawan)
	if err != nil {
		http.Error(w, custerr.CustomError(err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Gaji telah disimpan"}`))
}

// DeleteGaji is func
func (pc PenggajianController) DeleteGaji(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idPenggajian := vars["idPenggajian"]
	var penggajian models.Penggajian

	_ = penggajian.DeleteGaji(idToko, idPenggajian)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
