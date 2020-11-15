package fuzzy_test

import (
	"fmt"
	"testing"

	"github.com/gnames/levenshtein/entity/editdist"
	"github.com/stretchr/testify/assert"
)

func TestDist(t *testing.T) {
	testData := []struct {
		str1, str2 string
		dist       int
	}{
		{"Hello", "Hello", 0},
		{"Pomatomus", "Pom-tomus", 1},
		{"Pomatomus", "Poma  tomus", 2},
		{"Pomatomus", "Pomщtomus", 1},
		{"sitting", "kitten", 3},
	}

	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		dist, _, _ := editdist.ComputeDistance(v.str1, v.str2, false)
		assert.Equal(t, dist, v.dist, msg)
	}
}

// BenchmarkDist checks the speed of fuzzy matching. Run it with:
// `go test -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt`

func BenchmarkDist(b *testing.B) {
	var out int
	b.Run("CompareOnce", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out, _, _ = editdist.ComputeDistance("Pomatomus solatror", "Pomatomus saltator", false)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
}
