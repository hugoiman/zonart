package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// NotifikasiController is class
type NotifikasiController struct{}

// GetNotifikasis is func
func (nc NotifikasiController) GetNotifikasis(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	idPenerima := user.IDCustomer
	var notifikasi models.Notifikasi

	dataNotifikasi := notifikasi.GetNotifikasis(strconv.Itoa(idPenerima))
	message, _ := json.Marshal(dataNotifikasi)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// ReadNotifikasi is func
func (nc NotifikasiController) ReadNotifikasi(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)

	idNotifikasi := vars["idNotifikasi"]
	idPenerima := user.IDCustomer

	var notifikasi models.Notifikasi

	_ = notifikasi.ReadNotifikasi(idNotifikasi, strconv.Itoa(idPenerima))

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!"}`))
}
