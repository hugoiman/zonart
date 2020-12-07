package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mw "zonart/middleware"
	"zonart/pkg/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	origins := handlers.AllowedOrigins([]string{"*"})

	auth := router.PathPrefix("").Subrouter()
	auth.Use(mw.AuthToken)

	router.HandleFunc("/api/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/register", controllers.Register).Methods("POST")
	router.HandleFunc("/api/reset-password", controllers.ResetPassword).Methods("POST")

	auth.HandleFunc("/api/customer", controllers.GetCustomer).Methods("GET")
	auth.HandleFunc("/api/customer", controllers.UpdateProfil).Methods("PUT")
	auth.HandleFunc("/api/change-password", controllers.ChangePassword).Methods("PUT")
	// auth.HandleFunc("/api/customers", controllers.GetCustomers).Methods("GET")

	router.HandleFunc("/api/toko/{id}", controllers.GetToko).Methods("GET")
	auth.HandleFunc("/api/toko", controllers.CreateToko).Methods("POST")
	auth.HandleFunc("/api/toko/{idToko}", mw.AuthOwner(controllers.UpdateToko)).Methods("PUT")
	// GetMyListToko

	router.HandleFunc("/api/galeri/{idToko}", controllers.GetGaleris).Methods("GET")
	auth.HandleFunc("/api/galeri/{idToko}", mw.AuthOwnerAdmin(controllers.CreateGaleri)).Methods("POST")
	auth.HandleFunc("/api/galeri/{idToko}/{idGaleri}", mw.AuthOwnerAdmin(controllers.DeleteGaleri)).Methods("DELETE")

	router.HandleFunc("/api/faq/{idToko}", controllers.GetFaqs).Methods("GET")
	router.HandleFunc("/api/faq/{idToko}/{idFaq}", controllers.GetFaq).Methods("GET")
	auth.HandleFunc("/api/faq/{idToko}", mw.AuthOwnerAdmin(controllers.CreateFaq)).Methods("POST")
	auth.HandleFunc("/api/faq/{idToko}/{idFaq}", mw.AuthOwnerAdmin(controllers.DeleteFaq)).Methods("DELETE")

	auth.HandleFunc("/api/karyawan/{idToko}", mw.AuthOwnerAdmin(controllers.GetKaryawans)).Methods("GET")
	auth.HandleFunc("/api/karyawan/{idToko}/{idKaryawan}", mw.AuthOwner(controllers.GetKaryawan)).Methods("GET")
	auth.HandleFunc("/api/karyawan/{idToko}/{idKaryawan}", mw.AuthOwner(controllers.UpdateKaryawan)).Methods("PUT")

	auth.HandleFunc("/api/undangan/{idToko}", mw.AuthOwner(controllers.UndangKaryawan)).Methods("POST")
	auth.HandleFunc("/api/undangan-tolak/{idToko}/{idUndangan}", controllers.TolakUndangan).Methods("POST")
	auth.HandleFunc("/api/undangan-terima/{idToko}/{idUndangan}", controllers.TerimaUndangan).Methods("POST")
	auth.HandleFunc("/api/undangan-batal/{idToko}/{idUndangan}/{idCustomer}", mw.AuthOwner(controllers.BatalkanUndangan)).Methods("POST")

	router.HandleFunc("/api/produk/{idToko}", controllers.GetProduks).Methods("GET")
	router.HandleFunc("/api/produk/{idToko}/{idProduk}", controllers.GetProduk).Methods("GET")

	router.HandleFunc("/api/grup-opsi/{idToko}", controllers.GetGrupOpsis).Methods("GET")
	router.HandleFunc("/api/grup-opsi/{idToko}/{idGrupOpsi}", controllers.GetGrupOpsi).Methods("GET")

	router.HandleFunc("/api/grup-opsi-produk/{idToko}/{idProduk}", controllers.GetGrupOpsiProduk).Methods("GET")

	os.Setenv("PORT", "8080")
	port := "8080"

	fmt.Println("Server running at :", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(origins)(router)))
}
