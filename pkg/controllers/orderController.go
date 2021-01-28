package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"gopkg.in/go-playground/validator.v9"
)

// OrderController is class
type OrderController struct{}

// GetOrder is get detail order customer
func (oc OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	var order models.Order

	dataOrder, err := order.GetOrder(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataOrder)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOrders is get all order customer
func (oc OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var order models.Order

	dataOrder := order.GetOrders(strconv.Itoa(user.IDCustomer))

	message, _ := json.Marshal(dataOrder)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOrdersToko is get all order customer
func (oc OrderController) GetOrdersToko(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idToko := vars["idToko"]
	var order models.Order

	dataOrder := order.GetOrdersToko(idToko)

	message, _ := json.Marshal(dataOrder)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOrdersEditor is get all order customer
func (oc OrderController) GetOrdersEditor(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)

	idToko := vars["idToko"]
	var order models.Order

	dataOrder := order.GetOrdersEditor(idToko, strconv.Itoa(user.IDCustomer))

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

	grupOpsi, _ := json.Marshal(order)

	// Cek batas min & max yg diperbolehkan grupOpsi
	var produk models.Produk
	dataProduk, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, "Gagal! Produk tidak ditemukan", http.StatusBadRequest)
		return
	}

	for _, vGrupOpsi := range dataProduk.GrupOpsi {
		totalOpsiGrup := gjson.Get(string(grupOpsi), "opsiOrder.#(idGrupOpsi=="+strconv.Itoa(vGrupOpsi.IDGrupOpsi)+")#").Array()
		if len(totalOpsiGrup) < vGrupOpsi.Min {
			http.Error(w, "Gagal! "+vGrupOpsi.NamaGrup+" kurang dari batas minimal", http.StatusBadRequest)
			return
		} else if len(totalOpsiGrup) > vGrupOpsi.Max {
			http.Error(w, "Gagal! "+vGrupOpsi.NamaGrup+" melebihi batas maksimal", http.StatusBadRequest)
			return
		}
	}

	dtProduk, _ := json.Marshal(dataProduk)

	// Input detail(namaGrup, opsi, harga, berat, perProduk) ke OpsiOrder
	for k, v := range order.OpsiOrder {
		dataGop := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+")#").Array()
		namaGrup := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+").namaGrup").String()
		spesificRequest := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+").spesificRequest").Bool()
		dataOpsiProduk := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+")#.opsi.#(idOpsi=="+strconv.Itoa(v.IDOpsi)+")").Array()
		if len(dataGop) == 0 {
			http.Error(w, "Grup Opsi tidak ditemukan.", http.StatusBadRequest)
			return
		} else if len(dataOpsiProduk) == 0 && spesificRequest == false {
			http.Error(w, "Spesific Request tidak diizinkan.", http.StatusBadRequest)
			return
		} else if len(dataOpsiProduk) == 0 && spesificRequest == true {
			order.OpsiOrder[k].NamaGrup = namaGrup
			order.OpsiOrder[k].Harga = 0
			order.OpsiOrder[k].Berat = 0
			order.OpsiOrder[k].PerProduk = false
		} else {
			opsi := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+").opsi.#(idOpsi=="+strconv.Itoa(v.IDOpsi)+").opsi").String()
			harga := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+").opsi.#(idOpsi=="+strconv.Itoa(v.IDOpsi)+").harga").Int()
			berat := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+").opsi.#(idOpsi=="+strconv.Itoa(v.IDOpsi)+").berat").Int()
			perProduk := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.IDGrupOpsi)+").opsi.#(idOpsi=="+strconv.Itoa(v.IDOpsi)+").perProduk").Bool()

			order.OpsiOrder[k].NamaGrup = namaGrup
			order.OpsiOrder[k].Opsi = opsi
			order.OpsiOrder[k].Harga = int(harga)
			order.OpsiOrder[k].Berat = int(berat)
			order.OpsiOrder[k].PerProduk = perProduk
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
	oc.HitungHargaBeratOpsi(&order)

	// get ongkir, estimasi, jenis kurir dan set harga produk
	var rj RajaOngkir
	var toko models.Toko
	dataToko, _ := toko.GetToko(idToko)

	order.Pengiriman.Berat = (dataProduk.Berat * order.Pcs) + order.TotalBeratOpsi
	asal, _ := rj.GetIDKota(dataToko.Kota)
	tujuan, _ := rj.GetIDKota(order.Pengiriman.Kota)
	ongkir, estimasi, kurir, ok := rj.GetOngkir(asal, tujuan, order.Pengiriman.KodeKurir, order.Pengiriman.Service, strconv.Itoa(order.Pengiriman.Berat))
	order.Pengiriman.Kurir = kurir

	// simpan harga produk
	if order.JenisPesanan == "soft copy" {
		order.HargaProduk = dataProduk.HargaSoftCopy
	} else if order.JenisPesanan == "cetak" && ok {
		order.HargaProduk = dataProduk.HargaCetak
		order.Pengiriman.Ongkir = ongkir
		order.Pengiriman.Estimasi = estimasi
	} else {
		http.Error(w, "Gagal! Terjadi kesalahan. Mohon periksa data pengiriman.", http.StatusBadRequest)
		return
	}

	// hitung total belanja
	order.Total = (order.HargaProduk * order.Pcs) + order.TotalHargaWajah + order.TotalHargaOpsi + order.TotalTambahanBiaya + order.Pengiriman.Ongkir

	// buat order
	idOrder, err := order.CreateOrder(idToko, idProduk)
	if err != nil {
		_ = order.DeleteOrder(strconv.Itoa(idOrder))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to admin & owner

	var karyawan models.Karyawan
	admins := karyawan.GetAdmins(idToko)

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataToko.IDOwner)
	notif.IDPenerima = append(notif.IDPenerima, admins...)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Permintaan pesanan baru"
	notif.Pesan = notif.Pengirim + " telah memesan produk " + order.NamaProduk + ". Pesanan " + strconv.Itoa(idOrder)
	notif.Link = "/order/" + strconv.Itoa(idOrder)
	notif.CreatedAt = order.CreatedAt
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Mohon tunggu konfirmasi kami. Terimakasih.","idOrder":"` + strconv.Itoa(idOrder) + `"}`))
}

