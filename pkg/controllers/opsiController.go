package controllers

import (
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// OpsiController is class
type OpsiController struct{}

// DeleteOpsi is func
func (oc OpsiController) DeleteOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGrupOpsi := vars["idGrupOpsi"]
	idOpsi := vars["idOpsi"]

	var grupOpsi models.GrupOpsi
	var opsi models.Opsi

	_, err := grupOpsi.GetGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, "Grup opsi tidak ditemukan.", http.StatusBadRequest)
		return
	}

	err = opsi.DeleteOpsi(idGrupOpsi, idOpsi)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
