package controllers

// BiayaTambahanController is class
type BiayaTambahanController struct{}

// CreateBiayaTambahan is func
// func (btc BiayaTambahanController) CreateBiayaTambahan(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idOrder := vars["idOrder"]

// 	var tb models.BiayaTambahan
// 	var order models.Order

// 	if err := json.NewDecoder(r.Body).Decode(&tb); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	} else if err := validator.New().Struct(tb); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	dataOrder, err := order.GetOrder(idOrder)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	} else if dataOrder.Invoice.StatusPesanan == "selesai" {
// 		http.Error(w, "Status pesanan telah selesai.", http.StatusBadRequest)
// 		return
// 	}

// 	err = tb.CreateBiayaTambahan(idOrder)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	dataOrder.Invoice.TotalPembelian -= tb.Nominal
// 	dataOrder.Invoice.Tagihan += tb.Nominal
// 	dataOrder.Invoice.StatusPembayaran = "menunggu pembayaran"
// 	_ = dataOrder.UpdateBiayaOrder(idOrder)

// 	penerima := []int{}

// 	// send notif to customer
// 	var notif models.Notifikasi
// 	notif.IDPenerima = append(penerima, dataOrder.IDCustomer)
// 	notif.Pengirim = dataOrder.Invoice.NamaToko
// 	notif.Judul = "Terdapat Biaya Tambahan Baru"
// 	notif.Pesan = "Pesanan " + strconv.Itoa(dataOrder.IDOrder) + " mempunyai biaya tambahan baru."
// 	notif.Link = "/order/" + strconv.Itoa(dataOrder.IDOrder)
// 	notif.CreatedAt = time.Now().Format("2006-01-02")
// 	notif.CreateNotifikasi()

// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message":"Biaya tambahan telah ditambahkan."}`))
// }

// DeleteBiayaTambahan is func
// func (btc BiayaTambahanController) DeleteBiayaTambahan(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idOrder := vars["idOrder"]
// 	idBiayaTambahan := vars["idBiayaTambahan"]

// 	var bt models.BiayaTambahan
// 	var order models.Order
// 	// var newerOrder models.Order

// 	dataOrder, err := order.GetOrder(idOrder)
// 	if err != nil {
// 		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
// 		return
// 	} else if dataOrder.Invoice.StatusPesanan == "selesai" {
// 		http.Error(w, "Pesanan telah selesai.", http.StatusBadRequest)
// 		return
// 	}

// 	dataBT, err := bt.GetBiayaTambahan(idBiayaTambahan, idOrder)
// 	if err != nil {
// 		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
// 		return
// 	}

// 	err = dataBT.DeleteBiayaTambahan(idBiayaTambahan, idOrder)
// 	if err != nil {
// 		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
// 		return
// 	}

// 	dataNewerOrder, _ := order.GetOrder(idOrder)
// 	dataNewerOrder.Tagihan -= dataBT.Nominal
// 	dataNewerOrder.Total += dataBT.Nominal

// 	if dataNewerOrder.Tagihan <= 0 {
// 		dataNewerOrder.StatusPembayaran = "lunas"
// 	}

// 	_ = dataNewerOrder.UpdateBiayaOrder(idOrder)

// 	var notif models.Notifikasi
// 	notif.IDPenerima = append(notif.IDPenerima, dataOrder.IDCustomer)
// 	notif.Pengirim = dataNewerOrder.NamaToko
// 	notif.Judul = "Biaya tambahan telah dibatalkan"
// 	notif.Pesan = "Biaya tambahan berupa " + dataBT.Item + "(Rp " + strconv.Itoa(dataBT.Nominal) + ") telah dibatalkan. No pesanan " + strconv.Itoa(dataOrder.IDOrder)
// 	notif.Link = "/order/" + idOrder
// 	notif.CreatedAt = time.Now().Format("2006-01-02")
// 	notif.CreateNotifikasi()

// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

// }
