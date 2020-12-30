package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// ProdukController is class
type ProdukController struct{}

// GetProduks is func
func (pc ProdukController) GetProduks(w http.ResponseWriter, r *http.Request) {
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
func (pc ProdukController) GetProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idProduk := vars["idProduk"]
	idToko := vars["idToko"]
	var produk models.Produk

	dataProduk, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataProduk)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateProduk is func
func (pc ProdukController) CreateProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var produk models.Produk

	if err := json.NewDecoder(r.Body).Decode(&produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if produk.Cetak == false && produk.SoftCopy == false {
		http.Error(w, "Gagal! Pilih setidaknya satu jenis pemesanan", http.StatusBadRequest)
		return
	}

	produk.Status = true

	idProduk, err := produk.CreateProduk(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Produk telah ditambahkan.","idProduk":"` + strconv.Itoa(idProduk) + `"}`))
}

// UpdateProduk is func
func (pc ProdukController) UpdateProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idProduk := vars["idProduk"]
	var produk models.Produk

	if err := json.NewDecoder(r.Body).Decode(&produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if produk.Cetak == false && produk.SoftCopy == false {
		http.Error(w, "Gagal! Pilih setidaknya satu jenis pemesanan", http.StatusBadRequest)
		return
	}

	err := produk.UpdateProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data produk berhasil diperbarui!"}`))
}

// DeleteProduk is func
func (pc ProdukController) DeleteProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idProduk := vars["idProduk"]
	var produk models.Produk

	err := produk.DeleteProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Data telah dihapus!"}`))
}
