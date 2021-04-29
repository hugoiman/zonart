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

// HasilOrderController is class
type HasilOrderController struct{}

// AddHasilOrder is func
func (hoc HasilOrderController) AddHasilOrder(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var ho models.HasilOrder
	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&ho); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, _, err := r.FormFile("hasil"); err == http.ErrMissingFile {
		http.Error(w, "Mohon masukan gambar", http.StatusBadRequest)
		return
	}

	maxSize := int64(1024 * 1024 * 2) // 2 MB
	destinationFolder := "zonart/hasilOrder"
	var cloudinary Cloudinary
	images, err := cloudinary.UploadImages(r, maxSize, destinationFolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ho.SetHasil(images[0])

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)

	ho.SetCreatedAt(time.Now().Format("2006-01-02"))
	ho.SetStatus("menunggu persetujuan")
	if err := ho.AddHasilOrder(idOrder); err != nil {
		cloudinary.DeleteImages(images)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldImage := []string{dataOrder.GetHasilOrder().GetHasil()}
	cloudinary.DeleteImages(oldImage)

	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), dataOrder.GetPemesan()))
	notif.SetPengirim(dataOrder.GetInvoice().GetNamaToko())
	notif.SetJudul("Hasil pesanan sudah keluar")
	notif.SetPesan("Hasil pesanan " + dataOrder.GetInvoice().GetIDInvoice() + " sudah keluar. Segera beri tanggapan ke penjual.")
	notif.SetLink("/order?id=" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
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
	if dataOrder.GetInvoice().GetStatusPesanan() != "diproses" {
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
	notif.SetPenerima(append(notif.GetPenerima(), dataOrder.GetPenangan().GetIDPenangan()))
	notif.SetPengirim(dataCustomer.GetNama())
	notif.SetJudul("Hasil pesanan " + dataOrder.GetInvoice().GetIDInvoice() + " telah disetujui.")
	notif.SetPesan("")
	notif.SetLink("/pesanan/" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
	notif.CreateNotifikasi()

	var message = ""
	if dataOrder.GetJenisPesanan() == "cetak" {
		message = "Barang akan segera dikirim."
	}

	if dataOrder.GetInvoice().GetTagihan() > 0 {
		message += " Yuk segera selesaikan pembayaran kamu."
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Hasil telah disetujui. ` + message + ` . Terimakasih"}`))
}
