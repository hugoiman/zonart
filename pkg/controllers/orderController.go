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

// GetOrderToko is get detail order customer
func (oc OrderController) GetOrderToko(w http.ResponseWriter, r *http.Request) {
	position := context.Get(r, "position").(map[string]interface{})
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	idToko := vars["idToko"]
	var order models.Order

	dataOrder, err := order.GetOrder(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if strconv.Itoa(dataOrder.IDToko) != idToko || (position["position"] == "editor" && strconv.Itoa(dataOrder.Penangan.IDKaryawan) != position["idKaryawan"].(string)) {
		http.Error(w, "Pesanan tidak ditemukan.", http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(dataOrder)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOrderByInvoice is get detail order customer
func (oc OrderController) GetOrderByInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idInvoice := vars["idInvoice"]
	var order models.Order

	dataOrder, err := order.GetOrderByInvoice(idInvoice)
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
	position := context.Get(r, "position").(map[string]interface{})
	vars := mux.Vars(r)

	idToko := vars["idToko"]
	var order models.Order
	var dataOrder models.Orders

	if position["position"].(string) == "owner" || position["position"].(string) == "admin" {
		dataOrder = order.GetOrdersToko(idToko)
	} else if position["position"] == "editor" {
		dataOrder = order.GetOrdersEditor(idToko, position["idKaryawan"].(string))
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

	grupOpsi, _ := json.Marshal(order)

	// Cek batas min & max yg diperbolehkan grupOpsi
	var produk models.Produk
	dataProduk, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusBadRequest)
		return
	}

	for _, vGrupOpsi := range dataProduk.GrupOpsi {
		totalOpsiGrup := gjson.Get(string(grupOpsi), "opsiOrder.#(idGrupOpsi=="+strconv.Itoa(vGrupOpsi.IDGrupOpsi)+")#").Array()
		if len(totalOpsiGrup) < vGrupOpsi.Min {
			http.Error(w, ""+vGrupOpsi.NamaGrup+" kurang dari batas minimal", http.StatusBadRequest)
			return
		} else if len(totalOpsiGrup) > vGrupOpsi.Max {
			http.Error(w, ""+vGrupOpsi.NamaGrup+" melebihi batas maksimal", http.StatusBadRequest)
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

	order.Invoice.IDInvoice = time.Now().Format("020106150405")
	order.IDCustomer = user.IDCustomer
	order.ProdukOrder.NamaProduk = dataProduk.NamaProduk
	order.ProdukOrder.BeratProduk = dataProduk.Berat
	order.Invoice.StatusPesanan = "menunggu konfirmasi"
	order.Invoice.StatusPembayaran = "-"
	order.Invoice.TotalBayar = 0
	order.Invoice.Tagihan = 0
	order.ProdukOrder.HargaSatuanWajah = dataProduk.HargaWajah
	order.TglOrder = time.Now().Format("2006-01-02")

	// Total Harga Wajah
	totalHargaWajah := order.TambahanWajah * order.ProdukOrder.HargaSatuanWajah

	// Total Harga Opsi
	totalHargaOpsi, totalBeratOpsi := oc.HitungHargaBeratOpsi(order)

	// get ongkir, estimasi, jenis kurir dan set harga produk
	var rj RajaOngkir
	var toko models.Toko
	dataToko, _ := toko.GetToko(idToko)

	// simpan harga produk
	if order.JenisPesanan == "soft copy" {
		order.ProdukOrder.HargaProduk = int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="soft copy").harga`).Int())
	} else if order.JenisPesanan == "cetak" && order.Pengiriman.KodeKurir == "cod" {
		order.ProdukOrder.HargaProduk = int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="cetak").harga`).Int())
		order.Pengiriman.Kurir = "COD (Cash On Delivery)"
	} else if order.JenisPesanan == "cetak" && order.Pengiriman.KodeKurir != "cod" {
		order.Pengiriman.Berat = (dataProduk.Berat * order.Pcs) + totalBeratOpsi
		asal, _ := rj.GetIDKota(dataToko.Kota)
		tujuan, _ := rj.GetIDKota(order.Pengiriman.Kota)
		ongkir, estimasi, kurir, err := rj.GetOngkir(asal, tujuan, order.Pengiriman.KodeKurir, order.Pengiriman.Service, strconv.Itoa(order.Pengiriman.Berat))

		if !err {
			http.Error(w, "Terjadi kesalahan. Mohon periksa data pengiriman.", http.StatusBadRequest)
			return
		}

		order.Pengiriman.Kurir = kurir

		order.ProdukOrder.HargaProduk = int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="cetak").harga`).Int())
		order.Pengiriman.Ongkir = ongkir
		order.Pengiriman.Estimasi = estimasi
	}

	// hitung total belanja
	order.Invoice.TotalPembelian = (order.ProdukOrder.HargaProduk * order.Pcs) + totalHargaWajah + totalHargaOpsi + order.Pengiriman.Ongkir

	// create invoice
	order.Invoice.IDCustomer = user.IDCustomer
	err = order.Invoice.CreateInvoice(idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create order
	idOrder, err := order.CreateOrder(idToko, idProduk)
	if err != nil {
		_ = order.Invoice.DeleteInvoice(order.Invoice.IDInvoice)
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
	notif.Pesan = notif.Pengirim + " telah memesan produk " + order.ProdukOrder.NamaProduk + ". Pesanan #" + order.Invoice.IDInvoice
	notif.Link = "/pesanan/" + strconv.Itoa(idOrder)
	notif.CreatedAt = order.TglOrder
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Mohon tunggu konfirmasi penjual. Terimakasih.","idOrder":"` + strconv.Itoa(idOrder) + `"}`))
}

