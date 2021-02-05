package controllers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"zonart/internal/gomail"
	mw "zonart/middleware"
	"zonart/pkg/models"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v9"
)

// MyClaims is credential
type MyClaims = mw.MyClaims

// ClaimsResetPassword is Credential
type ClaimsResetPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

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
		http.Error(w, "Username atau password salah.", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	data, err := customer.GetCustomer(idCustomer)
	if err != nil {
		http.Error(w, "Terjadi error.", http.StatusInternalServerError)
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
	data := struct {
		Email       string `json:"email"  validate:"required,email"`
		NewPassword string `json:"newPassword" validate:"required,min=6"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var customer models.Customer
	dataCustomer, err := customer.GetCustomer(data.Email)
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusBadRequest)
		return
	}

	// Generate token
	var mySigningKey = mw.MySigningKey
	claims := ClaimsResetPassword{
		Email:    data.Email,
		Password: data.NewPassword,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)
	link := "http://localhost:5500/pages/verifikasi-reset-password.html?token=" + tokenString

	message := "Hallo " + dataCustomer.Nama + ",<br> Silahkan klik link dibawah ini untuk verifikasi. Link ini akan kadaluarsa dalam waktu 15 menit.<br><br>" + link
	err = gomail.SendEmail("Reset Password", data.Email, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Link verifikasi telah terkirim. Silahkan periksa email anda."}`))
}

// VerificationResetPassword is func
func (auth AuthController) VerificationResetPassword(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Token string `json:"token"  validate:"required"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claimsResetPassword := &ClaimsResetPassword{}

	token, err := jwt.ParseWithClaims(data.Token, claimsResetPassword, func(token *jwt.Token) (interface{}, error) {
		return mw.MySigningKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token invalid: "+err.Error(), http.StatusBadRequest) // Token expired/key tidak cocok(invalid)
		return
	}

	var customer models.Customer
	dataCustomer, err := customer.GetCustomer(claimsResetPassword.Email)
	if err != nil {
		http.Error(w, "User tidak ditemukan", http.StatusBadRequest)
		return
	}

	//  Enkripsi Password
	var pass = sha1.New()
	pass.Write([]byte(claimsResetPassword.Password))
	var encryptedPassword = fmt.Sprintf("%x", pass.Sum(nil))

	//  Update Password
	_ = customer.UpdatePassword(dataCustomer.IDCustomer, encryptedPassword)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Berhasil diverifikasi. Silahkan melakukan login."}`))
}

// Logout is func
