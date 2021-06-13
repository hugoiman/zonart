package controllers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"zonart/custerr"
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
func (ac AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth
	var login = struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(login); err != nil {
		http.Error(w, custerr.CustomError(err).Error(), http.StatusBadRequest)
		return
	}

	var sha = sha1.New()
	sha.Write([]byte(login.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	login.Password = encryptedString

	idCustomer, err := auth.Login(login.Username, login.Password)
	if err != nil {
		http.Error(w, "Username atau password salah.", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	data, _ := customer.GetCustomer(idCustomer)

	token := AuthController{}.createToken(data)

	type M map[string]interface{}
	message, _ := json.Marshal(M{"token": token})

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// createToken is Generate token
func (ac AuthController) createToken(customer models.Customer) string {
	var mdw mw.MiddleWare
	var mySigningKey = mdw.GetSigningKey()
	claims := MyClaims{
		IDCustomer: customer.GetIDCustomer(),
		Username:   customer.GetUsername(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

// ResetPassword is func
func (ac AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Email       string `json:"email"  validate:"required,email"`
		NewPassword string `json:"newPassword" validate:"required,min=6"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(data); err != nil {
		http.Error(w, custerr.CustomError(err).Error(), http.StatusBadRequest)
		return
	}

	var customer models.Customer
	dataCustomer, err := customer.GetCustomer(data.Email)
	if err != nil {
		http.Error(w, "Email tidak terdaftar", http.StatusBadRequest)
		return
	}

	// Generate token
	var mdw mw.MiddleWare
	var mySigningKey = mdw.GetSigningKey()
	claims := ClaimsResetPassword{
		Email:    data.Email,
		Password: data.NewPassword,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)
	link := "http://localhost:5500/verifikasi-reset-password?token=" + tokenString

	message := "Hallo " + dataCustomer.GetNama() + ",<br> Silahkan klik link dibawah ini untuk verifikasi. Link ini akan kadaluarsa dalam waktu 15 menit.<br><br>" + link
	var gomail Gomail
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
func (ac AuthController) VerificationResetPassword(w http.ResponseWriter, r *http.Request) {
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
		var mdw mw.MiddleWare
		return mdw.GetSigningKey(), nil
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
	_ = customer.UpdatePassword(dataCustomer.GetIDCustomer(), encryptedPassword)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Berhasil diverifikasi. Silahkan melakukan login."}`))
}

// Logout is func
func (ac AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Anda telah keluar dari sistem."}`))
}
