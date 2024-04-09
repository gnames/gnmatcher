package matcher

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/internal/ent/fuzzy"
)

func (m matcher) matchStem(ns nameString) (*mlib.Match, error) {
	matches, err := m.exactStemMatches(ns.CanonicalStemID, ns.CanonicalStem)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, nil
	}
	matchType := vlib.Fuzzy
	matchItems := make([]mlib.MatchItem, 0, len(matches))
	relax := m.cfg.WithRelaxedFuzzyMatch
	for _, v := range matches {
		v.InputStr = ns.Canonical
		if v.MatchStr == v.InputStr {
			v.MatchType = vlib.Exact
			matchType = vlib.Exact
		} else {
			editDistance := fuzzy.EditDistance(v.InputStr, v.MatchStr, relax)
			// editDistance went over threshold
			if editDistance == -1 {
				continue
			}
			v.EditDistance = editDistance
			v.MatchType = vlib.Fuzzy
		}
		matchItems = append(matchItems, v)
	}

	matchItems = m.filterDataSources(matchItems)
	if len(matchItems) == 0 {
		return nil, nil
	}

	res := &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  matchType,
		MatchItems: matchItems,
	}
	return res, nil
}

// matchVirus returns the "virus" name the way it was given, without matching.
func (m matcher) matchVirus(ns nameString) (*mlib.Match, error) {
	matchItems, err := m.virusMatcher.MatchVirus(ns.Name)
	if err != nil {
		return nil, err
	}

	matchType := vlib.Virus
	if len(matchItems) == 0 {
		matchType = vlib.NoMatch
	}
	for i := range matchItems {
		matchItems[i].InputStr = ns.Name
	}
	res := &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  matchType,
		MatchItems: matchItems,
	}
	return res, nil
}
