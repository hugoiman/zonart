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

// min > max
func Test_TC7(t *testing.T) {
	body := map[string]interface{}{
		"namaGrup":        "Pakaian",
		"required":        true,
		"min":             2,
		"max":             1,
		"spesificRequest": true,
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

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya minimal jumlah memilih > maksimal memilih")
}
