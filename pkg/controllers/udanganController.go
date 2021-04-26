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

// UndanganController is class
type UndanganController struct{}

// GetUndangans is func
func (uc UndanganController) GetUndangans(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var undangan models.Undangan

	dataUndangan := undangan.GetUndangans(idToko)
	message, _ := json.Marshal(dataUndangan)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetUndangan is func
func (uc UndanganController) GetUndanganCustomer(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idUndangan := vars["idUndangan"]

	var undangan models.Undangan

	dataUndangan, err := undangan.GetUndanganCustomer(idUndangan, strconv.Itoa(user.IDCustomer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataUndangan)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// UndangKaryawan is func
func (uc UndanganController) UndangKaryawan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]

	var customer models.Customer
	var undangan models.Undangan
	var toko models.Toko

	if err := json.NewDecoder(r.Body).Decode(&undangan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(undangan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataCustomer, err := customer.GetCustomer(undangan.Email)
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusBadRequest)
		return
	}

	idCustomer := strconv.Itoa(dataCustomer.IDCustomer)
	dataUndangan, _ := undangan.CheckUndangan(idToko, idCustomer)
	if dataUndangan.Status == "diterima" {
		http.Error(w, "User sudah pernah menjadi karyawan dari toko anda.", http.StatusBadRequest)
		return
	}

	dataToko, _ := toko.GetToko(idToko)

	if dataCustomer.IDCustomer == dataToko.Owner {
		http.Error(w, "Anda adalah pemilik dari toko ini.", http.StatusBadRequest)
		return
	}

	undangan.IDUndangan = dataUndangan.IDUndangan
	undangan.Status = "menunggu"
	undangan.Date = time.Now().Format("2006-01-02")

	idUndangan, err := undangan.UndangKaryawan(idToko, idCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to new karyawan
	var notif models.Notifikasi
	notif.Penerima = append(notif.Penerima, dataCustomer.IDCustomer)
	notif.Pengirim = dataToko.NamaToko
	notif.Judul = "Undangan rekrut pegawai dari " + notif.Pengirim
	notif.Pesan = notif.Pengirim + " telah mengundang anda sebagai " + undangan.Posisi + "."
	notif.Link = "/undangan-rekrut/" + strconv.Itoa(idUndangan)
	notif.CreatedAt = undangan.Date

	err = notif.CreateNotifikasi()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Undangan telah terkirim. Mohon tunggu persetujuan dari user."}`))
}

// TolakUndangan is func
func (uc UndanganController) TolakUndangan(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)

	vars := mux.Vars(r)
	idUndangan := vars["idUndangan"]
	idCustomer := strconv.Itoa(user.IDCustomer)

	var undangan models.Undangan
	var customer models.Customer

	dataCustomer, _ := customer.GetCustomer(idCustomer)
	dataUndangan, err := undangan.GetUndanganCustomer(idUndangan, idCustomer)
	if err != nil || dataUndangan.Status != "menunggu" {
		http.Error(w, "Undangan tidak tersedia", http.StatusBadRequest)
		return
	}

	_ = undangan.TolakUndangan(idUndangan)

	// send notif owner
	var toko models.Toko
	idToko, _ := toko.GetIDTokoByUndangan(idUndangan)
	dataToko, _ := toko.GetToko(idToko)

	// send notif to owner
	var notif models.Notifikasi
	notif.Penerima = append(notif.Penerima, dataToko.Owner)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Undangan Ditolak"
	notif.Pesan = notif.Pengirim + " telah menolak undangan anda"
	notif.Link = ""
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Undangan telah ditolak"}`))
}

// TerimaUndangan is func
func (uc UndanganController) TerimaUndangan(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)

	vars := mux.Vars(r)
	idUndangan := vars["idUndangan"]
	idCustomer := strconv.Itoa(user.IDCustomer)

	var undangan models.Undangan
	var karyawan models.Karyawan
	var customer models.Customer

	dataCustomer, _ := customer.GetCustomer(idCustomer)
	dataUndangan, err := undangan.GetUndanganCustomer(idUndangan, idCustomer)
	if err != nil || dataUndangan.Status != "menunggu" {
		http.Error(w, "Undangan tidak tersedia", http.StatusBadRequest)
		return
	}

	var toko models.Toko
	idToko, _ := toko.GetIDTokoByUndangan(idUndangan)

	karyawan.IDCustomer = user.IDCustomer
	karyawan.NamaKaryawan = dataCustomer.Nama
	karyawan.Email = dataCustomer.Email
	karyawan.Hp = ""
	karyawan.Alamat = ""
	karyawan.Posisi = dataUndangan.Posisi
	karyawan.Status = "aktif"
	karyawan.Bergabung = time.Now().Format("2006-01-02")

	_ = karyawan.CreateKaryawan(idToko)

	_ = undangan.TerimaUndangan(idUndangan)

	dataToko, _ := toko.GetToko(idToko)

	// send notif to owner
	var notif models.Notifikasi
	notif.Penerima = append(notif.Penerima, dataToko.Owner)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Undangan Diterima"
	notif.Pesan = notif.Pengirim + " telah menerima undangan anda"
	notif.Link = ""
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Berhasil menerima undangan. Anda menjadi karyawan dari ` + dataToko.NamaToko + `"}`))
}

// BatalkanUndangan is func
func (uc UndanganController) BatalkanUndangan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idUndangan := vars["idUndangan"]

	var undangan models.Undangan

	dataUndangan, err := undangan.GetUndangan(idUndangan)
	if err != nil {
		http.Error(w, "Undangan tidak ditemukan", http.StatusBadRequest)
		return
	} else if dataUndangan.Status != "menunggu" {
		http.Error(w, "Undangan hanya dapat dibatalkan jika status undangan adalah MENUNGGU.", http.StatusBadRequest)
		return
	}

	err = undangan.BatalkanUndangan(idUndangan, idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Undangan telah dibatalkan"}`))
}
