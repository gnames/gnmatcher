package matcher

import (
	"testing"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/stretchr/testify/assert"
)

// TestFuzzyLimit checks that edit distances larger than 2 are ignored.
func TestFuzzyLimit(t *testing.T) {
	ns := nameString{ID: "123", Name: "Pardosa maesta"}
	m := matcher{fuzzyMatcher: fuzzyMatcherMock{}}
	res := m.matchFuzzy("Pardosa maesta", "Pardosa maest", ns)
	assert.Equal(t, len(res.MatchItems), 1)
	assert.Equal(t, res.MatchItems[0].EditDistance, 1)
	ns = nameString{ID: "124", Name: "Acacia may"}
	res = m.matchFuzzy("Acacia may", "Acacia may", ns)
	assert.Nil(t, res)
}

var matchStemMock = map[string][]string{
	"Pardosa maest": {"Pardosa moest"},
	"Acacia may":    {"Acacia ma"},
}

var stemToMatchItemsMock = map[string][]mlib.MatchItem{
	// EditDistance 2
	"Pardosa moest": {{ID: "123", MatchStr: "Pardosa moesta"}},
	// EditDistance 6
	"Acacia ma": {{ID: "124", MatchStr: "Acacia maustica"}},
}

type fuzzyMatcherMock struct{}

func (fuzzyMatcherMock) Init() {}

func (fuzzyMatcherMock) MatchStem(stem string) []string {
	if stems, ok := matchStemMock[stem]; ok {
		return stems
	}
	return []string{}
}

func (fuzzyMatcherMock) MatchStemExact(stem string) bool {
	return true
}

func (fuzzyMatcherMock) StemToMatchItems(stem string) []mlib.MatchItem {
	res := []mlib.MatchItem{}
	if mis, ok := stemToMatchItemsMock[stem]; ok {
		return mis
	}
	return res
}
