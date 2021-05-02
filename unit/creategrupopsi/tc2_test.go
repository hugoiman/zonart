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

// gagal validasi
func Test_TC2(t *testing.T) {
	body := map[string]interface{}{
		"namaGrup":        "",
		"required":        false,
		"min":             0,
		"max":             1,
		"spesificRequest": true,
		"opsi": []map[string]interface{}{
			{
				"opsi":      "",
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

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya gagal validasi")
}
