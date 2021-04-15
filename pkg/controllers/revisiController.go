package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// RevisiController is class
type RevisiController struct{}

// CreateRevisi is func
func (rc RevisiController) CreateRevisi(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	var revisi models.Revisi

	if err := json.NewDecoder(r.Body).Decode(&revisi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(revisi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.Invoice.StatusPesanan != "diproses" {
		http.Error(w, "Status Pesanan tidak sedang dalam diproses", http.StatusBadRequest)
		return
	}

	revisi.CreatedAt = time.Now().Format("2006-01-02")

	err := revisi.CreateRevisi(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var karyawan models.Karyawan
	dataKaryawan, _ := karyawan.GetKaryawan(strconv.Itoa(dataOrder.IDToko), strconv.Itoa(dataOrder.Penangan.IDKaryawan))

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	// send notif to penangan
	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataKaryawan.IDCustomer)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Permintaan revisi pesanan #" + dataOrder.IDInvoice
	notif.Pesan = "Revisi pesanan #" + dataOrder.IDInvoice + " baru. Segera periksa pesanan."
	notif.Link = dataOrder.Invoice.SlugToko + "/pesanan/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Revisi telah terkirim."}`))
}
