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

	router.HandleFunc("/api/galeri/{idToko}", controllers.GetGaleris).Methods("GET")
	auth.HandleFunc("/api/galeri/{idToko}", mw.AuthOwnerAdmin(controllers.CreateGaleri)).Methods("POST")
	auth.HandleFunc("/api/galeri/{idToko}/{idGaleri}", mw.AuthOwnerAdmin(controllers.DeleteGaleri)).Methods("DELETE")

	router.HandleFunc("/api/faq/{idToko}", controllers.GetFaqs).Methods("GET")
	router.HandleFunc("/api/faq/{idToko}/{idFaq}", controllers.GetFaq).Methods("GET")
	auth.HandleFunc("/api/faq/{idToko}", mw.AuthOwnerAdmin(controllers.CreateFaq)).Methods("POST")
	auth.HandleFunc("/api/faq/{idToko}/{idFaq}", mw.AuthOwnerAdmin(controllers.DeleteFaq)).Methods("DELETE")

	os.Setenv("PORT", "8080")
	port := "8080"

	fmt.Println("Server running at :", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(origins)(router)))
}
