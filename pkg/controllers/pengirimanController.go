package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// PengirimanController is class
type PengirimanController struct{}

// SetResi is func
func (rc PengirimanController) SetResi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	var dataJSON map[string]interface{}
	json.NewDecoder(r.Body).Decode(&dataJSON)

	resi := fmt.Sprintf("%v", dataJSON["resi"])

	if err := validator.New().Var(resi, "required"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pengiriman models.Pengiriman
	pengiriman.SetResi(resi)
	err := pengiriman.InputResi(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Nomor resi telah disimpan."}`))
}