// HitungHargaBeratOpsi is func
func (oc OrderController) HitungHargaBeratOpsi(o *models.Order) {
	for _, valueOpsi := range o.OpsiOrder {
		if valueOpsi.PerProduk == false {
			o.TotalHargaOpsi += valueOpsi.Harga
			o.TotalBeratOpsi += valueOpsi.Berat
		} else {
			o.TotalHargaOpsi += valueOpsi.Harga * o.Pcs
			o.TotalBeratOpsi += valueOpsi.Berat * o.Pcs
		}

	}
}

// KonfirmasiOrder is func
func (oc OrderController) KonfirmasiOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idOrder := vars["idOrder"]
	var order models.Order

	dataOrder, err := order.GetOrder(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataOrder.Tagihan = dataOrder.Total
	dataOrder.StatusPembayaran = "menunggu pembayaran"
	dataOrder.StatusPesanan = "diproses"

	err = dataOrder.KonfirmasiOrder(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var notif models.Notifikasi

	notif.IDPenerima = append(notif.IDPenerima, dataOrder.IDCustomer)
	notif.Pengirim = dataOrder.NamaToko
	notif.Judul = "Pesanan telah dikonfirmasi"
	notif.Pesan = "Pesanan " + idOrder + " sedang diproses. Silahkan melakukan pembayaran."
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!"}`))
}

// SetWaktuPengerjaan is func
func (oc OrderController) SetWaktuPengerjaan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var dataJSON map[string]interface{}
	json.NewDecoder(r.Body).Decode(&dataJSON)

	waktu := fmt.Sprintf("%v", dataJSON["waktuPengerjaan"])
	var order models.Order

	if err := validator.New().Var(waktu, "required"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.WaktuPengerjaan = waktu

	if err := order.SetWaktuPengerjaan(idOrder); err != nil {
		http.Error(w, "Gagal!", http.StatusBadRequest)
		return
	}
}

// UploadHasilProduksi is func
func (oc OrderController) UploadHasilProduksi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var dataJSON map[string]interface{}
	json.NewDecoder(r.Body).Decode(&dataJSON)

	hasil := fmt.Sprintf("%v", dataJSON["hasil"])
	var order models.Order

	if err := validator.New().Var(hasil, "required"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.Hasil = hasil
	order.StatusPesanan = "menunggu persetujuan pembeli"

	if err := order.UploadHasilProduksi(idOrder); err != nil {
		http.Error(w, "Gagal!", http.StatusBadRequest)
		return
	}

	dataOrder, _ := order.GetOrder(idOrder)

	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataOrder.IDCustomer)
	notif.Pengirim = dataOrder.NamaToko
	notif.Judul = "Hasil pesanan sudah keluar"
	notif.Pesan = "Hasil pesanan " + idOrder + " sudah keluar. Segera beri tanggapan ke penjual"
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!"}`))
}

