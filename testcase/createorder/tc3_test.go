package createorder

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	mw "zonart/middleware"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// produk not found
func Test_TestCase3(t *testing.T) {
	body := map[string]interface{}{
		"jenisPesanan":  "cetak",
		"tambahanWajah": 2,
		"catatan":       "asdf asdfnjk asdfjnk",
		"pcs":           1,
		"rencanaPakai":  "5 Januari 2021",
		"gambar":        "tes1.jpg",
		"opsiOrder":     []map[string]interface{}{},
	}

	payload, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/api/order/13/1000", bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	// set parameter idToko di URL
	vars := map[string]string{
		"idToko":   "13",
		"idProduk": "1000",
	}

	request = mux.SetURLVars(request, vars)

	handler := http.HandlerFunc(order.CreateOrder)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 3, Username: "asdf"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v", response.Body)

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya produk tidak ditemukan")
}
