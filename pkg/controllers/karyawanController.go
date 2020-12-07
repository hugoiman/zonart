package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// GetKaryawans is func
func GetKaryawans(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var karyawan models.Karyawan

	dataKaryawan := karyawan.GetKaryawans(idToko)
	message, _ := json.Marshal(dataKaryawan)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetKaryawan is func
func GetKaryawan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idKaryawan := vars["idKaryawan"]
	idToko := vars["idToko"]
	var karyawan models.Karyawan

	dataKaryawan, err := karyawan.GetKaryawan(idToko, idKaryawan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataKaryawan)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// UpdateKaryawan is func
func UpdateKaryawan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idKaryawan := vars["idKaryawan"]
	var karyawan models.Karyawan

	if err := json.NewDecoder(r.Body).Decode(&karyawan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(karyawan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := karyawan.UpdateKaryawan(idToko, idKaryawan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data karyawan berhasil diperbarui!"}`))
}