// HitungHargaBeratOpsi is func
func (oc OrderController) HitungHargaBeratOpsi(o models.Order) (int, int) {
	var totalHargaOpsi, totalBeratOpsi int
	for _, valueOpsi := range o.OpsiOrder {
		if valueOpsi.PerProduk == false {
			totalHargaOpsi += valueOpsi.Harga
			totalBeratOpsi += valueOpsi.Berat
		} else {
			totalHargaOpsi += valueOpsi.Harga * o.Pcs
			totalBeratOpsi += valueOpsi.Berat * o.Pcs
		}
	}

	return totalHargaOpsi, totalBeratOpsi
}

// ProsesOrder is func
func (oc OrderController) ProsesOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idOrder := vars["idOrder"]
	var order models.Order

	dataOrder, err := order.GetOrder(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataOrder.Invoice.Tagihan = dataOrder.Invoice.TotalPembelian
	dataOrder.Invoice.StatusPembayaran = "menunggu pembayaran"
	dataOrder.Invoice.StatusPesanan = "diproses"

	err = dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var notif models.Notifikasi

	notif.IDPenerima = append(notif.IDPenerima, dataOrder.IDCustomer)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = "Pesanan telah dikonfirmasi"
	notif.Pesan = "Pesanan " + idOrder + " sedang diproses. Silahkan melakukan pembayaran."
	notif.Link = "/order?id=" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan diproses."}`))
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

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Waktu pengerjaan disimpan!"}`))
}

// CancelOrder is func
func (oc OrderController) CancelOrder(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	if dataOrder.Invoice.StatusPesanan != "menunggu konfirmasi" {
		http.Error(w, "Pesanan sedang diproses.", http.StatusBadRequest)
		return
	}

	dataOrder.Invoice.StatusPesanan = "dibatalkan"
	if err := dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice); err != nil {
		http.Error(w, "Gagal!", http.StatusBadRequest)
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
	notif.Pesan = "Pesanan #" + dataOrder.Invoice.IDInvoice + " dibatalkan oleh pembeli."
	notif.Link = "/order?id=" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan  #` + dataOrder.Invoice.IDInvoice + ` telah dibatalkan"}`))
}

// FinishOrder is func
func (oc OrderController) FinishOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	dataOrder.Invoice.StatusPesanan = "selesai"

	if dataOrder.Invoice.StatusPembayaran == "menunggu pembayaran" {
		dataOrder.Invoice.StatusPembayaran = "tidak lunas"
	}

	var pembukuan models.Pembukuan
	pembukuan.IDToko = dataOrder.IDToko
	pembukuan.Jenis = "pemasukan"
	pembukuan.Keterangan = "Pesanan #" + dataOrder.Invoice.IDInvoice + " telah selesai."
	pembukuan.Nominal = dataOrder.Invoice.TotalBayar - dataOrder.Pengiriman.Ongkir
	pembukuan.TglTransaksi = time.Now().Format("2006-01-02")

	err := dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = pembukuan.CreatePembukuan(strconv.Itoa(dataOrder.IDToko))

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan #` + dataOrder.Invoice.IDInvoice + ` telah diselesaikan"}`))
}
