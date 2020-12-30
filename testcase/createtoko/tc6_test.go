package createtoko

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

func Test_TestCase6(t *testing.T) {
	body := map[string]interface{}{
		"namaToko":  "Barokart",
		"emailToko": "barokart@gmail.com",
		"alamat":    "Jl. Cindera Mata no 53.",
		"kota":      "Jakarta Utara",
		"whatsapp":  "081234567898",
		"slug":      "barok-art",
	}
	payload, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/toko", bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(toko.CreateToko)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 3, Username: "asdf"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v", response.Body)

	assert.Equal(t, response.Code, http.StatusOK, "Sukses!")
}
