package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// RajaOngkir is class
type RajaOngkir struct {
}

const baseURL = "https://api.rajaongkir.com/starter"
const apiKey = "1999918691cd6b4137c9ec218333e3e2"

// GetIDKota is func
func (rj RajaOngkir) GetIDKota(kota string) (string, bool) {
	url := baseURL + "/city"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", false
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	idKota := gjson.Get(string(body), "rajaongkir.results.#(city_name=="+kota+").city_id")
	if idKota.String() == "" {
		return "", false
	}
	return idKota.String(), true
}

// GetOngkir is func
func (rj RajaOngkir) GetOngkir(asal, tujuan, kodeKurir, service, berat string) (int, string, string, bool) {
	url := baseURL + "/cost"
	payload := strings.NewReader("origin=" + asal + "&destination=" + tujuan + "&weight=" + berat + "&courier=" + kodeKurir)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("key", apiKey)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", "", false
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	ongkir := gjson.Get(string(body), "rajaongkir.results.0.costs.#(service=="+service+").cost.0.value").Int()
	kurir := gjson.Get(string(body), "rajaongkir.results.#(code=="+kodeKurir+").name").String()
	estimasi := gjson.Get(string(body), "rajaongkir.results.0.costs.#(service=="+service+").cost.0.etd").String()

	if estimasi == "" {
		return 0, "", "", false
	}

	return int(ongkir), estimasi, kurir, true
}
