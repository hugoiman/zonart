package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

	dataOrder, err := order.GetOrderToko(idOrder, idToko)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if position["position"] == "editor" && strconv.Itoa(dataOrder.Penangan.IDKaryawan) != position["idKaryawan"].(string) {
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
	payload := r.FormValue("payload")
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)

	idToko := vars["idToko"]
	idProduk := vars["idProduk"]

	var order models.Order
	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, _, err := r.FormFile("fileOrder"); err == http.ErrMissingFile {
		http.Error(w, "Silahkan masukan foto wajah", http.StatusBadRequest)
		return
	}

	// Cek batas min & max yg diperbolehkan grupOpsi
	var produk models.Produk
	dataProduk, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusBadRequest)
		return
	}

	err = oc.validateGrupOpsiOrder(order, dataProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dtProduk, _ := json.Marshal(dataProduk)

	// Input detail(namaGrup, opsi, harga, berat, perProduk) ke OpsiOrder
	err = oc.setOpsiOrder(&order, dtProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.Invoice.IDInvoice = time.Now().Format("020106150405")
	order.Pemesan = user.IDCustomer
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
	totalHargaOpsi, totalBeratOpsi := oc.hitungHargaBeratOpsi(order)

	// get ongkir, estimasi, jenis kurir dan set harga produk
	var toko models.Toko
	dataToko, _ := toko.GetToko(idToko)

	// set harga produk & pengiriman
	err = oc.setPengiriman(&order, dataProduk, dataToko, dtProduk, totalBeratOpsi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// hitung total belanja
	order.Invoice.TotalPembelian = (order.ProdukOrder.HargaProduk * order.Pcs) + totalHargaWajah + totalHargaOpsi + order.Pengiriman.Ongkir

	// create invoice
	err = order.Invoice.CreateInvoice(idToko, strconv.Itoa(user.IDCustomer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// upload fileOrder to cloud
	maxSize := int64(1024 * 1024 * 10) // 10 MB
	destinationFolder := "zonart/order"
	var cloudinary Cloudinary
	images, err := cloudinary.UploadImages(r, maxSize, destinationFolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fileOrder models.FileOrder
	for _, v := range images {
		fileOrder.Foto = v
		order.FileOrder = append(order.FileOrder, fileOrder)
	}

	// create order
	idOrder, err := order.CreateOrder(idToko, idProduk)
	if err != nil {
		_ = order.Invoice.DeleteInvoice(order.Invoice.IDInvoice)
		_ = cloudinary.DeleteImages(images)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to admin & owner
	var karyawan models.Karyawan
	admins := karyawan.GetAdmins(idToko)

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var notif models.Notifikasi
	notif.Penerima = append(notif.Penerima, dataToko.Owner)
	notif.Penerima = append(notif.Penerima, admins...)
	notif.Pengirim = dataCustomer.Nama
	notif.Judul = "Permintaan pesanan baru"
	notif.Pesan = notif.Pengirim + " telah memesan produk " + order.ProdukOrder.NamaProduk + ". Pesanan #" + order.Invoice.IDInvoice
	notif.Link = dataToko.Slug + "/pesanan/" + strconv.Itoa(idOrder)
	notif.CreatedAt = order.TglOrder
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Mohon tunggu konfirmasi penjual. Terimakasih.","idOrder":"` + strconv.Itoa(idOrder) + `"}`))
}

func (oc OrderController) validateGrupOpsiOrder(order models.Order, dataProduk models.Produk) error {
	grupOpsi, _ := json.Marshal(order)
	for _, vGrupOpsi := range dataProduk.GrupOpsi {
		totalOpsiGrup := gjson.Get(string(grupOpsi), "opsiOrder.#(idGrupOpsi=="+strconv.Itoa(vGrupOpsi.IDGrupOpsi)+")#").Array()
		if len(totalOpsiGrup) < vGrupOpsi.Min && ((order.JenisPesanan == "cetak" && vGrupOpsi.HardCopy) || (order.JenisPesanan == "soft copy" && vGrupOpsi.SoftCopy)) {
			return errors.New("Opsi " + vGrupOpsi.NamaGrup + " kurang dari batas minimal")
		} else if len(totalOpsiGrup) > vGrupOpsi.Max && ((order.JenisPesanan == "cetak" && vGrupOpsi.HardCopy) || (order.JenisPesanan == "soft copy" && vGrupOpsi.SoftCopy)) {
			return errors.New("Opsi " + vGrupOpsi.NamaGrup + " melebihi batas maksimal")
		} else if ((order.JenisPesanan == "cetak" && !vGrupOpsi.HardCopy) || (order.JenisPesanan == "soft copy" && !vGrupOpsi.SoftCopy)) && len(totalOpsiGrup) > 0 {
			return errors.New(vGrupOpsi.NamaGrup + " dapat diisi jika jenis pesanan selain " + order.JenisPesanan)
		}
	}
	return nil
}

func (oc OrderController) setOpsiOrder(order *models.Order, dtProduk []byte) error {
	for k, v := range order.OpsiOrder {
		dataGop := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+")#").Array()
		namaGrup := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+").namaGrup").String()
		spesificRequest := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+").spesificRequest").Bool()
		dataOpsiProduk := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+")#.opsi.#(idOpsi=="+v.Opsi+")").Array()
		if len(dataGop) == 0 {
			return errors.New("Grup Opsi tidak ditemukan.")
		} else if len(dataOpsiProduk) == 0 && spesificRequest == false {
			return errors.New("Spesific Request tidak diizinkan.")
		} else if len(dataOpsiProduk) == 0 && spesificRequest == true {
			order.OpsiOrder[k].NamaGrup = namaGrup
			order.OpsiOrder[k].Harga = 0
			order.OpsiOrder[k].Berat = 0
			order.OpsiOrder[k].PerProduk = false
		} else {
			opsi := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+").opsi.#(idOpsi=="+v.Opsi+").opsi").String()
			harga := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+").opsi.#(idOpsi=="+v.Opsi+").harga").Int()
			berat := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+").opsi.#(idOpsi=="+v.Opsi+").berat").Int()
			perProduk := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+v.NamaGrup+").opsi.#(idOpsi=="+v.Opsi+").perProduk").Bool()

			order.OpsiOrder[k].NamaGrup = namaGrup
			order.OpsiOrder[k].Opsi = opsi
			order.OpsiOrder[k].Harga = int(harga)
			order.OpsiOrder[k].Berat = int(berat)
			order.OpsiOrder[k].PerProduk = perProduk
		}
	}
	return nil
}

