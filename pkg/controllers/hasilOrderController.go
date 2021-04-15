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

// HasilOrderController is class
type HasilOrderController struct{}

// AddHasilOrder is func
func (hoc HasilOrderController) AddHasilOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var ho models.HasilOrder
	if err := json.NewDecoder(r.Body).Decode(&ho); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(ho); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ho.CreatedAt = time.Now().Format("2006-01-02")
	ho.Status = "menunggu persetujuan"
	if err := ho.AddHasilOrder(idOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)

	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataOrder.IDCustomer)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = "Hasil pesanan sudah keluar"
	notif.Pesan = "Hasil pesanan " + dataOrder.IDInvoice + " sudah keluar. Segera beri tanggapan ke penjual."
	notif.Link = "/order?id=" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Hasil telah terkirim. Mohon tunggu persetujuan pembeli"}`))
}

// SetujuiHasilOrder is func
func (hoc HasilOrderController) SetujuiHasilOrder(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.Invoice.StatusPesanan != "diproses" {
		http.Error(w, "Status pesanan tidak sedang dalam proses.", http.StatusBadRequest)
		return
	}

	var ho models.HasilOrder
	if err := ho.SetujuiHasilOrder(idOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataOrder.Penangan.IDPenangan)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Hasil pesanan " + dataOrder.Invoice.IDInvoice + " telah disetujui."
	notif.Pesan = ""
	notif.Link = "/pesanan/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi()

	var message = ""
	if dataOrder.JenisPesanan == "cetak" {
		message = "Barang akan segera dikirim."
	}

	if dataOrder.Invoice.Tagihan > 0 {
		message += " Yuk segera selesaikan pembayaran kamu."
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Hasil telah disetujui. ` + message + ` . Terimakasih"}`))
}
