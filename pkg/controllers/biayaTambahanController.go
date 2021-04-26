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

// CreateBiayaTambahan is func
func (btc BiayaTambahanController) CreateBiayaTambahan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	idToko := vars["idToko"]

	var bt models.BiayaTambahan
	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&bt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(bt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.Invoice.StatusPesanan != "diproses" {
		http.Error(w, "Pesanan tidak sedang dalam proses pengerjaan.", http.StatusBadRequest)
		return
	}

	err := bt.CreateBiayaTambahan(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newOngkir = 0
	if dataOrder.JenisPesanan == "cetak" && dataOrder.Pengiriman.KodeKurir != "cod" {
		// Update Berat Ongkir
		var rj RajaOngkir
		var toko models.Toko
		dataToko, _ := toko.GetToko(idToko)
		newBerat := dataOrder.Pengiriman.Berat + bt.Berat

		asal, _ := rj.GetIDKota(dataToko.Kota)
		tujuan, _ := rj.GetIDKota(dataOrder.Pengiriman.Kota)
		newOngkir, _, _, _ = rj.GetOngkir(asal, tujuan, dataOrder.Pengiriman.KodeKurir, dataOrder.Pengiriman.Service, strconv.Itoa(newBerat))

		var pengiriman models.Pengiriman
		_ = pengiriman.UpdateBeratOngkir(idOrder, newBerat, newOngkir)
	}

	dataOrder.Invoice.TotalPembelian = dataOrder.Invoice.TotalPembelian - dataOrder.Pengiriman.Ongkir + bt.Total + newOngkir
	dataOrder.Invoice.Tagihan = dataOrder.Invoice.TotalPembelian - dataOrder.Invoice.TotalBayar
	dataOrder.Invoice.StatusPembayaran = "menunggu pembayaran"
	_ = dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice)

	// send notif to customer
	penerima := []int{}
	var notif models.Notifikasi
	notif.Penerima = append(penerima, dataOrder.Pemesan)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = "Terdapat Biaya Tambahan Baru"
	notif.Pesan = "Pesanan " + dataOrder.Invoice.IDInvoice + " mempunyai biaya tambahan baru."
	notif.Link = "/order?id=" + strconv.Itoa(dataOrder.IDOrder)
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Biaya tambahan telah terkirim."}`))
}

// DeleteBiayaTambahan is func
func (btc BiayaTambahanController) DeleteBiayaTambahan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idOrder := vars["idOrder"]
	idBiayaTambahan := vars["idBiayaTambahan"]

	var bt models.BiayaTambahan
	var order models.Order

	dataOrder, err := order.GetOrder(idOrder)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	} else if dataOrder.Invoice.StatusPesanan != "diproses" {
		http.Error(w, "Pesanan tidak sedang dalam proses pengerjaan.", http.StatusBadRequest)
		return
	}

	dataBT, err := bt.GetBiayaTambahan(idBiayaTambahan, idOrder)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	var newOngkir = 0
	if dataOrder.JenisPesanan == "cetak" && dataOrder.Pengiriman.KodeKurir != "cod" {
		// Update Berat Ongkir
		var rj RajaOngkir
		var toko models.Toko
		dataToko, _ := toko.GetToko(idToko)
		newBerat := dataOrder.Pengiriman.Berat - dataBT.Berat

		asal, _ := rj.GetIDKota(dataToko.Kota)
		tujuan, _ := rj.GetIDKota(dataOrder.Pengiriman.Kota)
		newOngkir, _, _, _ = rj.GetOngkir(asal, tujuan, dataOrder.Pengiriman.KodeKurir, dataOrder.Pengiriman.Service, strconv.Itoa(newBerat))

		var pengiriman models.Pengiriman
		_ = pengiriman.UpdateBeratOngkir(idOrder, newBerat, newOngkir)
	}

	dataOrder.Invoice.TotalPembelian = dataOrder.Invoice.TotalPembelian - dataOrder.Pengiriman.Ongkir - dataBT.Total + newOngkir
	dataOrder.Invoice.Tagihan = dataOrder.Invoice.TotalPembelian - dataOrder.Invoice.TotalBayar
	if dataOrder.Invoice.Tagihan <= 0 {
		dataOrder.Invoice.StatusPembayaran = "lunas"
	}
	_ = dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice)

	_ = bt.DeleteBiayaTambahan(idBiayaTambahan, idOrder)

	var notif models.Notifikasi
	notif.Penerima = append(notif.Penerima, dataOrder.Pemesan)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = "Pembatalan biaya tambahan pada pesanan " + dataOrder.Invoice.IDInvoice
	notif.Pesan = "Biaya tambahan berupa " + dataBT.Item + "(Rp " + strconv.Itoa(dataBT.Total) + ") telah dibatalkan."
	notif.Link = "/order?id=" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
