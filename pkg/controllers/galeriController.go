package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"zonart/pkg/models"

	"github.com/gorilla/mux"
)

// GaleriController is class
type GaleriController struct{}

// GetGaleris is func
func (gc GaleriController) GetGaleris(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var galeri models.Galeri

	dataGaleri := galeri.GetGaleris(idToko)
	message, _ := json.Marshal(&dataGaleri)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

// CreateGaleri is func
func (gc GaleriController) CreateGaleri(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	var galeri models.Galeri

	if err := json.NewDecoder(strings.NewReader(payload)).Decode(&galeri); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, _, err := r.FormFile("gambar"); err == http.ErrMissingFile {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	maxSize := int64(1024 * 1024 * 2) // 2 MB
	destinationFolder := "zonart/galeri"
	var cloudinary Cloudinary
	images, err := cloudinary.UploadImages(r, maxSize, destinationFolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var idGaleris []string
	for _, v := range images {
		galeri.SetGambar(v)
		idGaleri, err := galeri.CreateGaleri(idToko)
		if err != nil {
			cloudinary.DeleteImages(images)
			_ = galeri.DeleteGaleri(idToko, idGaleris)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		idGaleris = append(idGaleris, strconv.Itoa(idGaleri))
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Gambar telah ditambahkan."}`))
}

// DeleteGaleri is func
func (gc GaleriController) DeleteGaleri(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idToko := vars["idToko"]
	idGaleri := vars["idGaleri"]
	idGaleris := []string{idGaleri}
	var galeri models.Galeri

	err := galeri.DeleteGaleri(idToko, idGaleris)
	if err != nil {
		http.Error(w, "Data tidak ditemukan.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data telah dihapus!"}`))

}
