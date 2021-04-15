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

// PembayaranController is class
type PembayaranController struct{}

// CreatePembayaran is func
func (pc PembayaranController) CreatePembayaran(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var cloudinary Cloudinary
	var pembayaran models.Pembayaran
	decode := json.NewDecoder(r.Body).Decode(&pembayaran)

	var images = []string{pembayaran.Bukti}
	if decode != nil {
		_ = cloudinary.DeleteImages(images)
		http.Error(w, decode.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(pembayaran); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.Invoice.StatusPesanan != "diproses" {
		http.Error(w, "Status pesanan tidak sedang dalam proses", http.StatusBadRequest)
		return
	}

	pembayaran.CreatedAt = time.Now().Format("2006-01-02")
	pembayaran.Status = "Menunggu Konfirmasi"

	err := pembayaran.CreatePembayaran(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var date time.Time
	date, _ = time.Parse("2006-01-02", pembayaran.CreatedAt)
	pembayaran.CreatedAt = date.Format("02 Jan 2006")

	// send Notif to admin and owner
	var toko models.Toko
	var customer models.Customer

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
	notif.Link = dataOrder.Invoice.SlugToko + "/pesanan/" + idOrder
	notif.CreatedAt = order.TglOrder
	notif.CreateNotifikasi()

	data, _ := json.Marshal(pembayaran)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pembayaran telah terkirim.","pembayaran":` + string(data) + `}`))
}

// KonfirmasiPembayaran is func
func (pc PembayaranController) KonfirmasiPembayaran(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPembayaran := vars["idPembayaran"]
	idOrder := vars["idOrder"]

	var pembayaran models.Pembayaran
	var order models.Order

	dataOrder, _ := order.GetOrder(idOrder)
	dataPembayaran, err := pembayaran.GetPembayaran(idPembayaran, idOrder)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if dataPembayaran.Status == "Sudah Dikonfirmasi" {
		http.Error(w, "Pembayaran sudah dikonfirmasi sebelumnya.", http.StatusBadRequest)
		return
	}

	dataOrder.Invoice.Tagihan -= dataPembayaran.Nominal
	dataOrder.Invoice.TotalBayar += dataPembayaran.Nominal

	if dataOrder.Invoice.Tagihan <= 0 {
		dataOrder.Invoice.StatusPembayaran = "lunas"
	}

	dataPembayaran.Status = "Sudah Dikonfirmasi"

	err = dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = dataPembayaran.UpdatePembayaran(idPembayaran, idOrder)

	// send notif to customer
	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataOrder.IDCustomer)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = notif.Pengirim + " telah mengonfirmasi pembayaran anda. Inv: " + idOrder
	notif.Pesan = ""
	notif.Link = "/order?id=" + idOrder
	notif.CreatedAt = order.TglOrder

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pembayaran telah dikonfirmasi."}`))
}
