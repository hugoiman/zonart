package controllers

import (
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// DeleteRekening is func
func DeleteRekening(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idRekening := vars["idRekening"]
	var rekening models.Rekening

	err := rekening.DeleteRekening(idToko, idRekening)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Data telah dihapus!"}`))

}
