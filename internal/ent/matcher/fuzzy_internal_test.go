package matcher

import (
	"testing"

	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TestFuzzyLimit checks that edit distances larger than 2 are ignored.
func TestFuzzyLimit(t *testing.T) {
	ns := nameString{ID: "123", Name: "Pardosa maesta"}
	m := matcher{fuzzyMatcher: fuzzyMatcherMock{}}
	res, err := m.matchFuzzy("Pardosa maesta", "Pardosa maest", ns)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res.MatchItems))
	assert.Equal(t, 1, res.MatchItems[0].EditDistance)
	ns = nameString{ID: "124", Name: "Acacia may"}
	res, err = m.matchFuzzy("Acacia may", "Acacia may", ns)
	assert.Nil(t, err)
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

func (fuzzyMatcherMock) Init() error { return nil }

func (fuzzyMatcherMock) SetConfig(cfg config.Config) {}

func (fuzzyMatcherMock) MatchStem(stem string) []string {
	if stems, ok := matchStemMock[stem]; ok {
		return stems
	}
	return []string{}
}

func (fuzzyMatcherMock) MatchStemExact(stem string) bool {
	return true
}

func (fuzzyMatcherMock) StemToMatchItems(
	stem string,
) ([]mlib.MatchItem, error) {
	res := []mlib.MatchItem{}
	if mis, ok := stemToMatchItemsMock[stem]; ok {
		return mis, nil
	}
	return res, nil
}
