package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
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
	}

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.GetInvoice().GetStatusPesanan() != "diproses" {
		http.Error(w, "Status Pesanan tidak sedang dalam diproses", http.StatusBadRequest)
		return
	}

	revisi.SetCreatedAt(time.Now().Format("2006-01-02"))

	err := revisi.CreateRevisi(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var toko models.Toko
	idToko, _ := toko.GetIDTokoByOrder(idOrder)

	var karyawan models.Karyawan
	dataKaryawan, _ := karyawan.GetKaryawan(idToko, strconv.Itoa(dataOrder.GetPenangan().GetIDKaryawan()))

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	id, _ := karyawan.GetIDCustomerByKaryawan(strconv.Itoa(dataKaryawan.GetIDKaryawan()))
	idCustomer, _ := strconv.Atoi(id)

	// send notif to penangan
	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), idCustomer))
	notif.SetPengirim(dataCustomer.GetNama())
	notif.SetJudul("Permintaan revisi pesanan #" + dataOrder.GetInvoice().GetIDInvoice())
	notif.SetPesan("Revisi pesanan #" + dataOrder.GetInvoice().GetIDInvoice() + " baru. Segera periksa pesanan.")
	notif.SetLink(dataOrder.GetInvoice().GetSlugToko() + "/pesanan/" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Revisi telah terkirim."}`))
}
