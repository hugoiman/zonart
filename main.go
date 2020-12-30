package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"zonart/middleware"
	"zonart/pkg/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	origins := handlers.AllowedOrigins([]string{"*"})

	api := router.PathPrefix("").Subrouter()

	// Instansiasi Class Controller
	var mw middleware.MiddleWare
	var auth controllers.AuthController
	var bt controllers.BiayaTambahanController
	var customer controllers.CustomerController
	var faq controllers.FaqController
	var galeri controllers.GaleriController
	var gopsi controllers.GrupOpsiController
	var goproduk controllers.GrupOpsiProdukController
	var karyawan controllers.KaryawanController
	var opsi controllers.OpsiController
	var order controllers.OrderController
	// var pembukuan controllers.PembukuanController
	var penangan controllers.PenanganController
	var produk controllers.ProdukController
	var toko controllers.TokoController
	var undangan controllers.UndanganController

	api.Use(mw.AuthToken)

	router.HandleFunc("/api/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/register", customer.Register).Methods("POST")
	router.HandleFunc("/api/reset-password", auth.ResetPassword).Methods("POST")

	api.HandleFunc("/api/customer", customer.GetCustomer).Methods("GET")
	api.HandleFunc("/api/customer", customer.UpdateProfil).Methods("PUT")
	api.HandleFunc("/api/change-password", customer.ChangePassword).Methods("PUT")
	// api.HandleFunc("/api/customers", controllers.GetCustomers).Methods("GET")

	router.HandleFunc("/api/toko/{id}", toko.GetToko).Methods("GET")
	api.HandleFunc("/api/toko", toko.CreateToko).Methods("POST")
	api.HandleFunc("/api/toko/{idToko}", mw.AuthOwner(toko.UpdateToko)).Methods("PUT")
	// GetMyListToko

	router.HandleFunc("/api/galeri/{idToko}", galeri.GetGaleris).Methods("GET")
	api.HandleFunc("/api/galeri/{idToko}", mw.AuthOwnerAdmin(galeri.CreateGaleri)).Methods("POST")
	api.HandleFunc("/api/galeri/{idToko}/{idGaleri}", mw.AuthOwnerAdmin(galeri.DeleteGaleri)).Methods("DELETE")

	router.HandleFunc("/api/faq/{idToko}", faq.GetFaqs).Methods("GET")
	router.HandleFunc("/api/faq/{idToko}/{idFaq}", faq.GetFaq).Methods("GET")
	api.HandleFunc("/api/faq/{idToko}", mw.AuthOwnerAdmin(faq.CreateFaq)).Methods("POST")
	api.HandleFunc("/api/faq/{idToko}/{idFaq}", mw.AuthOwnerAdmin(faq.DeleteFaq)).Methods("DELETE")

	api.HandleFunc("/api/karyawan/{idToko}", mw.AuthOwnerAdmin(karyawan.GetKaryawans)).Methods("GET")
	api.HandleFunc("/api/karyawan/{idToko}/{idKaryawan}", mw.AuthOwner(karyawan.GetKaryawan)).Methods("GET")
	api.HandleFunc("/api/karyawan/{idToko}/{idKaryawan}", mw.AuthOwner(karyawan.UpdateKaryawan)).Methods("PUT")

	api.HandleFunc("/api/undangan/{idToko}", mw.AuthOwner(undangan.GetUndangans)).Methods("GET")
	api.HandleFunc("/api/undangan/{idToko}/{idUndangan}", undangan.GetUndangan).Methods("GET")
	api.HandleFunc("/api/undangan/{idToko}", mw.AuthOwner(undangan.UndangKaryawan)).Methods("POST")
	api.HandleFunc("/api/undangan-tolak/{idToko}/{idUndangan}", undangan.TolakUndangan).Methods("POST")
	api.HandleFunc("/api/undangan-terima/{idToko}/{idUndangan}", undangan.TerimaUndangan).Methods("POST")
	api.HandleFunc("/api/undangan-batal/{idToko}/{idUndangan}/{idCustomer}", mw.AuthOwner(undangan.BatalkanUndangan)).Methods("POST")

	router.HandleFunc("/api/produk/{idToko}", produk.GetProduks).Methods("GET")
	router.HandleFunc("/api/produk/{idToko}/{idProduk}", produk.GetProduk).Methods("GET")
	api.HandleFunc("/api/produk/{idToko}", mw.AuthOwnerAdmin(produk.CreateProduk)).Methods("POST")
	api.HandleFunc("/api/produk/{idToko}/{idProduk}", mw.AuthOwnerAdmin(produk.UpdateProduk)).Methods("PUT")

	router.HandleFunc("/api/grup-opsi/{idToko}", gopsi.GetGrupOpsis).Methods("GET")
	router.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}", gopsi.GetGrupOpsi).Methods("GET")
	api.HandleFunc("/api/grup-opsi/{idToko}", mw.AuthOwnerAdmin(gopsi.CreateGrupOpsi)).Methods("POST")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}", mw.AuthOwnerAdmin(gopsi.UpdateGrupOpsi)).Methods("PUT")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}", mw.AuthOwnerAdmin(gopsi.DeleteGrupOpsi)).Methods("DELETE")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}/{idProduk}", mw.AuthOwnerAdmin(goproduk.SambungGrupOpsikeProduk)).Methods("POST")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}/{idProduk}", mw.AuthOwnerAdmin(goproduk.PutusGrupOpsidiProduk)).Methods("DELETE")

	api.HandleFunc("/api/opsi/{idToko}/{idGrupOpsi}/{idOpsi}", mw.AuthOwnerAdmin(opsi.DeleteOpsi)).Methods("DELETE")

	api.HandleFunc("/api/grup-opsi-produk/{idToko}/{idGrupOpsi}", mw.AuthOwnerAdmin(goproduk.GetGrupOpsiProduks)).Methods("GET")

	// my list order customer
	// api.HandleFunc("/api/order", controllers.GetOrders).Methods("GET")
	// detail order customer
	api.HandleFunc("/api/order/{idOrder}", order.GetOrder).Methods("GET")
	// detail order toko
	api.HandleFunc("/api/order/{idToko}/{idOrder}", order.GetOrderToko).Methods("GET")
	// list order toko
	// api.HandleFunc("/api/orders/{idToko}", controllers.GetOrder).Methods("GET")
	// list order editor
	// api.HandleFunc("/api/orders/{idToko}", controllers.GetOrder).Methods("GET")
	api.HandleFunc("/api/order/{idToko}/{idProduk}", order.CreateOrder).Methods("POST")

	api.HandleFunc("/api/biaya-tambahan/{idToko}/{idOrder}", mw.AuthOwnerAdmin(bt.CreateBiayaTambahans)).Methods("POST")
	api.HandleFunc("/api/biaya-tambahan/{idToko}/{idOrder}/{idBiayaTambahan}", mw.AuthOwnerAdmin(bt.DeleteBiayaTambahan)).Methods("DELETE")

	api.HandleFunc("/api/penangan/{idToko}/{idOrder}", mw.AuthOwnerAdmin(penangan.SetPenangan)).Methods("POST")

	os.Setenv("PORT", "8080")
	port := "8080"

	fmt.Println("Server running at :", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(origins)(router)))
}
