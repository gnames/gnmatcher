package matcher

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/ent/fuzzy"
)

func (m matcher) matchStem(ns nameString) *mlib.Match {
	matches := m.exactStemMatches(ns.CanonicalStemID, ns.CanonicalStem)
	if len(matches) == 0 {
		return nil
	}
	matchType := vlib.Fuzzy
	matchItems := make([]mlib.MatchItem, 0, len(matches))
	for _, v := range matches {
		if v.MatchStr == ns.Canonical {
			v.MatchType = vlib.Exact
			matchType = vlib.Exact
		} else {
			editDistance := fuzzy.EditDistance(ns.Name, v.MatchStr)
			// editDistance went over threshold
			if editDistance == -1 {
				continue
			}
			v.EditDistance = editDistance
			v.MatchType = vlib.Fuzzy
		}
		matchItems = append(matchItems, v)
	}
	return &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  matchType,
		MatchItems: matchItems,
	}
}

// matchVirus returns the "virus" name the way it was given, without matching.
func (m matcher) matchVirus(ns nameString) *mlib.Match {
	matchItems := m.virusMatcher.MatchVirus(ns.Name)
	matchType := vlib.Virus
	if len(matchItems) == 0 {
		matchType = vlib.NoMatch
	}
	return &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  matchType,
		MatchItems: matchItems,
	}
}
