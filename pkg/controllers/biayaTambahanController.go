package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
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
	}

	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.GetInvoice().GetStatusPesanan() != "diproses" {
		http.Error(w, "Pesanan tidak sedang dalam proses pengerjaan.", http.StatusBadRequest)
		return
	}

	err := bt.CreateBiayaTambahan(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newOngkir = 0
	if dataOrder.GetJenisPesanan() == "cetak" && dataOrder.GetPengiriman().GetKodeKurir() != "cod" {
		// Update Berat Ongkir
		var rj RajaOngkir
		var toko models.Toko
		dataToko, _ := toko.GetToko(idToko)
		newBerat := dataOrder.GetPengiriman().GetBerat() + bt.GetBerat()
		newOngkir, _, _, _ = rj.GetOngkir(dataToko.GetKota(), dataOrder.GetPengiriman().GetKota(), dataOrder.GetPengiriman().GetKodeKurir(), dataOrder.GetPengiriman().GetService(), newBerat)

		var pengiriman models.Pengiriman
		_ = pengiriman.UpdateBeratOngkir(idOrder, newBerat, newOngkir)
	}

	dataOrder.GetInvoice().SetTotalPembelian(dataOrder.GetInvoice().GetTotalPembelian() - dataOrder.GetPengiriman().GetOngkir() + bt.GetTotal() + newOngkir)
	dataOrder.GetInvoice().SetTagihan(dataOrder.GetInvoice().GetTotalPembelian() - dataOrder.GetInvoice().GetTotalBayar())
	dataOrder.GetInvoice().SetStatusPembayaran("menunggu pembayaran")
	_ = dataOrder.GetInvoice().UpdateInvoice(dataOrder.GetInvoice().GetIDInvoice())

	// send notif to customer
	penerima := []int{}
	var notif models.Notifikasi
	notif.SetPenerima(append(penerima, dataOrder.GetPemesan()))
	notif.SetPengirim(dataOrder.GetInvoice().GetNamaToko())
	notif.SetJudul("Terdapat Biaya Tambahan Baru")
	notif.SetPesan("Pesanan " + dataOrder.GetInvoice().GetIDInvoice() + " mempunyai biaya tambahan baru.")
	notif.SetLink("/order?id=" + strconv.Itoa(dataOrder.GetIDOrder()))
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
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
	} else if dataOrder.GetInvoice().GetStatusPesanan() != "diproses" {
		http.Error(w, "Pesanan tidak sedang dalam proses pengerjaan.", http.StatusBadRequest)
		return
	}

	dataBT, err := bt.GetBiayaTambahan(idBiayaTambahan, idOrder)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	var newOngkir = 0
	if dataOrder.GetJenisPesanan() == "cetak" && dataOrder.GetPengiriman().GetKodeKurir() != "cod" {
		// Update Berat Ongkir
		var rj RajaOngkir
		var toko models.Toko
		dataToko, _ := toko.GetToko(idToko)
		newBerat := dataOrder.GetPengiriman().GetBerat() - dataBT.GetBerat()

		newOngkir, _, _, _ = rj.GetOngkir(dataToko.GetKota(), dataOrder.GetPengiriman().GetKota(), dataOrder.GetPengiriman().GetKodeKurir(), dataOrder.GetPengiriman().GetService(), newBerat)

		var pengiriman models.Pengiriman
		_ = pengiriman.UpdateBeratOngkir(idOrder, newBerat, newOngkir)
	}

	dataOrder.GetInvoice().SetTotalPembelian(dataOrder.GetInvoice().GetTotalPembelian() - dataOrder.GetPengiriman().GetOngkir() - dataBT.GetTotal() + newOngkir)
	dataOrder.GetInvoice().SetTagihan(dataOrder.GetInvoice().GetTotalPembelian() - dataOrder.GetInvoice().GetTotalBayar())
	if dataOrder.GetInvoice().GetTagihan() <= 0 {
		dataOrder.GetInvoice().SetStatusPembayaran("lunas")
	}
	_ = dataOrder.GetInvoice().UpdateInvoice(dataOrder.GetInvoice().GetIDInvoice())

	_ = bt.DeleteBiayaTambahan(idBiayaTambahan, idOrder)

	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), dataOrder.GetPemesan()))
	notif.SetPengirim(dataOrder.GetInvoice().GetNamaToko())
	notif.SetJudul("Pembatalan biaya tambahan pada pesanan " + dataOrder.GetInvoice().GetIDInvoice())
	notif.SetPesan("Biaya tambahan berupa " + dataBT.GetItem() + "(Rp " + strconv.Itoa(dataBT.GetTotal()) + ") telah dibatalkan.")
	notif.SetLink("/order?id=" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
