package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gogna/gnparser"
)

// TestNameString checks that name strings are created properly.
func TestNameString(t *testing.T) {
	parser := gnparser.NewGNparser()
	names := []string{
		"Pomatomus",
		"Homo habilis L. Leakey, Tobias & Napier, 1964",
		"Pabstia viridis var. parviflora (Hoehne) Garay",
		"Abedus deinostoma herberti utahensis Menke 1960",
	}
	nstrings := make([]nameString, len(names))
	for i, name := range names {
		ns, _ := newNameString(parser, name)
		nstrings[i] = ns
	}

	uni := nstrings[0].Partial
	assert.Nil(t, uni)

	bi := nstrings[1].Partial
	assert.Equal(t, bi.Genus, "Homo")
	assert.Nil(t, bi.Multinomials)

	tri := nstrings[2].Partial
	assert.Equal(t, tri.Genus, "Pabstia")
	assert.Equal(t, len(tri.Multinomials), 1)
	assert.Equal(t, tri.Multinomials[0].Tail, "Pabstia parviflora")
	assert.Equal(t, tri.Multinomials[0].Head, "Pabstia viridis")

	quat := nstrings[3].Partial
	assert.Equal(t, quat.Genus, "Abedus")
	assert.Equal(t, len(quat.Multinomials), 2)
	assert.Equal(t, quat.Multinomials[0].Tail, "Abedus utahensis")
	assert.Equal(t, quat.Multinomials[0].Head, "Abedus deinostoma herberti")
}
