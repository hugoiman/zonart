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

	message, _ := json.Marshal(&dataOrder)

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
	} else if position["position"] == "editor" && strconv.Itoa(dataOrder.GetPenangan().GetIDKaryawan()) != position["idKaryawan"].(string) {
		http.Error(w, "Pesanan tidak ditemukan.", http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(&dataOrder)

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

	message, _ := json.Marshal(&dataOrder)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOrders is get all order customer
func (oc OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var order models.Order

	dataOrder := order.GetOrders(strconv.Itoa(user.IDCustomer))

	message, _ := json.Marshal(&dataOrder)

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

	message, _ := json.Marshal(&dataOrder)

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

	err = oc.validateGrupOpsiOrder(&order, dataProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dtProduk, _ := json.Marshal(&dataProduk)

	// Input detail(namaGrup, opsi, harga, berat, perProduk) ke OpsiOrder
	err = oc.setOpsiOrder(&order, dtProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.GetInvoice().SetIDInvoice(time.Now().Format("020106150405"))
	order.SetPemesan(user.IDCustomer)
	order.GetProdukOrder().SetNamaProduk(dataProduk.GetNamaProduk())
	order.GetProdukOrder().SetBeratProduk(dataProduk.GetBerat())
	order.GetInvoice().SetStatusPesanan("menunggu konfirmasi")
	order.GetInvoice().SetStatusPembayaran("-")
	order.GetInvoice().SetTotalBayar(0)
	order.GetInvoice().SetTagihan(0)
	order.GetProdukOrder().SetHargaSatuanWajah(dataProduk.GetHargaWajah())
	order.SetTglOrder(time.Now().Format("2006-01-02"))

	// Total Harga Wajah
	totalHargaWajah := order.GetTambahanWajah() * order.GetProdukOrder().GetHargaSatuanWajah()

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
	order.GetInvoice().SetTotalPembelian((order.GetProdukOrder().GetHargaProduk() * order.GetPcs()) + totalHargaWajah + totalHargaOpsi + order.GetPengiriman().GetOngkir())

	// create invoice
	err = order.GetInvoice().CreateInvoice(idToko, strconv.Itoa(user.IDCustomer))
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
		fileOrder.SetFoto(v)
		order.SetFileOrder(append(order.GetFileOrder(), fileOrder))
	}

	// create order
	idOrder, err := order.CreateOrder(idToko, idProduk)
	if err != nil {
		_ = order.GetInvoice().DeleteInvoice(order.GetInvoice().GetIDInvoice())
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
	notif.SetPenerima(append(notif.GetPenerima(), dataToko.GetOwner()))
	notif.SetPenerima(append(notif.GetPenerima(), admins...))
	notif.SetPengirim(dataCustomer.GetNama())
	notif.SetJudul("Permintaan pesanan baru")
	notif.SetPesan(notif.GetPengirim() + " telah memesan produk " + order.GetProdukOrder().GetNamaProduk() + ". Pesanan #" + order.GetInvoice().GetIDInvoice())
	notif.SetLink(dataToko.GetSlug() + "/pesanan/" + strconv.Itoa(idOrder))
	notif.SetCreatedAt(order.GetTglOrder())
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Mohon tunggu konfirmasi penjual. Terimakasih.","idOrder":"` + strconv.Itoa(idOrder) + `"}`))
}

func (oc OrderController) validateGrupOpsiOrder(order *models.Order, dataProduk models.Produk) error {
	grupOpsi, _ := json.Marshal(order)
	for _, vGrupOpsi := range dataProduk.GetGrupOpsi() {
		totalOpsiGrup := gjson.Get(string(grupOpsi), "opsiOrder.#(idGrupOpsi=="+strconv.Itoa(vGrupOpsi.GetIDGrupOpsi())+")#").Array()
		if len(totalOpsiGrup) < vGrupOpsi.GetMin() && ((order.GetJenisPesanan() == "cetak" && vGrupOpsi.GetHardCopy()) || (order.GetJenisPesanan() == "soft copy" && vGrupOpsi.GetSoftCopy())) {
			return errors.New("Opsi " + vGrupOpsi.GetNamaGrup() + " kurang dari batas minimal")
		} else if len(totalOpsiGrup) > vGrupOpsi.GetMax() && ((order.GetJenisPesanan() == "cetak" && vGrupOpsi.GetHardCopy()) || (order.GetJenisPesanan() == "soft copy" && vGrupOpsi.GetSoftCopy())) {
			return errors.New("Opsi " + vGrupOpsi.GetNamaGrup() + " melebihi batas maksimal")
		} else if ((order.GetJenisPesanan() == "cetak" && !vGrupOpsi.GetHardCopy()) || (order.GetJenisPesanan() == "soft copy" && !vGrupOpsi.GetSoftCopy())) && len(totalOpsiGrup) > 0 {
			return errors.New(vGrupOpsi.GetNamaGrup() + " dapat diisi jika jenis pesanan selain " + order.GetJenisPesanan())
		}
	}
	return nil
}

func (oc OrderController) setOpsiOrder(order *models.Order, dtProduk []byte) error {
	for k, v := range order.GetOpsiOrder() {
		dataGop := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+")#").Array()
		namaGrup := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+").namaGrup").String()
		spesificRequest := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+").spesificRequest").Bool()
		dataOpsiProduk := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+")#.opsi.#(idOpsi=="+strconv.Itoa(v.GetIDOpsiOrder())+")").Array()
		if len(dataGop) == 0 {
			return errors.New("Grup Opsi tidak ditemukan.")
		} else if len(dataOpsiProduk) == 0 && spesificRequest == false {
			return errors.New("Spesific Request tidak diizinkan.")
		} else if len(dataOpsiProduk) == 0 && spesificRequest == true {
			order.GetOpsiOrder()[k].SetNamaGrup(namaGrup)
			order.GetOpsiOrder()[k].SetHarga(0)
			order.GetOpsiOrder()[k].SetBerat(0)
			order.GetOpsiOrder()[k].SetPerProduk(false)
		} else {
			opsi := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+").opsi.#(idOpsi=="+strconv.Itoa(v.GetIDOpsiOrder())+").opsi").String()
			harga := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+").opsi.#(idOpsi=="+strconv.Itoa(v.GetIDOpsiOrder())+").harga").Int()
			berat := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+").opsi.#(idOpsi=="+strconv.Itoa(v.GetIDOpsiOrder())+").berat").Int()
			perProduk := gjson.Get(string(dtProduk), "grupOpsi.#(idGrupOpsi=="+strconv.Itoa(v.GetIDGrupOpsi())+").opsi.#(idOpsi=="+strconv.Itoa(v.GetIDOpsiOrder())+").perProduk").Bool()

			order.GetOpsiOrder()[k].SetNamaGrup(namaGrup)
			order.GetOpsiOrder()[k].SetOpsi(opsi)
			order.GetOpsiOrder()[k].SetHarga(int(harga))
			order.GetOpsiOrder()[k].SetBerat(int(berat))
			order.GetOpsiOrder()[k].SetPerProduk(perProduk)
		}
	}
	return nil
}

