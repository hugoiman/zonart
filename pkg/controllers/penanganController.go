package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// PenanganController is class
type PenanganController struct{}

// SetPenangan is func
func (pc PenanganController) SetPenangan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idOrder := vars["idOrder"]

	var penangan models.Penangan

	if err := json.NewDecoder(r.Body).Decode(&penangan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var karyawan models.Karyawan
	dataKaryawan, err := karyawan.GetKaryawan(idToko, strconv.Itoa(penangan.GetIDKaryawan()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if dataKaryawan.GetStatus() != "aktif" {
		http.Error(w, "Status karyawan tidak aktif.", http.StatusBadRequest)
		return
	} else if dataKaryawan.GetPosisi() != "editor" {
		http.Error(w, "Posisi karyawan bukanlah editor.", http.StatusBadRequest)
		return
	}

	err = penangan.SetPenangan(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)

	var toko models.Toko
	dataToko, _ := toko.GetToko(idToko)

	id, _ := karyawan.GetIDCustomerByKaryawan(strconv.Itoa(dataKaryawan.GetIDKaryawan()))
	idCustomer, _ := strconv.Atoi(id)

	// send notif to karyawan penangan
	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), idCustomer))
	notif.SetPengirim(dataToko.GetNamaToko())
	notif.SetJudul("Pengerjaan Pesanan")
	notif.SetPesan("Anda telah diberi tugas untuk mengerjakan pesanan #" + dataOrder.GetInvoice().GetIDInvoice())
	notif.SetLink(dataToko.GetSlug() + "/pesanan/" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan diteruskan ke ` + dataKaryawan.GetNamaKaryawan() + `"}`))
}
