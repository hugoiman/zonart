package rajaongkir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TestCase2(t *testing.T) {
	_, err := rj.GetIDKota("Jakarta Pusat")
	t.Logf("result message:  %v", err)
	assert.Nil(t, err, "seharusnya tidak error")
}
