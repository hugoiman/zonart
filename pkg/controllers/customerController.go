package controllers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"zonart/pkg/models"

	"github.com/gorilla/context"
	"gopkg.in/go-playground/validator.v9"
)

// CustomerController is class
type CustomerController struct{}

// GetCustomer is func
func (cc CustomerController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*MyClaims)
	var customer models.Customer

	data, err := customer.GetCustomer(strconv.Itoa(user.IDCustomer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// Register is func
func (cc CustomerController) Register(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer

	data := struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email"  validate:"required,email"`
		Nama     string `json:"nama" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pass = sha1.New()
	pass.Write([]byte(data.Password))
	var encryptedPassword = fmt.Sprintf("%x", pass.Sum(nil))

	err := customer.Register(data.Username, data.Email, data.Nama, encryptedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Registrasi sukses!"}`))

}

// UpdateProfil is func
func (cc CustomerController) UpdateProfil(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	user := context.Get(r, "user").(*MyClaims)

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := customer.UpdateProfil(user.IDCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Profil berhasil diperbarui!"}`))
}

// ChangePassword is func
func (cc CustomerController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	user := context.Get(r, "user").(*MyClaims)

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	if err := validator.New().Var(fmt.Sprintf("%v", data["newPassword"]), "required,min=6,max=18"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Var(fmt.Sprintf("%v", data["oldPassword"]), "required,min=6,max=18"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var oldPass = sha1.New()
	oldPass.Write([]byte(fmt.Sprintf("%v", data["oldPassword"])))
	var encryptedOldPass = fmt.Sprintf("%x", oldPass.Sum(nil))

	isValid := models.CheckOldPassword(user.IDCustomer, encryptedOldPass)
	if !isValid {
		http.Error(w, "Password lama tidak sesuai", http.StatusBadRequest)
		return
	}

	var newPass = sha1.New()
	newPass.Write([]byte(fmt.Sprintf("%v", data["newPassword"])))
	var encryptedPass = fmt.Sprintf("%x", newPass.Sum(nil))

	_ = customer.UpdatePassword(user.IDCustomer, encryptedPass)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Password berhasil diperbarui!"}`))
}