func (oc OrderController) setPengiriman(order *models.Order, dataProduk models.Produk, dataToko models.Toko, dtProduk []byte, totalBeratOpsi int) error {
	if order.GetJenisPesanan() == "soft copy" {
		order.GetProdukOrder().SetHargaProduk(int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="soft copy").harga`).Int()))
		order.GetPengiriman().SetKota(dataToko.GetKota())
		order.GetPengiriman().SetKurir("")
		order.GetPengiriman().SetAlamat("")
		order.GetPengiriman().SetLabel("")
		order.GetPengiriman().SetService("")
	} else if order.GetJenisPesanan() == "cetak" && order.GetPengiriman().GetKodeKurir() == "cod" {
		order.GetProdukOrder().SetHargaProduk(int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="cetak").harga`).Int()))
		order.GetPengiriman().SetKota(dataToko.GetKota())
		order.GetPengiriman().SetKurir("COD (Cash On Delivery)")
		order.GetPengiriman().SetAlamat("")
		order.GetPengiriman().SetLabel("")
		order.GetPengiriman().SetService("")
	} else if order.GetJenisPesanan() == "cetak" && order.GetPengiriman().GetKodeKurir() != "cod" {
		var rj RajaOngkir
		order.GetPengiriman().SetBerat((dataProduk.GetBerat() * order.GetPcs()) + totalBeratOpsi)
		asal, _ := rj.GetIDKota(dataToko.GetKota())
		tujuan, _ := rj.GetIDKota(order.GetPengiriman().GetKota())
		ongkir, estimasi, kurir, err := rj.GetOngkir(asal, tujuan, order.GetPengiriman().GetKodeKurir(), order.GetPengiriman().GetService(), strconv.Itoa(order.GetPengiriman().GetBerat()))
		if !err {
			return errors.New("Terjadi kesalahan. Mohon periksa data pengiriman.")
		}

		order.GetPengiriman().SetKurir(kurir)
		order.GetProdukOrder().SetHargaProduk(int(gjson.Get(string(dtProduk), `jenisPemesanan.#(jenis=="cetak").harga`).Int()))
		order.GetPengiriman().SetOngkir(ongkir)
		order.GetPengiriman().SetEstimasi(estimasi)
	}
	return nil
}

