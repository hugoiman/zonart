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

// SetPenangan is func
func SetPenangan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idOrder := vars["idOrder"]

	var penangan models.Penangan
	var karyawan models.Karyawan
	var toko models.Toko

	if err := json.NewDecoder(r.Body).Decode(&penangan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(penangan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataKaryawan, err := karyawan.GetKaryawan(idToko, strconv.Itoa(penangan.IDKaryawan))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if dataKaryawan.Status != "aktif" {
		http.Error(w, "Gagal! Status karyawan tidak aktif.", http.StatusBadRequest)
		return
	} else if dataKaryawan.Posisi != "editor" {
		http.Error(w, "Gagal! Posisi karyawan bukanlah editor.", http.StatusBadRequest)
		return
	}

	err = penangan.SetPenangan(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataToko, _ := toko.GetToko(idToko)

	// send notif to karyawan penangan
	var notif models.Notifikasi
	notif.IDPenerima = dataKaryawan.IDCustomer
	notif.Pengirim = dataToko.NamaToko
	notif.Judul = "Pengerjaan Pesanan"
	notif.Pesan = notif.Pengirim + "telah menugaskan Anda untuk mengerjakan pesanan dengan no invoice: " + idOrder
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")
	_ = notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!"}`))
}
