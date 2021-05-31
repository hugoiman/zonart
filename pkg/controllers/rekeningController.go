package controllers

import (
	"net/http"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// RekeningController is class
type RekeningController struct{}

// DeleteRekening is func
func (rc RekeningController) DeleteRekening(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idRekening := vars["idRekening"]
	var rekening models.Rekening

	_ = rekening.DeleteRekening(idToko, idRekening)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Rekening telah dihapus!"}`))

}
