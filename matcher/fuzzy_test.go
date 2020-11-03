package matcher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	. "github.com/gnames/gnmatcher/matcher"
)

var _ = Describe("Fuzzy", func() {
	Describe("MatchFuzzy", func() {
		It("skips matches with edit distance 3 or more", func() {
			ns := NameString{ID: "123", Name: "Pardosa maesta"}
			m := Matcher{FuzzyMatcher: fuzzyMatcherMock{}}
			res := m.MatchFuzzy("Pardosa maesta", "Pardosa maest", ns)
			Expect(len(res.MatchItems)).To(Equal(1))
			Expect(res.MatchItems[0].EditDistance).To(Equal(1))
			ns = NameString{ID: "124", Name: "Acacia may"}
			res = m.MatchFuzzy("Acacia may", "Acacia may", ns)
			Expect(res).To(BeNil())
		})
	})
})

var matchStemMock = map[string][]string{
	"Pardosa maest": {"Pardosa moest"},
	"Acacia may":    {"Acacia ma"},
}

var stemToMatchItemsMock = map[string][]mlib.MatchItem{
	"Pardosa moest": {{ID: "123", MatchStr: "Pardosa moesta"}},
	"Acacia ma":     {{ID: "124", MatchStr: "Acacia maustica"}},
}

type fuzzyMatcherMock struct{}

func (fuzzyMatcherMock) MatchStem(stem string, maxED int) []string {
	if stems, ok := matchStemMock[stem]; ok {
		return stems
	}
	return []string{}
}

func (fuzzyMatcherMock) StemToMatchItems(stem string) []mlib.MatchItem {
	res := []mlib.MatchItem{}
	if mis, ok := stemToMatchItemsMock[stem]; ok {
		return mis
	}
	return res
}
