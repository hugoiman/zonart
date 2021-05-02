package rajaongkir

import (
	"testing"
	"zonart/pkg/controllers"

	"github.com/stretchr/testify/assert"
)

var rj controllers.RajaOngkir

func Test_TestCase1(t *testing.T) {
	_, err := rj.GetIDKota("Jakarta Tengah")

	t.Logf("result message:  %v", err)
	assert.NotNil(t, err, "seharusnya terdapat error")
}
