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
func (uc UndanganController) GetUndangan(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idUndangan := vars["idUndangan"]
	idCustomer := strconv.Itoa(user.IDCustomer)

	var undangan models.Undangan

	dataUndangan, err := undangan.GetUndangan(idUndangan, idToko, idCustomer)
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
	if dataUndangan.Status == "disetujui" {
		http.Error(w, "User sudah pernah menjadi karyawan dari toko anda.", http.StatusBadRequest)
		return
	}

	dataToko, _ := toko.GetToko(idToko)

	if dataCustomer.IDCustomer == dataToko.IDOwner {
		http.Error(w, "Gagal! Anda adalah pemilik dari toko ini.", http.StatusBadRequest)
		return
	}

	undangan.IDUndangan = dataUndangan.IDUndangan
	undangan.IDCustomer = dataCustomer.IDCustomer
	undangan.Status = "menunggu"
	undangan.Date = time.Now().Format("2006-01-02")

	idUndangan, err := undangan.UndangKaryawan(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to new karyawan
	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataCustomer.IDCustomer)
	notif.Pengirim = dataToko.NamaToko
	notif.Judul = notif.Pengirim + " telah mengundang anda sebagai " + undangan.Posisi + "."
	notif.Pesan = "Pesanan sedang diproses. Silahkan melakukan pembayaran."
	notif.Link = "/undangan/" + strconv.Itoa(idUndangan)
	notif.CreatedAt = undangan.Date

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Undangan telah terkirim. Mohon tunggu persetujuan dari user."}`))
}

// TolakUndangan is func
func (uc UndanganController) TolakUndangan(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)

	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idUndangan := vars["idUndangan"]
	idCustomer := strconv.Itoa(user.IDCustomer)

	var undangan models.Undangan

	dataUndangan, err := undangan.GetUndangan(idUndangan, idToko, idCustomer)
	if err != nil || dataUndangan.Status != "menunggu" {
		http.Error(w, "Undangan tidak tersedia", http.StatusBadRequest)
		return
	}

	_ = undangan.TolakUndangan(idUndangan, idToko, idCustomer)

	// send notif owner

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Undangan telah ditolak"}`))
}

// TerimaUndangan is func
func (uc UndanganController) TerimaUndangan(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)

	vars := mux.Vars(r)
	idUndangan := vars["idUndangan"]
	idToko := vars["idToko"]
	idCustomer := strconv.Itoa(user.IDCustomer)

	var undangan models.Undangan
	var karyawan models.Karyawan
	var customer models.Customer

	dataCustomer, _ := customer.GetCustomer(idCustomer)
	dataUndangan, err := undangan.GetUndangan(idUndangan, idToko, idCustomer)
	if err != nil || dataUndangan.Status != "menunggu" {
		http.Error(w, "Undangan tidak tersedia", http.StatusBadRequest)
		return
	}

	karyawan.IDCustomer = user.IDCustomer
	karyawan.IDToko = dataUndangan.IDToko
	karyawan.NamaKaryawan = dataCustomer.Nama
	karyawan.Email = dataCustomer.Email
	karyawan.Hp = ""
	karyawan.Alamat = ""
	karyawan.Posisi = dataUndangan.Posisi
	karyawan.Status = "aktif"
	karyawan.Bergabung = time.Now().Format("2006-01-02")

	_ = karyawan.CreateKaryawan()

	_ = undangan.TerimaUndangan(idUndangan, idToko, idCustomer)

	var toko models.Toko
	dataToko, err := toko.GetToko(idToko)

	// send notif to owner
	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataToko.IDOwner)
	notif.Pengirim = dataToko.NamaToko
	notif.Judul = notif.Pengirim + " telah menerima undangan anda"
	notif.Pesan = "Pesanan sedang diproses. Silahkan melakukan pembayaran."
	notif.Link = "/undangan/" + idUndangan
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!"}`))
}

// BatalkanUndangan is func
func (uc UndanganController) BatalkanUndangan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idUndangan := vars["idUndangan"]
	idCustomer := vars["idCustomer"]

	var undangan models.Undangan

	dataUndangan, err := undangan.GetUndangan(idUndangan, idToko, idCustomer)
	if err != nil {
		http.Error(w, "Undangan tidak ditemukan", http.StatusBadRequest)
		return
	} else if dataUndangan.Status != "menunggu" {
		http.Error(w, "Gagal! Undangan hanya dapat dibatalkan jika status undangan adalah MENUNGGU.", http.StatusBadRequest)
		return
	}

	err = undangan.BatalkanUndangan(idUndangan, idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif owner

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Undangan telah dibatalkan"}`))
}