func (oc OrderController) setPengiriman(order *models.Order, dataProduk models.Produk, dataToko models.Toko, dtProduk []byte, totalBeratOpsi int) error {
	if order.JenisPesanan == "soft copy" {
		order.ProdukOrder.HargaProduk = int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="soft copy").harga`).Int())
		order.Pengiriman.Kota = dataToko.Kota
		order.Pengiriman.Kurir = ""
		order.Pengiriman.Alamat = ""
		order.Pengiriman.Label = ""
		order.Pengiriman.Service = ""
	} else if order.JenisPesanan == "cetak" && order.Pengiriman.KodeKurir == "cod" {
		order.ProdukOrder.HargaProduk = int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="cetak").harga`).Int())
		order.Pengiriman.Kota = dataToko.Kota
		order.Pengiriman.Kurir = "COD (Cash On Delivery)"
		order.Pengiriman.Alamat = ""
		order.Pengiriman.Label = ""
		order.Pengiriman.Service = ""
	} else if order.JenisPesanan == "cetak" && order.Pengiriman.KodeKurir != "cod" {
		var rj RajaOngkir
		order.Pengiriman.Berat = (dataProduk.Berat * order.Pcs) + totalBeratOpsi
		asal, _ := rj.GetIDKota(dataToko.Kota)
		tujuan, _ := rj.GetIDKota(order.Pengiriman.Kota)
		ongkir, estimasi, kurir, err := rj.GetOngkir(asal, tujuan, order.Pengiriman.KodeKurir, order.Pengiriman.Service, strconv.Itoa(order.Pengiriman.Berat))
		if !err {
			return errors.New("Terjadi kesalahan. Mohon periksa data pengiriman.")
		}

		order.Pengiriman.Kurir = kurir
		order.ProdukOrder.HargaProduk = int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="cetak").harga`).Int())
		order.Pengiriman.Ongkir = ongkir
		order.Pengiriman.Estimasi = estimasi
	}
	return nil
}

// hitungHargaBeratOpsi is func
func (oc OrderController) hitungHargaBeratOpsi(o models.Order) (int, int) {
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

	notif.Penerima = append(notif.Penerima, dataOrder.Pemesan)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = "Pesanan telah dikonfirmasi"
	notif.Pesan = "Pesanan #" + dataOrder.Invoice.IDInvoice + " sedang diproses. Silahkan melakukan pembayaran."
	notif.Link = "/order?id=" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan diproses."}`))
}

// TolakOrder is func
func (oc OrderController) TolakOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	var order models.Order

	data := struct {
		Keterangan string `json:"keterangan" validate:"required,max=100"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataOrder, err := order.GetOrder(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataOrder.Invoice.StatusPesanan = "ditolak"

	err = dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var notif models.Notifikasi

	notif.Penerima = append(notif.Penerima, dataOrder.Pemesan)
	notif.Pengirim = dataOrder.Invoice.NamaToko
	notif.Judul = "Pesanan ditolak"
	notif.Pesan = "Pesanan #" + dataOrder.Invoice.IDInvoice + " ditolak. Keterangan: " + data.Keterangan
	notif.Link = "/order?id=" + idOrder
	notif.CreatedAt = time.Now().Format("2006-01-02")

	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan ditolak."}`))
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
		http.Error(w, "Status pesanan tidak sedang menunggu konfirmasi", http.StatusBadRequest)
		return
	}

	dataOrder.Invoice.StatusPesanan = "dibatalkan"
	if err := dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var toko models.Toko
	idToko, _ := toko.GetIDTokoByOrder(idOrder)
	dataToko, _ := toko.GetToko(idToko)

	var karyawan models.Karyawan
	admins := karyawan.GetAdmins(idToko)

	var customer models.Customer
	dataCustomer, _ := customer.GetCustomer(strconv.Itoa(user.IDCustomer))

	var notif models.Notifikasi
	notif.Penerima = append(notif.Penerima, dataToko.Owner)
	notif.Penerima = append(notif.Penerima, admins...)
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
	idToko := vars["idToko"]

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	dataOrder.Invoice.StatusPesanan = "selesai"

	if dataOrder.Invoice.StatusPembayaran == "menunggu pembayaran" {
		dataOrder.Invoice.StatusPembayaran = "tidak lunas"
	}

	var pembukuan models.Pembukuan
	pembukuan.Jenis = "pemasukan"
	pembukuan.Keterangan = "Pesanan #" + dataOrder.Invoice.IDInvoice + " telah selesai."
	pembukuan.Nominal = dataOrder.Invoice.TotalBayar - dataOrder.Pengiriman.Ongkir
	pembukuan.TglTransaksi = time.Now().Format("2006-01-02")

	err := dataOrder.Invoice.UpdateInvoice(dataOrder.Invoice.IDInvoice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = pembukuan.CreatePembukuan(idToko)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan #` + dataOrder.Invoice.IDInvoice + ` telah diselesaikan"}`))
}
