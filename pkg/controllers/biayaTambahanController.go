package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// BiayaTambahanController is class
type BiayaTambahanController struct{}

// CreateBiayaTambahans is func
func (btc BiayaTambahanController) CreateBiayaTambahans(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	idToko := vars["idToko"]

	var tb models.BiayaTambahan
	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&tb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(tb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataOrder, err := order.GetOrderToko(idOrder, idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if dataOrder.StatusPesanan == "selesai" {
		http.Error(w, "Gagal! Pesanan telah selesai.", http.StatusBadRequest)
		return
	}

	err = tb.CreateBiayaTambahan(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	penerima := []int{}
	penerima = append(penerima, dataOrder.IDCustomer)

	// send notif to customer
	var notif models.Notifikasi
	notif.Pengirim = dataOrder.NamaToko
	notif.Judul = "Anda mempunyai biaya tambahan baru"
	notif.Pesan = notif.Pengirim + "telah menambahkan biaya lain pada pesanan Anda. No invoice: " + strconv.Itoa(dataOrder.IDOrder)
	notif.Link = "/order/" + strconv.Itoa(dataOrder.IDOrder)
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi(penerima)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Biaya tambahan telah ditambahkan."}`))
}

// DeleteBiayaTambahan is func
func (btc BiayaTambahanController) DeleteBiayaTambahan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idOrder := vars["idOrder"]
	idBiayaTambahan := vars["idBiayaTambahan"]

	var bt models.BiayaTambahan
	var order models.Order
	var newerOrder models.Order

	dataOrder, err := order.GetOrderToko(idOrder, idToko)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan3.", http.StatusBadRequest)
		return
	} else if dataOrder.StatusPesanan == "selesai" {
		http.Error(w, "Gagal! Pesanan telah selesai.", http.StatusBadRequest)
		return
	}

	dataBT, err := bt.GetBiayaTambahan(idBiayaTambahan, idOrder)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan2.", http.StatusBadRequest)
		return
	}

	err = dataBT.DeleteBiayaTambahan(idBiayaTambahan, idOrder)
	if err != nil {
		http.Error(w, "Gagal! Data tidak ditemukan1.", http.StatusBadRequest)
		return
	}

	dataNewerOrder, _ := newerOrder.GetOrderToko(idOrder, idToko)
	if dataNewerOrder.Tagihan <= 0 {
		dataNewerOrder.StatusPembayaran = "lunas"
		_ = dataNewerOrder.UpdateStatusOrder(idOrder)
	}

	penerima := []int{}
	penerima = append(penerima, dataOrder.IDCustomer)

	var notif models.Notifikasi
	notif.IDPenerima = dataNewerOrder.IDCustomer
	notif.Pengirim = dataNewerOrder.NamaToko
	notif.Judul = "Biaya tambahan telah dibatalkan"
	notif.Pesan = notif.Pengirim + "telah membatalkan biaya tambahan berupa " + dataBT.Item + "(Rp " + strconv.Itoa(dataBT.Nominal) + "). No invoice: " + strconv.Itoa(dataOrder.IDOrder)
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi(penerima)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Data telah dihapus!"}`))

}
