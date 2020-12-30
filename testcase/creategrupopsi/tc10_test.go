package creategrupopsi

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

// opsi gagal disimpan di basis data (char = 51)
func Test_TC10(t *testing.T) {
	body := map[string]interface{}{
		"namaGrup":        "Pakaians",
		"required":        true,
		"min":             1,
		"max":             2,
		"spesificRequest": true,
		"opsi": []map[string]interface{}{
			{
				"opsi":      "asdfghjkl asdfghjkl asdfghjkl asdfghjkl asdfghjkl a",
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
	request, _ := http.NewRequest(http.MethodPost, "/api/grup-opsi/13", bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	// set parameter idToko di URL
	vars := map[string]string{
		"idToko": "13",
	}

	request = mux.SetURLVars(request, vars)

	handler := http.HandlerFunc(grupOpsi.CreateGrupOpsi)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 3, Username: "asdf"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v", response.Body)

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya gagal menyimpan opsi pada database")
}
