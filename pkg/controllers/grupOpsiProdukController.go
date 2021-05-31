package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// GrupOpsiProdukController is class
type GrupOpsiProdukController struct{}

// GetGrupOpsiProduks is get all produk in a grup opsi
func (gopc GrupOpsiProdukController) GetGrupOpsiProduks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGrupOpsi := vars["idGrupOpsi"]

	var gop models.GrupOpsiProduk

	dataGOP := gop.GetGrupOpsiProduks(idToko, idGrupOpsi)
	message, _ := json.Marshal(&dataGOP)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"grupOpsiProduk":` + string(message) + `}`))
}

// GetGrupOpsiProduksByProduk is get all produk in a grup opsi
func (gopc GrupOpsiProdukController) GetGrupOpsiProduksByProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idProduk := vars["idProduk"]

	var gop models.GrupOpsiProduk

	dataGOP := gop.GetGrupOpsiProduksByProduk(idToko, idProduk)
	message, _ := json.Marshal(&dataGOP)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"grupOpsiProduk":` + string(message) + `}`))
}

// SambungGrupOpsikeProduk is func
func (gopc GrupOpsiProdukController) SambungGrupOpsikeProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idProduk := vars["idProduk"]
	idGrupOpsi := vars["idGrupOpsi"]

	var produk models.Produk
	var grupOpsi models.GrupOpsi
	var gop models.GrupOpsiProduk

	_, errP := produk.GetProduk(idToko, idProduk)        // produk not found
	_, errGo := grupOpsi.GetGrupOpsi(idToko, idGrupOpsi) // grup opsi tidak ditemukan
	if errP != nil || errGo != nil {
		http.Error(w, "produk atau grup opsi tidak ditemukan.", http.StatusBadRequest)
		return
	}

	isAny := gop.CheckSambunganGrupOpsi(idProduk, idGrupOpsi)
	if isAny == false {
		_ = gop.SambungGrupOpsikeProduk(idProduk, idGrupOpsi)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Berhasil disambungkan."}`))
}

// PutusGrupOpsidiProduk is func
func (gopc GrupOpsiProdukController) PutusGrupOpsidiProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idProduk := vars["idProduk"]
	idGrupOpsi := vars["idGrupOpsi"]

	var produk models.Produk
	var gop models.GrupOpsiProduk

	_, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, "Produk tidak ditemukan.", http.StatusBadRequest)
		return
	}

	err = gop.PutusGrupOpsidiProduk(idProduk, idGrupOpsi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sambungan telah terputus."}`))
}
