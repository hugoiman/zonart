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

// PembayaranController is class
type PembayaranController struct{}

// CreatePembayaran is func
func (pc PembayaranController) CreatePembayaran(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var pembayaran models.Pembayaran

	if err := json.NewDecoder(r.Body).Decode(&pembayaran); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(pembayaran); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := pembayaran.CreatePembayaran(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send Notif to admin and owner
	var order models.Order
	var toko models.Toko
	var customer models.Customer

	dataOrder, _ := order.GetOrder(idOrder)
	dataToko, _ := toko.GetToko(strconv.Itoa(dataOrder.IDToko))
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var karyawan models.Karyawan
	admins := karyawan.GetAdmins(strconv.Itoa(dataOrder.IDToko))

	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataToko.IDOwner)
	notif.IDPenerima = append(notif.IDPenerima, admins...)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Pembayaran Masuk"
	notif.Pesan = notif.Pengirim + " telah melakukan pembayaran Rp " + strconv.Itoa(pembayaran.Nominal) + ". No invoice:" + idOrder
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = order.TglOrder
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pembayaran telah terkirim."}`))
}

// KonfirmasiPembayaran is func
func (pc PembayaranController) KonfirmasiPembayaran(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPembayaran := vars["idPembayaran"]
	idOrder := vars["idOrder"]
	var pembayaran models.Pembayaran

	if err := json.NewDecoder(r.Body).Decode(&pembayaran); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(pembayaran); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataPembayaran, err := pembayaran.GetPembayaran(idPembayaran, idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	dataOrder.Invoice.Tagihan -= dataPembayaran.Nominal
	dataOrder.Invoice.TotalBayar += dataPembayaran.Nominal

	if order.Invoice.Tagihan == 0 {
		order.Invoice.StatusPembayaran = "lunas"
	}

	pembayaran.Status = "diterima"
	pembayaran.Nominal = dataPembayaran.Nominal
	err = pembayaran.KonfirmasiPembayaran(idPembayaran, idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to customer
	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataOrder.IDCustomer)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = notif.Pengirim + " telah mengonfirmasi pembayaran anda. Inv: " + idOrder
	notif.Pesan = ""
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = order.TglOrder

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pembayaran telah dikonfirmasi."}`))
}
