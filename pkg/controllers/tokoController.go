package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// TokoController is class
type TokoController struct{}

// GetToko is func
func (tc TokoController) GetToko(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var toko models.Toko

	dataToko, err := toko.GetToko(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := json.Marshal(dataToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetTokos is func
func (tc TokoController) GetTokos(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var toko models.Toko

	dataToko := toko.GetMyToko(strconv.Itoa(user.IDCustomer))
	for _, v := range toko.GetTokoByEmploye(strconv.Itoa(user.IDCustomer)).Tokos {
		dataToko.Tokos = append(dataToko.Tokos, v)
	}

	message, err := json.Marshal(dataToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateToko is func
func (tc TokoController) CreateToko(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var toko models.Toko
	var rj RajaOngkir
	regexSlug := regexp.MustCompile(`^([a-z])([a-z0-9-]{1,48})([a-z0-9])$`)

	if err := json.NewDecoder(r.Body).Decode(&toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !regexSlug.MatchString(toko.Slug) {
		http.Error(w, "Domain hanya dapat mengandung huruf, angka atau strip(-) & terdiri 3-50 karakter.", http.StatusBadRequest)
		return
	}
	_, ok := rj.GetIDKota(toko.Kota)
	if !ok {
		http.Error(w, "Kota tidak ditemukan", http.StatusBadRequest)
		return
	}

	toko.Owner = user.IDCustomer
	toko.Foto = "https://res.cloudinary.com/dbddhr9rz/image/upload/v1612894274/zonart/toko/toko_jhecxf.png"
	toko.CreatedAt = time.Now().Format("2006-01-02")

	_, err := toko.CreateToko()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Toko berhasil terdaftar!", "domain":"` + toko.Slug + `"}`))
}

// UpdateToko is func
func (tc TokoController) UpdateToko(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	vars := mux.Vars(r)
	idToko := vars["idToko"]

	var toko models.Toko
	var rj RajaOngkir
	var rekening models.Rekening
	var jpt models.JasaPengirimanToko
	regexSlug := regexp.MustCompile(`^([a-z])([a-z0-9-]{1,48})([a-z0-9])$`)

	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !regexSlug.MatchString(toko.Slug) {
		http.Error(w, "Domain hanya dapat mengandung huruf, angka atau strip(-) & terdiri 3-50 karakter.", http.StatusBadRequest)
		return
	}

	_, ok := rj.GetIDKota(toko.Kota)
	if !ok {
		http.Error(w, "Kota tidak ditemukan", http.StatusBadRequest)
		return
	}

	_, _, existImage := r.FormFile("foto")
	var oldToko models.Toko
	var images []string
	var cloudinary Cloudinary
	if existImage != http.ErrMissingFile {
		maxSize := int64(1024 * 1024 * 2) // 2 MB
		destinationFolder := "zonart/toko"
		images, err := cloudinary.UploadImages(r, maxSize, destinationFolder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		oldToko, _ = toko.GetToko(idToko)
		toko.Foto = images[0]
	}

	// Update main data toko
	err := toko.UpdateToko(idToko)
	if err != nil {
		cloudinary.DeleteImages(images)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if existImage != http.ErrMissingFile && oldToko.Foto != "https://res.cloudinary.com/dbddhr9rz/image/upload/v1612894274/zonart/toko/toko_jhecxf.png" {
		oldImage := []string{oldToko.Foto}
		cloudinary.DeleteImages(oldImage)
	}

	// update data pengiriman toko
	for k := range toko.JasaPengirimanToko {
		jpt = toko.JasaPengirimanToko[k]
		_ = jpt.CreateUpdatePengirimanToko(idToko)
	}

	// update rekening toko
	for x := range toko.Rekening {
		rekening = toko.Rekening[x]
		_ = rekening.CreateUpdateRekening(idToko)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil diperbarui!."}`))
}
