package creategrupopsi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	mw "zonart/middleware"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

// spesificRequest = false, req = false && min != 0
func Test_TC3(t *testing.T) {
	body := map[string]interface{}{
		"namaGrup":        "Pakaian",
		"required":        false,
		"min":             1,
		"max":             2,
		"spesificRequest": false,
		"opsi": []map[string]interface{}{
			{
				"opsi":      "Kemeja",
				"harga":     0,
				"berat":     0,
				"perProduk": false,
				"status":    true,
			},
			{
				"opsi":      "Jas",
				"harga":     0,
				"berat":     0,
				"perProduk": false,
				"status":    true,
			},
		},
	}

	payload, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/grup-opsi/13", bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(grupOpsi.CreateGrupOpsi)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 3, Username: "asdf"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v", response.Body)

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya required false && minimal memilih != 0")
}
