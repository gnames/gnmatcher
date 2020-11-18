package fuzzy

import (
	"strings"

	"github.com/gnames/levenshtein/entity/editdist"
)

const (
	charsPerED      = 5
	maxEditDistance = 2
)

// EditDistance calculates edit distance (**ed**) according to Levenshtein algorithm.
// It also runs additional checks and if they fail, returns -1.
//
// Checks:
// - result should not exceed maxEditDistance
// - number of characters divided by ed should be bigger than charsPerED
//
// It assumes that checks have to be applied only to the second string:
//  EditDistance("Pomatomus", "Pom atomus")
// returns -1
//  EditDistance("Pom atomus", "Pomatomus")
// returns 1
//
// It also assumes that number os spaces between words was already
// normalized to 1 space, and that s1 and s2 always have the same number of
// words.
func EditDistance(s1, s2 string) int {
	ed, _, _ := editdist.ComputeDistance(s1, s2, false)
	if ed == 0 {
		return ed
	}

	if ed > maxEditDistance {
		return -1
	}
	return checkED(s1, s2, ed)
}

func checkED(s1, s2 string, ed int) int {
	words1 := strings.Split(s1, " ")
	words2 := strings.Split(s2, " ")
	if len(words1) != len(words2) {
		return ed
	}
	for i, w := range words2 {
		r := []rune(w)
		if len(r)/ed < charsPerED {
			wordED, _, _ := editdist.ComputeDistance(w, words1[i], false)
			if wordED > 0 && len(r)/wordED < charsPerED {
				return -1
			}
		}
	}
	return ed
}
