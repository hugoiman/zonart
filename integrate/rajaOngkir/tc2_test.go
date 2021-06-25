package rajaongkir

import (
	"testing"
	"zonart/pkg/controllers"

	"github.com/stretchr/testify/assert"
)

func Test_TestCase2(t *testing.T) {
	var rj controllers.RajaOngkir
	_, err := rj.GetIDKota("Jakarta Pusat")
	t.Logf("result message:  %v", err)
	assert.Nil(t, err, "seharusnya tidak error")
}
