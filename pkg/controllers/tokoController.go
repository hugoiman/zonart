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

// GetToko is func
func GetToko(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var toko models.Toko
	var pengirimanToko models.PengirimanToko
	var rekeningToko models.Rekening

	dataToko, err := toko.GetToko(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataPengirimanToko := pengirimanToko.GetPengirimanToko(id)
	dataToko.PengirimanToko = dataPengirimanToko

	dataRekeningToko := rekeningToko.GetRekening(id)
	dataToko.Rekening = dataRekeningToko

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
func CreateToko(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var toko models.Toko
	var pengirimanToko models.PengirimanToko
	regexSlug := regexp.MustCompile(`^([a-z])([a-z0-9-]{1,48})([a-z0-9])$`)

	if err := json.NewDecoder(r.Body).Decode(&toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !regexSlug.MatchString(toko.Slug) {
		http.Error(w, "Gagal! Domain hanya dapat mengandung huruf, angka atau strip(-) & terdiri 3-50 karakter.", http.StatusBadRequest)
		return
	}

	toko.IDOwner = user.IDCustomer
	toko.Foto = "toko.jpg"
	toko.SetKaryawan = false
	toko.CreatedAt = time.Now().Format("2006-01-02")

	idToko, err := toko.CreateToko()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = pengirimanToko.InitializePengirimanToko(idToko)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!","domain":"` + toko.Slug + `"}`))
}

// UpdateToko is func
func UpdateToko(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]

	var toko models.Toko
	var pengirimanToko models.PengirimanToko
	var rekening models.Rekening
	regexSlug := regexp.MustCompile(`^([a-z])([a-z0-9-]{1,48})([a-z0-9])$`)

	if err := json.NewDecoder(r.Body).Decode(&toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(toko); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !regexSlug.MatchString(toko.Slug) {
		http.Error(w, "Gagal! Domain hanya dapat mengandung huruf, angka atau strip(-) & terdiri 3-50 karakter.", http.StatusBadRequest)
		return
	}

	// Update main data toko
	err := toko.UpdateToko(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update data pengiriman toko
	pengirimanToko = toko.PengirimanToko
	_ = pengirimanToko.UpdatePengirimanToko(idToko)

	// update rekening toko
	for x := range toko.Rekening {
		rekening = toko.Rekening[x]
		_ = rekening.CreateUpdateRekening(idToko)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!"}`))
}
