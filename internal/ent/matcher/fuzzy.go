package matcher

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/internal/ent/fuzzy"
)

// matchFuzzy tries to get fuzzy matching of a stemmed name-string to canonical
// forms from the gnames database.
func (m matcher) matchFuzzy(
	canonical,
	stem string,
	ns nameString,
) (*mlib.Match, error) {
	relax := m.cfg.WithRelaxedFuzzyMatch

	matchType := vlib.Fuzzy
	if relax {
		matchType = vlib.FuzzyRelaxed
	}

	stemMatches := m.fuzzyMatcher.MatchStem(stem)
	if len(stemMatches) == 0 {
		return nil, nil
	}

	res := &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  matchType,
		MatchItems: make([]mlib.MatchItem, 0, len(stemMatches)*2),
	}

	for _, stemMatch := range stemMatches {
		editDistanceStem := fuzzy.EditDistance(stemMatch, stem, relax)
		// -1 means edit distance got over threshold
		if editDistanceStem == -1 {
			continue
		}
		matchItems, err := m.fuzzyMatcher.StemToMatchItems(stemMatch)
		if err != nil {
			return nil, err
		}

		for _, matchItem := range matchItems {
			matchItem.InputStr = canonical
			// runs edit distance with checks, returns -1 if checks failed.
			editDistance := fuzzy.EditDistance(
				matchItem.InputStr,
				matchItem.MatchStr,
				relax,
			)
			// skip matches that failed edit distance checks.
			if editDistance == -1 {
				continue
			}

			matchItem.EditDistance = editDistance
			matchItem.EditDistanceStem = editDistanceStem
			matchItem.MatchType = matchType
			res.MatchItems = append(res.MatchItems, matchItem)
		}
	}

	res.MatchItems = m.filterDataSources(res.MatchItems)
	if len(res.MatchItems) == 0 {
		return nil, nil
	}

	return res, nil
}
