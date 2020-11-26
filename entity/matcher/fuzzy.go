package matcher

import (
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnmatcher/entity/fuzzy"
)

// matchFuzzy tries to get fuzzy matching of a stemmed name-string to canonical
// forms from the gnames database.
func (m matcher) matchFuzzy(name, stem string,
	ns nameString) *mlib.Match {
	stemMatches := m.fuzzyMatcher.MatchStem(stem)
	if len(stemMatches) == 0 {
		return nilResult
	}

	res := &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  vlib.Fuzzy,
		MatchItems: make([]mlib.MatchItem, 0, len(stemMatches)*2),
	}
	for _, stemMatch := range stemMatches {
		editDistanceStem := fuzzy.EditDistance(stemMatch, stem)
		if editDistanceStem == -1 {
			continue
		}
		matchItems := m.fuzzyMatcher.StemToMatchItems(stemMatch)
		for _, matchItem := range matchItems {
			// runs edit distance with checks, returns -1 if checks failed.
			editDistance := fuzzy.EditDistance(name, matchItem.MatchStr)
			// skip matches that failed edit distance checks.
			if editDistance == -1 {
				continue
			}
			matchItem.EditDistance = editDistance
			matchItem.EditDistanceStem = editDistanceStem
			res.MatchItems = append(res.MatchItems, matchItem)
		}
	}
	if len(res.MatchItems) == 0 {
		return nil
	}
	return res
}
