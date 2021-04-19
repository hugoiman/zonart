package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
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
	if err != nil || dataProduk.Status == "dihapus" {
		http.Error(w, "Produk tidak ditemukan.", http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataProduk)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateProduk is func
func (pc ProdukController) CreateProduk(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var produk models.Produk

	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, _, err := r.FormFile("gambar"); err == http.ErrMissingFile {
		http.Error(w, "Silahkan masukan foto produk", http.StatusBadRequest)
		return
	}

	maxSize := int64(1024 * 1024 * 2) // 2 MB
	destinationFolder := "zonart/produk"
	var cloudinary Cloudinary
	images, err := cloudinary.UploadImages(r, maxSize, destinationFolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	produk.Gambar = images[0]
	produk.Slug = slug.Make(produk.NamaProduk)

	_, err = produk.CreateProduk(idToko)
	if err != nil {
		cloudinary.DeleteImages(images)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Produk telah ditambahkan.","produk":"` + produk.Slug + `"}`))
}

// UpdateProduk is func
func (pc ProdukController) UpdateProduk(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idProduk := vars["idProduk"]
	var produk models.Produk

	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(produk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _, existImage := r.FormFile("gambar")

	var oldProduk models.Produk
	var images []string
	var cloudinary Cloudinary
	if existImage != http.ErrMissingFile {
		maxSize := int64(1024 * 1024 * 2) // 2 MB
		destinationFolder := "zonart/produk"
		images, err := cloudinary.UploadImages(r, maxSize, destinationFolder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		oldProduk, _ = produk.GetProduk(idToko, idProduk)
		produk.Gambar = images[0]
	}

	produk.Slug = slug.Make(produk.NamaProduk)
	err := produk.UpdateProduk(idToko, idProduk)
	if err != nil {
		cloudinary.DeleteImages(images)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if existImage != http.ErrMissingFile {
		oldImage := []string{oldProduk.Gambar}
		cloudinary.DeleteImages(oldImage)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data produk berhasil diperbarui!","produk":"` + produk.Slug + `"}`))
}

// DeleteProduk is func
func (pc ProdukController) DeleteProduk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idProduk := vars["idProduk"]
	var produk models.Produk

	err := produk.DeleteProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))
}
