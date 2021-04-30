package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// GrupOpsiController is class
type GrupOpsiController struct{}

// GetGrupOpsis is func
func (goc GrupOpsiController) GetGrupOpsis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var grupOpsi models.GrupOpsi

	dataGrupOpsi := grupOpsi.GetGrupOpsis(idToko)
	message, _ := json.Marshal(&dataGrupOpsi)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"grupopsi":` + string(message) + `}`))
}

// GetGrupOpsi is func
func (goc GrupOpsiController) GetGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idGrupOpsi := vars["idGrupOpsi"]
	idToko := vars["idToko"]
	var gop models.GrupOpsi

	dataGrupOpsi, err := gop.GetGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(&dataGrupOpsi)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateGrupOpsi is func
func (goc GrupOpsiController) CreateGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var grupOpsi models.GrupOpsi

	if err := json.NewDecoder(r.Body).Decode(&grupOpsi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hitung Jumlah Opsi
	totalOpsi := len(grupOpsi.GetOpsi())
	if grupOpsi.GetSpesificRequest() == true {
		totalOpsi++
	}

	if grupOpsi.GetRequired() == false && grupOpsi.GetMin() != 0 {
		http.Error(w, "Jika pilihan tidak wajib diisi, maka minimal pilihan harus 0", http.StatusBadRequest)
		return
	} else if grupOpsi.GetRequired() == true && grupOpsi.GetMin() < 1 {
		http.Error(w, "Jika pilihan wajib diisi, maka minimal pilihan setidaknya 1", http.StatusBadRequest)
		return
	} else if grupOpsi.GetMax() > totalOpsi {
		http.Error(w, "Maksimal jumlah memilih melebihi batas jumlah opsi", http.StatusBadRequest)
		return
	} else if grupOpsi.GetMin() > grupOpsi.GetMax() {
		http.Error(w, "Minimal jumlah memilih harus kurang dari samadengan maksimal jumlah memilih", http.StatusBadRequest)
		return
	} else if grupOpsi.GetHardCopy() == false && grupOpsi.GetSoftCopy() == false {
		http.Error(w, "Pilih minimal satu jenis pemesanan", http.StatusBadRequest)
		return
	}

	idGrupOpsi, err := grupOpsi.CreateGrupOpsi(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"GrupOpsi telah ditambahkan.","idGrupOpsi":"` + strconv.Itoa(idGrupOpsi) + `"}`))
}

// UpdateGrupOpsi is func
func (goc GrupOpsiController) UpdateGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGrupOpsi := vars["idGrupOpsi"]
	var grupOpsi models.GrupOpsi

	if err := json.NewDecoder(r.Body).Decode(&grupOpsi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hitung Jumlah Opsi
	totalOpsi := len(grupOpsi.GetOpsi())
	if grupOpsi.GetSpesificRequest() == true {
		totalOpsi++
	}

	if grupOpsi.GetRequired() == false && grupOpsi.GetMin() != 0 {
		http.Error(w, "Jika pilihan tidak wajib diisi, maka minimal pilihan harus 0", http.StatusBadRequest)
		return
	} else if grupOpsi.GetRequired() == true && grupOpsi.GetMin() < 1 {
		http.Error(w, "Jika pilihan wajib diisi, maka minimal pilihan setidaknya 1", http.StatusBadRequest)
		return
	} else if grupOpsi.GetMax() > totalOpsi {
		http.Error(w, "Maksimal jumlah memilih melebihi batas jumlah opsi", http.StatusBadRequest)
		return
	} else if grupOpsi.GetMin() > grupOpsi.GetMax() {
		http.Error(w, "Minimal jumlah memilih harus kurang dari samadengan maksimal memilih", http.StatusBadRequest)
		return
	}

	err := grupOpsi.UpdateGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range grupOpsi.GetOpsi() {
		_ = v.CreateUpdateOpsi(idGrupOpsi)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"GrupOpsi telah diperbarui."}`))
}

// DeleteGrupOpsi is func
func (goc GrupOpsiController) DeleteGrupOpsi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGrupOpsi := vars["idGrupOpsi"]
	var grupOpsi models.GrupOpsi

	err := grupOpsi.DeleteGrupOpsi(idToko, idGrupOpsi)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))
}
