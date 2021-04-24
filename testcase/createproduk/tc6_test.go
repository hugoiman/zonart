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

func Test_TestCase6(t *testing.T) {
	body := map[string]interface{}{
		"namaProduk": "Bantal Karikatur",
		"berat":      500,
		"status":     "aktif",
		"hargaWajah": 75000,
		"jenisPemesanan": []map[string]interface{}{
			{
				"idJenisPemesanan": 1,
				"harga":            90000,
				"status":           true,
			},
			{
				"idJenisPemesanan": 2,
				"harga":            0,
				"status":           false,
			},
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

	file := "./test.jpg"
	fw, err := w.CreateFormFile("gambar", file)
	if err != nil {
		t.Error(err)
	}
	fd, err := os.Open(file)
	if err != nil {
		t.Error(err)
	}
	defer fd.Close()

	_, err = io.Copy(fw, fd)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	request, _ := http.NewRequest(http.MethodPost, "/produk/idToko", buffer)
	request = mux.SetURLVars(request, map[string]string{"idToko": "37"})
	request.Header.Set("Content-Type", w.FormDataContentType())
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(produk.CreateProduk)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 13, Username: "geno"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v\n status code: %v", response.Body, response.Result().StatusCode)

	assert.Equal(t, response.Code, http.StatusOK, "Seharusnya berhasil menyimpan")
}
