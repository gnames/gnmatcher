package fuzzy_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/gnames/gnmatcher/internal/ent/fuzzy"
	"github.com/stretchr/testify/assert"
)

// EditDist without constraints
func TestDist(t *testing.T) {
	// to hide warnings
	oldLevel := slog.SetLogLoggerLevel(10)
	defer slog.SetLogLoggerLevel(oldLevel)

	testData := []struct {
		str1, str2 string
		dist       int
	}{
		{"Hello", "Hello", 0},
		{"Pomatomus", "Pom-tomus", 1},
		{"Pomatomus", "Pomщtomus", 1},
		{"Pom atomus", "Poma tomus", 2},
		// ed = 3, it was not allowed before, it is now, because
		// we assume that high edit distance comes from suffix of a stemmed
		// match
		{"sitting", "kitten", 3},
		// words are too small
		{"Acacia mal", "Acacia may", -1},
		// differnt number of words is not covered yet
		{"Pomatomus", "Poma  tomus", 2},
		// edge cases that should not happen
		// more than one empty space
		{"Pomatomus saltator", "Pomatomus  saltator", 1},
		{"Vesicaria creticum", "Vesicaria cretica", 2},
	}

	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		dist := fuzzy.EditDistance(v.str1, v.str2, false)
		assert.Equal(t, v.dist, dist, msg)
	}
}

// BenchmarkDist checks the speed of fuzzy matching. Run it with:
// `go test -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt`

func BenchmarkDist(b *testing.B) {
	var out int
	b.Run("CompareOnce", func(b *testing.B) {
		for range b.N {
			out = fuzzy.EditDistance("Pomatomus solatror", "Pomatomus saltator", false)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
}
