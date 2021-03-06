package fuzzy_test

import (
	"fmt"
	"testing"

	"github.com/gnames/gnmatcher/ent/fuzzy"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// EditDist without constraints
func TestDist(t *testing.T) {
	// to hide warnings
	log.SetLevel(log.FatalLevel)

	testData := []struct {
		str1, str2 string
		dist       int
	}{
		{"Hello", "Hello", 0},
		{"Pomatomus", "Pom-tomus", 1},
		{"Pomatomus", "Pomщtomus", 1},
		// ed = 3, too big
		{"sitting", "kitten", -1},
		// words are too small
		{"Pom atomus", "Poma tomus", -1},
		{"Acacia mal", "Acacia may", -1},
		// differnt number of words is not covered yet
		{"Pomatomus", "Poma  tomus", 2},
		// edge cases that should not happen
		// more than one empty space
		{"Pomatomus saltator", "Pomatomus  saltator", 1},
	}

	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		dist := fuzzy.EditDistance(v.str1, v.str2)
		assert.Equal(t, dist, v.dist, msg)
	}
}

// BenchmarkDist checks the speed of fuzzy matching. Run it with:
// `go test -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt`

func BenchmarkDist(b *testing.B) {
	var out int
	b.Run("CompareOnce", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out = fuzzy.EditDistance("Pomatomus solatror", "Pomatomus saltator")
		}
		_ = fmt.Sprintf("%d\n", out)
	})
}
