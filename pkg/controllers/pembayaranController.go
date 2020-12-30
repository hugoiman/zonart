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

	dataOrder, _ := order.GetOrder(idOrder, strconv.Itoa(user.IDCustomer))
	dataToko, _ := toko.GetToko(strconv.Itoa(dataOrder.IDToko))
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	penerima := []int{}
	penerima = append(penerima, dataToko.IDOwner)

	var karyawan models.Karyawan
	dataKaryawan := karyawan.GetKaryawans(strconv.Itoa(dataToko.IDToko))
	for _, vKaryawan := range dataKaryawan.Karyawans {
		if vKaryawan.Posisi == "admin" {
			penerima = append(penerima, vKaryawan.IDCustomer)
		}
	}

	var notif models.Notifikasi
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Pembayaran Masuk"
	notif.Pesan = notif.Pengirim + " telah melakukan pembayaran Rp " + strconv.Itoa(pembayaran.Nominal) + ". No invoice:" + idOrder
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = order.CreatedAt

	for _, vPenerima := range penerima {
		notif.IDPenerima = vPenerima
		_ = notif.CreateNotifikasi()
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Pembayaran telah terkirim."}`))
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

	pembayaran.Status = "diterima"
	pembayaran.Nominal = dataPembayaran.Nominal
	err = pembayaran.KonfirmasiPembayaran(idPembayaran, idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to customer

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Pembayaran telah diterima."}`))
}
