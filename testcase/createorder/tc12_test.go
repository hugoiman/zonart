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

// soft copy & createOrder error (kolom rencana pakai terlalu panjang)
func Test_TestCase12(t *testing.T) {
	body := map[string]interface{}{
		"jenisPesanan":  "soft copy",
		"tambahanWajah": 2,
		"catatan":       "asdf asdfnjk asdfjnk",
		"pcs":           1,
		"rencanaPakai":  "asdfghjkl asdfghjkl asdfghjkl asdfghjkl asdfghjkl a",
		"gambar":        "tes1.jpg",
		"opsiOrder": []map[string]interface{}{
			{
				"idGrupOpsi": 3,
				"idOpsi":     9,
				"opsi":       "",
			},
		},
	}

	payload, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/api/order/13/1", bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	// set parameter idToko di URL
	vars := map[string]string{
		"idToko":   "13",
		"idProduk": "1",
	}

	request = mux.SetURLVars(request, vars)

	handler := http.HandlerFunc(order.CreateOrder)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 3, Username: "asdf"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v", response.Body)

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya soft copy & createOrder error")
}
