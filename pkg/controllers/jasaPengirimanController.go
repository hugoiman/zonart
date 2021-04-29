package controllers

import (
	"encoding/json"
	"net/http"
	"zonart/pkg/models"
)

// JasaPengirimanController is class
type JasaPengirimanController struct{}

// GetJasaPengirimans is func
func (jpc JasaPengirimanController) GetJasaPengirimans(w http.ResponseWriter, r *http.Request) {
	var jp models.JasaPengiriman

	dataJasaPengiriman := jp.GetJasaPengirimans()
	message, _ := json.Marshal(&dataJasaPengiriman)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}
