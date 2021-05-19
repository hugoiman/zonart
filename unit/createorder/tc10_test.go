package createtoko

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	mw "zonart/middleware"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_TestCase10(t *testing.T) {
	body := map[string]interface{}{
		"jenisPesanan":  "cetak",
		"tambahanWajah": 2,
		"pcs":           2,
		"rencanaPakai":  "24 November 2021",
		"opsiOrder": []map[string]interface{}{
			{
				"idGrupOpsi": 32,
				"idOpsi":     0,
				"opsi":       "kue ultah, terompet, lilin",
			},
			{
				"idGrupOpsi": 33,
				"idOpsi":     36,
			},
			{
				"idGrupOpsi": 30,
				"idOpsi":     29,
			},
		},
		"pengiriman": map[string]interface{}{
			"penerima":  "Roy",
			"telp":      "08123456",
			"alamat":    "Jl. ikan no 23",
			"kota":      "Jakarta Timur",
			"label":     "Rumah",
			"kodeKurir": "tiki",
			"service":   "ECO",
		},
	}

	payload, _ := json.Marshal(body)
	buffer := new(bytes.Buffer)
	w := multipart.NewWriter(buffer)
	data, err := w.CreateFormField("payload")
	if err != nil {
		t.Error(err)
	}
	data.Write(payload)

	files := []string{"./avatar-1.png", "./avatar-2.png"}
	for _, file := range files {
		fw, err := w.CreateFormFile("fileOrder", file)
		if err != nil {
			t.Error(err)
		}
		fd, err := os.Open(file)
		if err != nil {
			t.Error(err)
		}

		_, err = io.Copy(fw, fd)
		if err != nil {
			t.Error(err)
		}

		fd.Close()
	}

	w.Close()

	request, _ := http.NewRequest(http.MethodPost, "/order/idToko/idProduk", buffer)
	request = mux.SetURLVars(request, map[string]string{"idToko": "37", "idProduk": "10"})
	request.Header.Set("Content-Type", w.FormDataContentType())
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(order.CreateOrder)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 11, Username: "roy"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v\n status code: %v", response.Body, response.Result().StatusCode)

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya gagal mendecode")
}
