package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

// CreateRevisi is func
func CreateRevisi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrder := vars["idOrder"]
	var revisi models.Revisi

	if err := json.NewDecoder(r.Body).Decode(&revisi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(revisi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := revisi.CreateRevisi(idOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send notif to penangan

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Sukses! Revisi telah terkirim."}`))
}
