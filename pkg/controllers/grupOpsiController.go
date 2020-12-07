package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// GetGrupOpsis is func
func GetGrupOpsis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var grupOpsi models.GrupOpsi

	dataGrupOpsi := grupOpsi.GetGrupOpsis(idToko)
	message, _ := json.Marshal(dataGrupOpsi)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetGrupOpsi is func
func GetGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idGrupOpsi := vars["idGrupOpsi"]
	idToko := vars["idToko"]
	var gop models.GrupOpsi
	var opsi models.Opsi

	dataGrupOpsi, err := gop.GetGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataOpsi := opsi.GetOpsis(idGrupOpsi)
	dataGrupOpsi.Opsi = dataOpsi

	message, _ := json.Marshal(dataGrupOpsi)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetGrupOpsiProduk is func
func GetGrupOpsiProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idProduk := vars["idProduk"]
	idToko := vars["idToko"]
	var gop models.GrupOpsi

	dataGrupOpsiProduk := gop.GetGrupOpsiProduk(idToko, idProduk)

	message, _ := json.Marshal(dataGrupOpsiProduk)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}
