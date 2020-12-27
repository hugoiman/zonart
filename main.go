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
	var btc controllers.BiayaTambahanController
	var cc controllers.CustomerController

	api.Use(mw.AuthToken)

	router.HandleFunc("/api/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/register", cc.Register).Methods("POST")
	router.HandleFunc("/api/reset-password", auth.ResetPassword).Methods("POST")

	api.HandleFunc("/api/customer", cc.GetCustomer).Methods("GET")
	api.HandleFunc("/api/customer", cc.UpdateProfil).Methods("PUT")
	api.HandleFunc("/api/change-password", cc.ChangePassword).Methods("PUT")
	// api.HandleFunc("/api/customers", controllers.GetCustomers).Methods("GET")

	router.HandleFunc("/api/toko/{id}", controllers.GetToko).Methods("GET")
	api.HandleFunc("/api/toko", controllers.CreateToko).Methods("POST")
	api.HandleFunc("/api/toko/{idToko}", mw.AuthOwner(controllers.UpdateToko)).Methods("PUT")
	// GetMyListToko

	router.HandleFunc("/api/galeri/{idToko}", controllers.GetGaleris).Methods("GET")
	api.HandleFunc("/api/galeri/{idToko}", mw.AuthOwnerAdmin(controllers.CreateGaleri)).Methods("POST")
	api.HandleFunc("/api/galeri/{idToko}/{idGaleri}", mw.AuthOwnerAdmin(controllers.DeleteGaleri)).Methods("DELETE")

	router.HandleFunc("/api/faq/{idToko}", controllers.GetFaqs).Methods("GET")
	router.HandleFunc("/api/faq/{idToko}/{idFaq}", controllers.GetFaq).Methods("GET")
	api.HandleFunc("/api/faq/{idToko}", mw.AuthOwnerAdmin(controllers.CreateFaq)).Methods("POST")
	api.HandleFunc("/api/faq/{idToko}/{idFaq}", mw.AuthOwnerAdmin(controllers.DeleteFaq)).Methods("DELETE")

	api.HandleFunc("/api/karyawan/{idToko}", mw.AuthOwnerAdmin(controllers.GetKaryawans)).Methods("GET")
	api.HandleFunc("/api/karyawan/{idToko}/{idKaryawan}", mw.AuthOwner(controllers.GetKaryawan)).Methods("GET")
	api.HandleFunc("/api/karyawan/{idToko}/{idKaryawan}", mw.AuthOwner(controllers.UpdateKaryawan)).Methods("PUT")

	api.HandleFunc("/api/undangan/{idToko}", mw.AuthOwner(controllers.GetUndangans)).Methods("GET")
	api.HandleFunc("/api/undangan/{idToko}/{idUndangan}", controllers.GetUndangan).Methods("GET")
	api.HandleFunc("/api/undangan/{idToko}", mw.AuthOwner(controllers.UndangKaryawan)).Methods("POST")
	api.HandleFunc("/api/undangan-tolak/{idToko}/{idUndangan}", controllers.TolakUndangan).Methods("POST")
	api.HandleFunc("/api/undangan-terima/{idToko}/{idUndangan}", controllers.TerimaUndangan).Methods("POST")
	api.HandleFunc("/api/undangan-batal/{idToko}/{idUndangan}/{idCustomer}", mw.AuthOwner(controllers.BatalkanUndangan)).Methods("POST")

	router.HandleFunc("/api/produk/{idToko}", controllers.GetProduks).Methods("GET")
	router.HandleFunc("/api/produk/{idToko}/{idProduk}", controllers.GetProduk).Methods("GET")
	api.HandleFunc("/api/produk/{idToko}", mw.AuthOwnerAdmin(controllers.CreateProduk)).Methods("POST")
	api.HandleFunc("/api/produk/{idToko}/{idProduk}", mw.AuthOwnerAdmin(controllers.UpdateProduk)).Methods("PUT")

	router.HandleFunc("/api/grup-opsi/{idToko}", controllers.GetGrupOpsis).Methods("GET")
	router.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}", controllers.GetGrupOpsi).Methods("GET")
	api.HandleFunc("/api/grup-opsi/{idToko}", mw.AuthOwnerAdmin(controllers.CreateGrupOpsi)).Methods("POST")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}", mw.AuthOwnerAdmin(controllers.UpdateGrupOpsi)).Methods("PUT")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}", mw.AuthOwnerAdmin(controllers.DeleteGrupOpsi)).Methods("DELETE")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}/{idProduk}", mw.AuthOwnerAdmin(controllers.SambungGrupOpsikeProduk)).Methods("POST")
	api.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}/{idProduk}", mw.AuthOwnerAdmin(controllers.PutusGrupOpsidiProduk)).Methods("DELETE")

	api.HandleFunc("/api/opsi/{idToko}/{idGrupOpsi}/{idOpsi}", mw.AuthOwnerAdmin(controllers.DeleteOpsi)).Methods("DELETE")

	api.HandleFunc("/api/grup-opsi-produk/{idToko}/{idGrupOpsi}", mw.AuthOwnerAdmin(controllers.GetGrupOpsiProduks)).Methods("GET")

	// my list order customer
	// api.HandleFunc("/api/order", controllers.GetOrders).Methods("GET")
	// detail order customer
	api.HandleFunc("/api/order/{idOrder}", controllers.GetOrder).Methods("GET")
	// detail order toko
	api.HandleFunc("/api/order/{idToko}/{idOrder}", controllers.GetOrderToko).Methods("GET")
	// list order toko
	// api.HandleFunc("/api/orders/{idToko}", controllers.GetOrder).Methods("GET")
	// list order editor
	// api.HandleFunc("/api/orders/{idToko}", controllers.GetOrder).Methods("GET")
	api.HandleFunc("/api/order/{idToko}/{idProduk}", controllers.CreateOrder).Methods("POST")

	api.HandleFunc("/api/biaya-tambahan/{idToko}/{idOrder}", mw.AuthOwnerAdmin(btc.CreateBiayaTambahans)).Methods("POST")
	api.HandleFunc("/api/biaya-tambahan/{idToko}/{idOrder}/{idBiayaTambahan}", mw.AuthOwnerAdmin(btc.DeleteBiayaTambahan)).Methods("DELETE")

	api.HandleFunc("/api/penangan/{idToko}/{idOrder}", mw.AuthOwnerAdmin(controllers.SetPenangan)).Methods("POST")

	os.Setenv("PORT", "8080")
	port := "8080"

	fmt.Println("Server running at :", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(origins)(router)))
}
