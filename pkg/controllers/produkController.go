package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// GetProduks is func
func GetProduks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var produk models.Produk

	dataProduk := produk.GetProduks(idToko)
	message, _ := json.Marshal(dataProduk)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetProduk is func
func GetProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idProduk := vars["idProduk"]
	idToko := vars["idToko"]
	var produk models.Produk
	var gop models.GrupOpsi

	dataProduk, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataGrupOpsiProduk := gop.GetGrupOpsiProduk(idToko, idProduk)
	dataProduk.GrupOpsi = dataGrupOpsiProduk

	message, _ := json.Marshal(dataProduk)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}
