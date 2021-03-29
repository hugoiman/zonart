package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// KaryawanController is class
type KaryawanController struct{}

// GetKaryawans is func
func (kc KaryawanController) GetKaryawans(w http.ResponseWriter, r *http.Request) {
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
func (kc KaryawanController) GetKaryawan(w http.ResponseWriter, r *http.Request) {
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

// GetKaryawanByIDCustomer is func
func (kc KaryawanController) GetKaryawanByIDCustomer(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var karyawan models.Karyawan

	dataKaryawan, err := karyawan.GetKaryawanByIDCustomer(idToko, strconv.Itoa(user.IDCustomer))
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
func (kc KaryawanController) UpdateKaryawan(w http.ResponseWriter, r *http.Request) {
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
