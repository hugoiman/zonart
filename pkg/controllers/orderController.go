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

// OrderController is class
type OrderController struct{}

// GetOrder is get detail order customer
func (oc OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	var order models.Order

	dataOrder, err := order.GetOrder(idOrder, strconv.Itoa(user.IDCustomer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataOrder)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOrderToko is get detail order toko
func (oc OrderController) GetOrderToko(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	idToko := vars["idToko"]
	var order models.Order

	dataOrder, err := order.GetOrderToko(idOrder, idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataOrder)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateOrder is func
func (oc OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)

	idToko := vars["idToko"]
	idProduk := vars["idProduk"]

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get total opsi yg dipilih berdasarkan grupOpsi
	var totalOpsiGrup = make(map[int]int)
	for _, v := range order.OpsiOrder {
		totalOpsiGrup[v.IDGrupOpsi]++
	}

	// Cek batas min & max yg diperbolehkan grupOpsi
	var produk models.Produk
	dataProduk, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, "Gagal! Produk tidak ditemukan", http.StatusBadRequest)
		return
	}

	for _, vGrupOpsi := range dataProduk.GrupOpsi {
		if totalOpsiGrup[vGrupOpsi.IDGrupOpsi] < vGrupOpsi.Min {
			http.Error(w, "Gagal! "+vGrupOpsi.NamaGrup+" kurang dari batas minimal", http.StatusBadRequest)
			return
		} else if totalOpsiGrup[vGrupOpsi.IDGrupOpsi] > vGrupOpsi.Max {
			http.Error(w, "Gagal! "+vGrupOpsi.NamaGrup+" melebihi batas maksimal", http.StatusBadRequest)
			return
		}
	}

	// simpan dataGrupOpsiProduk in var map
	type grupOpsiProduk models.GrupOpsi
	var dataGop = map[int]grupOpsiProduk{}
	for _, vGrupOpsiProduk := range dataProduk.GrupOpsi {
		dataGop[vGrupOpsiProduk.IDGrupOpsi] = grupOpsiProduk{
			NamaGrup:        vGrupOpsiProduk.NamaGrup,
			Required:        vGrupOpsiProduk.Required,
			Min:             vGrupOpsiProduk.Min,
			Max:             vGrupOpsiProduk.Max,
			SpesificRequest: vGrupOpsiProduk.SpesificRequest,
		}
	}

	// simpan dataOpsiProduk in var bertipe multidimensi map -> map[idGrupOpsi]map[idOpsi]Opsi
	type opsiProduk models.Opsi
	var dataOpsiProduk = map[int]map[int]opsiProduk{}
	for _, vGOP := range dataProduk.GrupOpsi {
		dataOpsiProduk[vGOP.IDGrupOpsi] = map[int]opsiProduk{}
		for _, vOpsi := range vGOP.Opsi {
			dataOpsiProduk[vGOP.IDGrupOpsi][vOpsi.IDOpsi] = opsiProduk{
				NamaGrup:  vGOP.NamaGrup,
				Opsi:      vOpsi.Opsi,
				Harga:     vOpsi.Harga,
				Berat:     vOpsi.Berat,
				PerProduk: vOpsi.PerProduk,
			}
		}
	}

	// Input detail(namaGrup, opsi, harga, berat, perProduk) ke OpsiOrder
	for k, v := range order.OpsiOrder {
		if _, isExist := dataGop[v.IDGrupOpsi]; !isExist {
			http.Error(w, "Grup Opsi tidak ditemukan.", http.StatusBadRequest)
			return
		} else if _, isExist := dataOpsiProduk[v.IDGrupOpsi][v.IDOpsi]; !isExist && dataGop[v.IDGrupOpsi].SpesificRequest == false {
			http.Error(w, "Spesific Request tidak diizinkan.", http.StatusBadRequest)
			return
		} else if _, isExist := dataOpsiProduk[v.IDGrupOpsi][v.IDOpsi]; !isExist && dataGop[v.IDGrupOpsi].SpesificRequest == true {
			order.OpsiOrder[k].NamaGrup = dataGop[v.IDGrupOpsi].NamaGrup
			order.OpsiOrder[k].Harga = 0
			order.OpsiOrder[k].Berat = 0
			order.OpsiOrder[k].PerProduk = false
			dataGop[v.IDGrupOpsi] = grupOpsiProduk{
				SpesificRequest: false,
			}
		} else {
			order.OpsiOrder[k].NamaGrup = dataGop[v.IDGrupOpsi].NamaGrup
			order.OpsiOrder[k].Opsi = dataOpsiProduk[v.IDGrupOpsi][v.IDOpsi].Opsi
			order.OpsiOrder[k].Harga = dataOpsiProduk[v.IDGrupOpsi][v.IDOpsi].Harga
			order.OpsiOrder[k].Berat = dataOpsiProduk[v.IDGrupOpsi][v.IDOpsi].Berat
			order.OpsiOrder[k].PerProduk = dataOpsiProduk[v.IDGrupOpsi][v.IDOpsi].PerProduk
		}
	}

	order.IDCustomer = user.IDCustomer
	order.NamaProduk = dataProduk.NamaProduk
	order.StatusPesanan = "menunggu konfirmasi"
	order.StatusPembayaran = "-"
	order.Dibayar = 0
	order.Tagihan = 0
	order.HargaWajah = dataProduk.HargaWajah
	order.CreatedAt = time.Now().Format("2006-01-02")

	// Total Harga Wajah
	order.TotalHargaWajah = order.TambahanWajah * order.HargaWajah

	// Total Harga Opsi
	for _, valueOpsi := range order.OpsiOrder {
		if valueOpsi.PerProduk == false {
			order.TotalHargaOpsi += valueOpsi.Harga
		} else {
			order.TotalHargaOpsi += valueOpsi.Harga * order.Pcs
		}

		order.TotalBeratOpsi += valueOpsi.Berat
	}

	order.Pengiriman.Berat = order.TotalBeratOpsi

	// get ongkir dan set harga produk
	var rj RajaOngkir
	var toko models.Toko
	dataToko, _ := toko.GetToko(idToko)

	asal, _ := rj.GetIDKota(dataToko.Kota)
	tujuan, _ := rj.GetIDKota(order.Pengiriman.Kota)
	ongkir, estimasi, ok := rj.GetOngkir(asal, tujuan, order.Pengiriman.KodeKurir, order.Pengiriman.Service, strconv.Itoa(order.Pengiriman.Berat))
	order.Pengiriman.Ongkir = ongkir
	order.Pengiriman.Estimasi = estimasi

	// simpan harga produk
	if order.JenisPesanan == "soft copy" {
		order.HargaProduk = dataProduk.HargaSoftCopy
	} else if order.JenisPesanan == "cetak" && ok {
		order.HargaProduk = dataProduk.HargaCetak
	} else if order.JenisPesanan == "cetak" && !ok {
		http.Error(w, "Gagal! Terjadi kesalahan. Mohon periksa data pengiriman.", http.StatusBadRequest)
		return
	}

	// hitung total belanja
	order.Total = (order.HargaProduk * order.Pcs) + order.TotalHargaWajah + order.TotalHargaOpsi + order.TotalTambahanBiaya + order.Pengiriman.Ongkir

	// buat order
	idOrder, err := order.CreateOrder(idToko, idProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add opsi order
	for _, vOpsiOrder := range order.OpsiOrder {
		err = vOpsiOrder.CreateOpsiOrder(strconv.Itoa(idOrder))
		if err != nil {
			_ = order.DeleteOrder(strconv.Itoa(idOrder))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Simpan data pengiriman
	if order.JenisPesanan == "cetak" {
		err = order.Pengiriman.CreatePengiriman(strconv.Itoa(idOrder))
		if err != nil {
			_ = order.DeleteOrder(strconv.Itoa(idOrder))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// send notif to admin/owner
	penerima := []int{}
	penerima = append(penerima, dataToko.IDOwner)

	var karyawan models.Karyawan
	dataKaryawan := karyawan.GetKaryawans(idToko)
	for _, vKaryawan := range dataKaryawan.Karyawans {
		if vKaryawan.Posisi == "admin" {
			penerima = append(penerima, vKaryawan.IDCustomer)
		}
	}

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var notif models.Notifikasi
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Pesanan Baru"
	notif.Pesan = notif.Pengirim + " telah memesan produk " + order.NamaProduk
	notif.Link = "/order/" + strconv.Itoa(idOrder)
	notif.CreatedAt = order.CreatedAt

	for _, vPenerima := range penerima {
		notif.IDPenerima = vPenerima
		_ = notif.CreateNotifikasi()
	}

	message, _ := json.Marshal(order)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOngkir is setter
// func GetOngkir(asal, tujuan, kurir string, berat int) (int, error) {
// 	if kurir == "cod" {
// 		return 15000, nil
// 	}
// 	return 10000, nil
// }

// KonfirmasiOrder,

// KonfirmasiOrder is func
// func KonfirmasiOrder(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	idToko := vars["idToko"]
// 	idOrder := vars["idOrder"]

// 	dataOrder, err := order.GetOrder(idOrder)

// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message":"Sukses!"}`))

// }
