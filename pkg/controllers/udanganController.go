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

// UndanganController is class
type UndanganController struct{}

// GetUndangans is func
func (uc UndanganController) GetUndangans(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var undangan models.Undangan

	dataUndangan := undangan.GetUndangans(idToko)
	message, _ := json.Marshal(&dataUndangan)

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

	message, _ := json.Marshal(&dataUndangan)

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
	}

	dataCustomer, err := customer.GetCustomer(undangan.GetEmail())
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusBadRequest)
		return
	}

	idCustomer := strconv.Itoa(dataCustomer.GetIDCustomer())
	dataUndangan, _ := undangan.CheckUndangan(idToko, idCustomer)
	if dataUndangan.GetStatus() == "diterima" {
		http.Error(w, "User sudah pernah menjadi karyawan dari toko anda.", http.StatusBadRequest)
		return
	}

	dataToko, _ := toko.GetToko(idToko)

	if dataCustomer.GetIDCustomer() == dataToko.GetOwner() {
		http.Error(w, "Anda adalah pemilik dari toko ini.", http.StatusBadRequest)
		return
	}

	undangan.SetIDUndangan(dataUndangan.GetIDUndangan())
	undangan.SetStatus("menunggu")
	undangan.SetDate(time.Now().Format("2006-01-02"))

	idUndangan, err := undangan.UndangKaryawan(idToko, idCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to new karyawan
	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), dataCustomer.GetIDCustomer()))
	notif.SetPengirim(dataToko.GetNamaToko())
	notif.SetJudul("Undangan rekrut pegawai dari " + notif.GetPengirim())
	notif.SetPesan(notif.GetPengirim() + " telah mengundang anda sebagai " + undangan.GetPosisi() + ".")
	notif.SetLink("/undangan-rekrut/" + strconv.Itoa(idUndangan))
	notif.SetCreatedAt(undangan.GetDate())

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
	if err != nil || dataUndangan.GetStatus() != "menunggu" {
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
	notif.SetPenerima(append(notif.GetPenerima(), dataToko.GetOwner()))
	notif.SetPengirim(dataCustomer.GetNama())
	notif.SetJudul("Undangan Ditolak")
	notif.SetPesan(notif.GetPengirim() + " telah menolak undangan anda")
	notif.SetLink("")
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
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
	if err != nil || dataUndangan.GetStatus() != "menunggu" {
		http.Error(w, "Undangan tidak tersedia", http.StatusBadRequest)
		return
	}

	var toko models.Toko
	idToko, _ := toko.GetIDTokoByUndangan(idUndangan)

	karyawan.SetNamaKaryawan(dataCustomer.GetNama())
	karyawan.SetEmail(dataCustomer.GetEmail())
	karyawan.SetHP("")
	karyawan.SetAlamat("")
	karyawan.SetPosisi(dataUndangan.GetPosisi())
	karyawan.SetStatus("aktif")
	karyawan.SetBergabung(time.Now().Format("2006-01-02"))

	_ = karyawan.CreateKaryawan(idToko, idCustomer)

	_ = undangan.TerimaUndangan(idUndangan)

	dataToko, _ := toko.GetToko(idToko)

	// send notif to owner
	var notif models.Notifikasi
	notif.SetPenerima(append(notif.GetPenerima(), dataToko.GetOwner()))
	notif.SetPengirim(dataCustomer.GetNama())
	notif.SetJudul("Undangan Diterima")
	notif.SetPesan(notif.GetPengirim() + " telah menerima undangan anda")
	notif.SetLink("")
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Berhasil menerima undangan. Anda resmi menjadi karyawan di ` + dataToko.GetNamaToko() + `"}`))
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
	} else if dataUndangan.GetStatus() != "menunggu" {
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
