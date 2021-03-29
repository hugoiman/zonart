package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"zonart/pkg/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// MySigningKey is signature
var MySigningKey = []byte("jwt super secret keys")

// MyClaims is Credential
type MyClaims struct {
	IDCustomer int    `json:"idCustomer"`
	Username   string `json:"username"`
	jwt.StandardClaims
}

// MiddleWare is class
type MiddleWare struct {
}

// AuthToken is middleware
func (mw MiddleWare) AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			http.Error(w, "Dibutuhkan otentikasi. Silahkan melakukan login.", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		claims := &MyClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return MySigningKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Token invalid: "+err.Error(), http.StatusBadRequest) // Token expired/key tidak cocok(invalid)
			return
		}

		context.Set(r, "user", claims)
		// fmt.Printf("%+v", claims)
		next.ServeHTTP(w, r)
	})
}

// AuthOwner is middleware
func (mw MiddleWare) AuthOwner(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idToko := vars["idToko"]
		var toko models.Toko

		user := context.Get(r, "user").(*MyClaims)

		dataToko, err := toko.GetToko(idToko)
		if err != nil || user.IDCustomer != dataToko.IDOwner {
			http.Error(w, "Anda tidak memiliki hak akses.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// AuthOwnerAdmin is middleware
func (mw MiddleWare) AuthOwnerAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idToko := vars["idToko"]
		var toko models.Toko
		var karyawan models.Karyawan

		user := context.Get(r, "user").(*MyClaims)

		dataToko, _ := toko.GetToko(idToko)
		// if err != nil {
		// 	http.Error(w, "Tidak ditemukan", http.StatusBadRequest)
		// 	return
		// }

		dataKaryawan, err := karyawan.GetKaryawanByIDCustomer(idToko, strconv.Itoa(user.IDCustomer))
		if user.IDCustomer != dataToko.IDOwner && (dataKaryawan.Posisi != "admin" || dataKaryawan.Status != "aktif" || err != nil) {
			http.Error(w, "Anda tidak memiliki hak akses.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// AuthEditor is middleware
func (mw MiddleWare) AuthEditor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idToko := vars["idToko"]
		var karyawan models.Karyawan

		user := context.Get(r, "user").(*MyClaims)

		dataKaryawan, err := karyawan.GetKaryawanByIDCustomer(idToko, strconv.Itoa(user.IDCustomer))
		if dataKaryawan.Posisi != "editor" || dataKaryawan.Status != "aktif" || err != nil {
			http.Error(w, "Anda tidak memiliki hak akses.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// CustomerOrder is middleware
func (mw MiddleWare) CustomerOrder(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idOrder := vars["idOrder"]
		var order models.Order

		user := context.Get(r, "user").(*MyClaims)

		dataOrder, err := order.GetOrder(idOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if user.IDCustomer != dataOrder.IDCustomer {
			http.Error(w, "Anda tidak memiliki otoritas pada pesanan ini.", http.StatusForbidden)
			return
		} else if r.Method != http.MethodGet && dataOrder.Invoice.StatusPesanan == "selesai" {
			http.Error(w, "Pesanan sudah diselesaikan.", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// OwnerAdminOrder is middleware
func (mw MiddleWare) OwnerAdminOrder(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idOrder := vars["idOrder"]
		idToko := vars["idToko"]
		var order models.Order

		dataOrder, err := order.GetOrder(idOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if strconv.Itoa(dataOrder.IDToko) != idToko {
			http.Error(w, "Anda tidak memiliki otoritas pada pesanan ini.", http.StatusForbidden)
			return
		} else if r.Method != http.MethodGet && dataOrder.Invoice.StatusPesanan == "selesai" {
			http.Error(w, "Pesanan sudah diselesaikan.", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// PenanganOrder is middleware
func (mw MiddleWare) PenanganOrder(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idOrder := vars["idOrder"]

		var order models.Order
		var toko models.Toko

		user := context.Get(r, "user").(*MyClaims)

		dataOrder, _ := order.GetOrder(idOrder)
		dataToko, err := toko.GetToko(strconv.Itoa(dataOrder.IDToko))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if user.IDCustomer != dataToko.IDOwner && dataOrder.Penangan.IDPenangan != user.IDCustomer {
			http.Error(w, "Anda tidak memiliki otoritas pada pesanan ini.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// SetUserPosition is middleware
func (mw MiddleWare) SetUserPosition(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*MyClaims)
		vars := mux.Vars(r)
		idToko := vars["idToko"]
		var toko models.Toko
		var position map[string]interface{}

		dataToko, err := toko.GetToko(idToko)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if dataToko.IDOwner == user.IDCustomer {
			position = map[string]interface{}{
				"position": "owner",
			}
		} else {
			var karyawan models.Karyawan
			dataKaryawan, err := karyawan.GetKaryawanByIDCustomer(idToko, strconv.Itoa(user.IDCustomer))
			if dataKaryawan.Status != "aktif" || err != nil {
				http.Error(w, "Anda tidak memiliki hak akses.", http.StatusForbidden)
				return
			}
			position = map[string]interface{}{
				"position":   dataKaryawan.Posisi,
				"idKaryawan": strconv.Itoa(dataKaryawan.IDKaryawan),
			}
		}

		context.Set(r, "position", position)
		next.ServeHTTP(w, r)
	}
}
