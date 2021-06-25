package rajaongkir

import (
	"testing"
	"time"
	"zonart/pkg/models"

	"github.com/stretchr/testify/assert"
)

func Test_TestCase2(t *testing.T) {
	var toko models.Toko
	toko.SetOwner(13)
	toko.SetNamaToko("Creative Art")
	toko.SetEmailToko("creativeart@gmail.com")
	toko.SetFoto("https://res.cloudinary.com/dbddhr9rz/image/upload/v1612894274/zonart/toko/toko_jhecxf.png")
	toko.SetKota("Jakarta Timur")
	toko.SetSlug("creative-art")
	toko.SetCreatedAt(time.Now().Format("2006-01-02"))
	_, err := toko.CreateToko()

	t.Logf("result message:  %v", err)
	assert.Nil(t, err, "seharusnya tidak error")
}
