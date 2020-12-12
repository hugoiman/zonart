package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
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

	dataGrupOpsi, err := gop.GetGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataGrupOpsi)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateGrupOpsi is func
func CreateGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var grupOpsi models.GrupOpsi

	if err := json.NewDecoder(r.Body).Decode(&grupOpsi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(grupOpsi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hitung Jumlah Opsi
	totalOpsi := len(grupOpsi.Opsi)
	if grupOpsi.SpesificRequest == true {
		totalOpsi++
	}

	if grupOpsi.Required == false && grupOpsi.Min != 0 {
		http.Error(w, "Gagal! Jika pilihan tidak wajib diisi, maka minimal pilihan harus 0", http.StatusBadRequest)
		return
	} else if grupOpsi.Required == true && grupOpsi.Min < 1 {
		http.Error(w, "Gagal! Jika pilihan wajib diisi, maka minimal pilihan setidaknya 1", http.StatusBadRequest)
		return
	} else if grupOpsi.Max > totalOpsi {
		http.Error(w, "Gagal! Maksimal jumlah memilih melebihi batas jumlah opsi", http.StatusBadRequest)
		return
	} else if grupOpsi.Min > grupOpsi.Max {
		http.Error(w, "Gagal! Minimal jumlah memilih harus kurang dari samadengan maksimal memilih", http.StatusBadRequest)
		return
	}

	idGrupOpsi, err := grupOpsi.CreateGrupOpsi(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range grupOpsi.Opsi {
		err = v.CreateUpdateOpsi(strconv.Itoa(idGrupOpsi))
		if err != nil {
			_ = grupOpsi.DeleteGrupOpsi(idToko, strconv.Itoa(idGrupOpsi))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! GrupOpsi telah ditambahkan.","idGrupOpsi":"` + strconv.Itoa(idGrupOpsi) + `"}`))
}

// UpdateGrupOpsi is func
func UpdateGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGrupOpsi := vars["idGrupOpsi"]
	var grupOpsi models.GrupOpsi

	if err := json.NewDecoder(r.Body).Decode(&grupOpsi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(grupOpsi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hitung Jumlah Opsi
	totalOpsi := len(grupOpsi.Opsi)
	if grupOpsi.SpesificRequest == true {
		totalOpsi++
	}

	if grupOpsi.Required == false && grupOpsi.Min != 0 {
		http.Error(w, "Gagal! Jika pilihan tidak wajib diisi, maka minimal pilihan harus 0", http.StatusBadRequest)
		return
	} else if grupOpsi.Required == true && grupOpsi.Min < 1 {
		http.Error(w, "Gagal! Jika pilihan wajib diisi, maka minimal pilihan setidaknya 1", http.StatusBadRequest)
		return
	} else if grupOpsi.Max > totalOpsi {
		http.Error(w, "Gagal! Maksimal jumlah memilih melebihi batas jumlah opsi", http.StatusBadRequest)
		return
	} else if grupOpsi.Min > grupOpsi.Max {
		http.Error(w, "Gagal! Minimal jumlah memilih harus kurang dari samadengan maksimal memilih", http.StatusBadRequest)
		return
	}

	err := grupOpsi.UpdateGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range grupOpsi.Opsi {
		_ = v.CreateUpdateOpsi(idGrupOpsi)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! GrupOpsi telah diperbarui."}`))
}

// DeleteGrupOpsi is func
func DeleteGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGrupOpsi := vars["idGrupOpsi"]
	var grupOpsi models.GrupOpsi

	err := grupOpsi.DeleteGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Data telah dihapus!"}`))
}
