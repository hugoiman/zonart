package createtoko

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	mw "zonart/middleware"
	"zonart/pkg/controllers"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func Test_TestCase4(t *testing.T) {
	body := map[string]interface{}{
		"namaToko":  "Barokart",
		"emailToko": "barokart@gmail.com",
		"alamat":    "Jl. Cindera Mata no 53.",
		"kota":      "Jakarta Utara",
		"whatsapp":  "081234567898",
		"slug":      "geno-art", // domain tidak unik
	}
	payload, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "/toko", bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	var toko controllers.TokoController
	handler := http.HandlerFunc(toko.CreateToko)

	// set identitas user
	context.Set(request, "user", &mw.MyClaims{IDCustomer: 13, Username: "geno"})

	handler.ServeHTTP(response, request)
	t.Logf("response message:  %v\n status code: %v", response.Body, response.Result().StatusCode)

	assert.Equal(t, response.Code, http.StatusBadRequest, "Seharusnya slug tidak unik!")
}
