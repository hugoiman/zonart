package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// PembayaranController is class
type PembayaranController struct{}

// CreatePembayaran is func
func (pc PembayaranController) CreatePembayaran(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	var pembayaran models.Pembayaran
	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&pembayaran); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, _, err := r.FormFile("bukti"); err == http.ErrMissingFile {
		http.Error(w, "Masukkan bukti pembayaran", http.StatusBadRequest)
		return
	}

	var cloudinary Cloudinary
	maxSize := int64(1024 * 1024 * 2) // 2 MB
	destinationFolder := "zonart/pembayaran"
	images, err := cloudinary.UploadImages(r, maxSize, destinationFolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pembayaran.SetBukti(images[0])

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.GetInvoice().GetStatusPesanan() != "diproses" {
		http.Error(w, "Status pesanan tidak sedang dalam proses", http.StatusBadRequest)
		return
	}

	pembayaran.SetCreatedAt(time.Now().Format("2006-01-02"))
	pembayaran.SetStatus("Menunggu Konfirmasi")

	err = pembayaran.CreatePembayaran(idOrder)
	if err != nil {
		cloudinary.DeleteImages(images)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var date time.Time
	date, _ = time.Parse("2006-01-02", pembayaran.GetCreatedAt())
	pembayaran.SetCreatedAt(date.Format("02 Jan 2006"))

	// send Notif to admin and owner
	var toko models.Toko
	var customer models.Customer

	idToko, _ := toko.GetIDTokoByOrder(idOrder)
	dataToko, _ := toko.GetToko(idToko)
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var karyawan models.Karyawan
	admins := karyawan.GetAdmins(idToko)

	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), dataToko.GetOwner()))
	notif.SetPenerima(append(notif.GetPenerima(), admins...))
	notif.SetPengirim(dataCustomer.GetNama())
	notif.SetJudul("Pembayaran Masuk")
	notif.SetPesan(notif.GetPengirim() + " telah melakukan pembayaran Rp " + strconv.Itoa(pembayaran.GetNominal()) + ". No invoice:" + idOrder)
	notif.SetLink(dataOrder.GetInvoice().GetSlugToko() + "/pesanan/" + idOrder)
	notif.SetCreatedAt(order.GetTglOrder())
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
	} else if dataPembayaran.GetStatus() == "Sudah Dikonfirmasi" {
		http.Error(w, "Pembayaran sudah dikonfirmasi sebelumnya.", http.StatusBadRequest)
		return
	}

	dataOrder.GetInvoice().SetTagihan(dataOrder.GetInvoice().GetTagihan() - dataPembayaran.GetNominal())
	dataOrder.GetInvoice().SetTotalBayar(dataOrder.GetInvoice().GetTotalBayar() + dataPembayaran.GetNominal())

	if dataOrder.GetInvoice().GetTagihan() <= 0 {
		dataOrder.GetInvoice().SetStatusPembayaran("lunas")
	}

	dataPembayaran.SetStatus("Sudah Dikonfirmasi")

	err = dataOrder.GetInvoice().UpdateInvoice(dataOrder.GetInvoice().GetIDInvoice())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = dataPembayaran.UpdatePembayaran(idPembayaran, idOrder)

	// send notif to customer
	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), dataOrder.GetPemesan()))
	notif.SetPengirim(dataOrder.GetInvoice().GetNamaToko())
	notif.SetJudul(notif.GetPengirim() + " telah mengonfirmasi pembayaran anda. Pesanan #" + dataOrder.GetInvoice().GetIDInvoice())
	notif.SetPesan("")
	notif.SetLink("/order?id=" + idOrder)
	notif.SetCreatedAt(order.GetTglOrder())

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pembayaran telah dikonfirmasi."}`))
}