// SetujuiHasilProduksi is func
func (oc OrderController) SetujuiHasilProduksi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var order models.Order

	order.StatusPesanan = "sudah disetujui pembeli"

	if err := order.SetujuiHasilProduksi(idOrder); err != nil {
		http.Error(w, "Gagal!", http.StatusBadRequest)
		return
	}

	dataOrder, _ := order.GetOrder(idOrder)

	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataOrder.Penangan.IDPenangan)
	notif.Pengirim = dataOrder.NamaToko
	notif.Judul = "Hasil pesanan sudah disetujui"
	notif.Pesan = "Hasil pesanan " + idOrder + " sudah disetujui pembeli."
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses!"}`))
}

// CancelOrder is func
func (oc OrderController) CancelOrder(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var order models.Order

	order.StatusPesanan = "dibatalkan"

	if err := order.CancelOrder(idOrder); err != nil {
		http.Error(w, "Gagal!", http.StatusBadRequest)
		return
	}

	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.StatusPesanan != "menunggu konfirmasi" {
		http.Error(w, "Gagal! Pesanan sedang diproses.", http.StatusBadRequest)
		return
	}

	var toko models.Toko
	dataToko, _ := toko.GetToko(strconv.Itoa(dataOrder.IDToko))

	var karyawan models.Karyawan
	admins := karyawan.GetAdmins(strconv.Itoa(dataOrder.IDToko))

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var notif models.Notifikasi
	notif.IDPenerima = append(notif.IDPenerima, dataToko.IDOwner)
	notif.IDPenerima = append(notif.IDPenerima, admins...)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Pesanan dibatalkan"
	notif.Pesan = "Pesanan " + idOrder + " dibatalkan oleh pembeli."
	notif.Link = "/order/" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Pesanan telah dibatalkan"}`))
}

// FinishOrder is func
func (oc OrderController) FinishOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	idToko, err := strconv.Atoi(vars["idToko"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	dataOrder.StatusPesanan = "selesai"

	var pembukuan models.Pembukuan
	pembukuan.IDToko = idToko
	pembukuan.Jenis = "pemasukan"
	pembukuan.Keterangan = "Pesanan " + idOrder + " telah selesai"
	pembukuan.Nominal = dataOrder.Dibayar - dataOrder.Pengiriman.Ongkir
	pembukuan.TglTransaksi = time.Now().Format("2006-01-02")

	err = order.FinishOrder(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = pembukuan.CreatePembukuan(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Pesanan telah dibatalkan"}`))
}
