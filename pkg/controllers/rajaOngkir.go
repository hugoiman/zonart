package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"gopkg.in/go-playground/validator.v9"
)

// RajaOngkir is class
type RajaOngkir struct {
	baseURL string
	apiKey  string
}

// setVariable is setter
func (rj *RajaOngkir) setVariable() {
	rj.baseURL = os.Getenv("RAJAONGKIR_BASE_URL")
	rj.apiKey = os.Getenv("RAJAONGKIR_API_KEY")
}

// GetIDKota is func
func (rj RajaOngkir) GetIDKota(kota string) (string, error) {
	rj.setVariable()
	uRL := rj.baseURL + "/city"
	req, _ := http.NewRequest("GET", uRL, nil)
	req.Header.Add("key", rj.apiKey)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	idKota := gjson.Get(string(body), "rajaongkir.results.#(city_name=="+kota+").city_id")
	if idKota.String() == "" {
		return "", errors.New("Kota tidak ditemukan")
	}
	return idKota.String(), nil
}

// GetOngkir is func
func (rj RajaOngkir) GetOngkir(asal, tujuan, kodeKurir, service string, berat int) (int, string, string, error) {
	rj.setVariable()
	kota_asal, err := rj.GetIDKota(asal)
	if err != nil {
		return 0, "", "", err
	}

	kota_tujuan, err := rj.GetIDKota(tujuan)
	if err != nil {
		return 0, "", "", err
	}

	uRL := rj.baseURL + "/cost"
	payload := strings.NewReader("origin=" + kota_asal + "&destination=" + kota_tujuan + "&weight=" + strconv.Itoa(berat) + "&courier=" + kodeKurir)
	req, _ := http.NewRequest("POST", uRL, payload)

	req.Header.Add("key", rj.apiKey)
	req.Header.Add("content-type", "application/x-www-form-uRLencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	ongkir := gjson.Get(string(body), "rajaongkir.results.0.costs.#(service=="+service+").cost.0.value").Int()
	kurir := gjson.Get(string(body), "rajaongkir.results.#(code=="+kodeKurir+").name").String()
	estimasi := gjson.Get(string(body), "rajaongkir.results.0.costs.#(service=="+service+").cost.0.etd").String()

	if estimasi == "" {
		return 0, "", "", err
	}

	var regex = regexp.MustCompile(`(?i)hari|jam`)
	if !regex.MatchString(estimasi) {
		estimasi += " hari"
	}
	subStr := strings.NewReplacer("JAM", "jam", "HARI", "hari")
	estimasi = subStr.Replace(estimasi)

	return int(ongkir), estimasi, kurir, nil
}

// GetAllKota is func
func (rj RajaOngkir) GetAllKota(w http.ResponseWriter, r *http.Request) {
	rj.setVariable()
	url := rj.baseURL + "/city"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("key", rj.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (rj RajaOngkir) GetCost(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Origin      string `json:"origin"  validate:"required"`
		Destination string `json:"destination" validate:"required"`
		Weight      string `json:"weight" validate:"required"`
		Courier     string `json:"courier" validate:"required"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if err := validator.New().Struct(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rj.setVariable()
	url := rj.baseURL + "/cost"

	payload := strings.NewReader("origin=" + data.Origin + "&destination=" + data.Destination + "&weight=" + data.Weight + "&courier=" + data.Courier)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("key", rj.apiKey)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
