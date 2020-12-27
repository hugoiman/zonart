package controllers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	mw "zonart/middleware"
	"zonart/pkg/models"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v9"
)

// MyClaims is credential
type MyClaims = mw.MyClaims

// AuthController is class
type AuthController struct{}

// Login is func
func (auth AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var login models.Auth
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sha = sha1.New()
	sha.Write([]byte(login.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	login.Password = encryptedString

	idCustomer, err := login.Login()
	if err != nil {
		http.Error(w, "Gagal! Username atau password salah.", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	data, err := customer.GetCustomer(idCustomer)
	if err != nil {
		http.Error(w, "Gagal! Terjadi error.", http.StatusInternalServerError)
		return
	}

	token := AuthController{}.CreateToken(data)

	type M map[string]interface{}
	message, _ := json.Marshal(M{"token": token})

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateToken is Generate token
func (auth AuthController) CreateToken(customer models.Customer) string {
	var mySigningKey = mw.MySigningKey
	claims := MyClaims{
		IDCustomer: customer.IDCustomer,
		Username:   customer.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

// ResetPassword is func
func (auth AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var dataJSON map[string]interface{}
	json.NewDecoder(r.Body).Decode(&dataJSON)

	email := fmt.Sprintf("%v", dataJSON["email"])
	var customer models.Customer

	if err := validator.New().Var(email, "required,email"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := customer.GetCustomer(email)
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusBadRequest)
		return
	}

	// Generate Random String
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randomStr := make([]rune, 10)
	for i := range randomStr {
		randomStr[i] = letters[rand.Intn(len(letters))]
	}
	newPass := string(randomStr)

	var pass = sha1.New()
	pass.Write([]byte(newPass))
	var encryptedPass = fmt.Sprintf("%x", pass.Sum(nil))

	_ = customer.UpdatePassword(data.IDCustomer, encryptedPass)

	// message := "Hallo " + data.Nama + ", your new password is <b>" + newPass + "</b>"
	// err := SendEmail("New Password", email, message)
	// if err != nil {
	// 	http.Error(w, "Gagal! Coba beberapa saat lagi.", http.StatusBadRequest)
	// 	return
	// }

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"New password has been sent to your email ` + newPass + `."}`))
}
