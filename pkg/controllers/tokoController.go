package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
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
	} else if _, ok := rj.GetIDKota(toko.Kota); !ok {
		http.Error(w, "Kota tidak ditemukan", http.StatusBadRequest)
		return
	}

	toko.IDOwner = user.IDCustomer
	toko.Foto = "toko.jpg"
	toko.SetKaryawan = false
	toko.CreatedAt = time.Now().Format("2006-01-02")

	_, err := toko.CreateToko()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"","domain":"` + toko.Slug + `"}`))
}

// UpdateToko is func
func (tc TokoController) UpdateToko(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]

	var toko models.Toko
	var rekening models.Rekening
	var jpt models.JasaPengirimanToko
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

	// Update main data toko
	err := toko.UpdateToko(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
	w.Write([]byte(`{"message":""}`))
}
