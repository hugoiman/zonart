package middleware

import (
	"net/http"
	"strconv"
	"zonart/pkg/models"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// AuthOwner is middleware
func AuthOwner(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idToko := vars["idToko"]
		var toko models.Toko

		user := context.Get(r, "user").(*MyClaims)

		dataToko, err := toko.GetToko(idToko)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if user.IDCustomer != dataToko.IDOwner {
			http.Error(w, "Gagal! Anda tidak memiliki otoritas menjalankan fungsi ini.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// AuthOwnerAdmin is middleware
func AuthOwnerAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idToko := vars["idToko"]
		var toko models.Toko
		var karyawan models.Karyawan

		user := context.Get(r, "user").(*MyClaims)

		dataToko, err := toko.GetToko(idToko)
		if err != nil {
			http.Error(w, "Tidak ditemukan", http.StatusBadRequest)
			return
		}

		dataKaryawan, err := karyawan.AuthKaryawan(idToko, strconv.Itoa(user.IDCustomer))
		if user.IDCustomer != dataToko.IDOwner {
			if dataKaryawan.Posisi != "admin" || dataKaryawan.Status != "aktif" {
				http.Error(w, "Gagal! Anda tidak memiliki otoritas menjalankan fungsi ini.", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
