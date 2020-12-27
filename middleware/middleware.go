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
			http.Error(w, "Gagal! Dibutuhkan otentikasi. Silahkan melakukan login.", http.StatusUnauthorized)
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
			http.Error(w, "Gagal! Anda tidak memiliki hak akses.", http.StatusForbidden)
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

		dataKaryawan, err := karyawan.AuthKaryawan(idToko, strconv.Itoa(user.IDCustomer))
		if err != nil || (user.IDCustomer != dataToko.IDOwner && (dataKaryawan.Posisi != "admin" || dataKaryawan.Status != "aktif")) {
			http.Error(w, "Gagal! Anda tidak memiliki hak akses.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// AuthCustomerOrder is middleware
func (mw MiddleWare) AuthCustomerOrder(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idOrder := vars["idOrder"]
		var order models.Order

		user := context.Get(r, "user").(*MyClaims)

		dataOrder, err := order.GetOrder(idOrder, strconv.Itoa(user.IDCustomer))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if user.IDCustomer != dataOrder.IDCustomer {
			http.Error(w, "Gagal! Anda tidak memiliki otoritas pada pesanan ini.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