// hitungHargaBeratOpsi is func
func (oc OrderController) hitungHargaBeratOpsi(o models.Order) (int, int) {
	var totalHargaOpsi, totalBeratOpsi int
	for _, valueOpsi := range o.GetOpsiOrder() {
		if valueOpsi.GetPerProduk() == false {
			totalHargaOpsi += valueOpsi.GetHarga()
			totalBeratOpsi += valueOpsi.GetBerat()
		} else {
			totalHargaOpsi += valueOpsi.GetHarga() * o.GetPcs()
			totalBeratOpsi += valueOpsi.GetBerat() * o.GetPcs()
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

	dataOrder.GetInvoice().SetTagihan(dataOrder.GetInvoice().GetTotalPembelian())
	dataOrder.GetInvoice().SetStatusPembayaran("menunggu pembayaran")
	dataOrder.GetInvoice().SetStatusPesanan("diproses")

	err = dataOrder.GetInvoice().UpdateInvoice(dataOrder.GetInvoice().GetIDInvoice())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var notif models.Notifikasi

	notif.SetPenerima(append(notif.GetPenerima(), dataOrder.GetPemesan()))
	notif.SetPengirim(dataOrder.GetInvoice().GetNamaToko())
	notif.SetJudul("Pesanan telah dikonfirmasi")
	notif.SetPesan("Pesanan #" + dataOrder.GetInvoice().GetIDInvoice() + " sedang diproses. Silahkan melakukan pembayaran.")
	notif.SetLink("/order?id=" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
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

	dataOrder.GetInvoice().SetStatusPesanan("ditolak")

	err = dataOrder.GetInvoice().UpdateInvoice(dataOrder.GetInvoice().GetIDInvoice())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var notif models.Notifikasi

	notif.SetPenerima(append(notif.GetPenerima(), dataOrder.GetPemesan()))
	notif.SetPengirim(dataOrder.GetInvoice().GetNamaToko())
	notif.SetJudul("Pesanan ditolak")
	notif.SetPesan("Pesanan #" + dataOrder.GetInvoice().GetIDInvoice() + " ditolak. Keterangan: " + data.Keterangan)
	notif.SetLink("/order?id=" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))

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

	order.SetWaktu(waktu)

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
	if dataOrder.GetInvoice().GetStatusPesanan() != "menunggu konfirmasi" {
		http.Error(w, "Status pesanan tidak sedang menunggu konfirmasi", http.StatusBadRequest)
		return
	}

	dataOrder.GetInvoice().SetStatusPesanan("dibatalkan")
	if err := dataOrder.GetInvoice().UpdateInvoice(dataOrder.GetInvoice().GetIDInvoice()); err != nil {
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
	notif.SetPenerima(append(notif.GetPenerima(), dataToko.GetOwner()))
	notif.SetPenerima(append(notif.GetPenerima(), admins...))
	notif.SetPengirim(dataCustomer.GetNama())
	notif.SetJudul("Pesanan dibatalkan")
	notif.SetPesan("Pesanan #" + dataOrder.GetInvoice().GetIDInvoice() + " dibatalkan oleh pembeli.")
	notif.SetLink("/order?id=" + idOrder)
	notif.SetCreatedAt(time.Now().Format("2006-01-02"))
	notif.CreateNotifikasi()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan  #` + dataOrder.GetInvoice().GetIDInvoice() + ` telah dibatalkan"}`))
}

// FinishOrder is func
func (oc OrderController) FinishOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	idToko := vars["idToko"]

	var order models.Order
	dataOrder, _ := order.GetOrder(idOrder)
	dataOrder.GetInvoice().SetStatusPesanan("selesai")

	if dataOrder.GetInvoice().GetStatusPembayaran() == "menunggu pembayaran" {
		dataOrder.GetInvoice().SetStatusPembayaran("tidak lunas")
	}

	var pembukuan models.Pembukuan
	pembukuan.SetJenis("pemasukan")
	pembukuan.SetKeterangan("Pesanan #" + dataOrder.GetInvoice().GetIDInvoice() + " telah selesai.")
	pembukuan.SetNominal(dataOrder.GetInvoice().GetTotalBayar() - dataOrder.GetPengiriman().GetOngkir())
	pembukuan.SetTglTransaksi(time.Now().Format("2006-01-02"))

	err := dataOrder.GetInvoice().UpdateInvoice(dataOrder.GetInvoice().GetIDInvoice())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = pembukuan.CreatePembukuan(idToko)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Pesanan #` + dataOrder.GetInvoice().GetIDInvoice() + ` telah diselesaikan"}`))
}
