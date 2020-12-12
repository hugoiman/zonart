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

// CreateOrder is func
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	vars := mux.Vars(r)

	idToko := vars["idToko"]
	idProduk := vars["idProduk"]
	idCustomer := user.IDCustomer

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var produk models.Produk
	dataProduk, err := produk.GetProduk(idToko, idProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get total opsi yg dipilih berdasarkan grupOpsi
	var totalOpsiGrup = make(map[int]int)
	for _, v := range order.OpsiOrder {
		totalOpsiGrup[v.IDGrupOpsi]++
	}

	// Cek batas min & max yg diperbolehkan grupOpsi
	for _, vGrupOpsi := range dataProduk.GrupOpsi {
		if totalOpsiGrup[vGrupOpsi.IDGrupOpsi] < vGrupOpsi.Min {
			http.Error(w, "Gagal! "+vGrupOpsi.NamaGrup+" kurang dari batas minimal", http.StatusBadRequest)
			return
		} else if totalOpsiGrup[vGrupOpsi.IDGrupOpsi] > vGrupOpsi.Max {
			http.Error(w, "Gagal! "+vGrupOpsi.NamaGrup+" melebihi batas maximal", http.StatusBadRequest)
			return
		}
	}

	// store dataGrupOpsiProduk in var bertipe multidimensi map > map[idGrupOpsi]GrupOpsi
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

	// store dataOpsiProduk in var bertipe multidimensi map > map[idGrupOpsi]map[idOpsi]Opsi
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

	order.IDCustomer = idCustomer
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

	if order.JenisPesanan == "cetak" {
		var toko models.Toko
		dataToko, _ := toko.GetToko(idToko)
		ongkir := GetOngkir(dataToko.Kota, order.Pengiriman.Kota, order.Pengiriman.Kurir, order.Pengiriman.Berat)
		order.Pengiriman.Ongkir = ongkir
		order.HargaProduk = dataProduk.HargaCetak
	} else if order.JenisPesanan == "soft copy" {
		order.HargaProduk = dataProduk.HargaSoftCopy
	}

	order.Total = (order.HargaProduk * order.Pcs) + order.TotalHargaWajah + order.TotalHargaOpsi + order.TotalTambahanBiaya

	idOrder, err := order.CreateOrder(idToko, idProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, vOpsiOrder := range order.OpsiOrder {
		err = vOpsiOrder.CreateOpsiOrder(strconv.Itoa(idOrder))
		if err != nil {
			_ = order.DeleteOrder(strconv.Itoa(idOrder))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = order.Pengiriman.CreatePengiriman(strconv.Itoa(idOrder))
	if err != nil {
		_ = order.DeleteOrder(strconv.Itoa(idOrder))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := json.Marshal(order)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// GetOngkir is setter
func GetOngkir(asal, tujuan, kurir string, berat int) int {
	return 10000
}
