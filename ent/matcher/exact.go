package matcher

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
)

// match tries to match a canonical form of a name-string exactly to canonical
// from from gnames database.
func (m matcher) match(ns nameString) *mlib.Match {
	isIn := m.isExactMatch(ns.CanonicalID, ns.CanonicalStem)
	if isIn {
		return &mlib.Match{
			ID:        ns.ID,
			Name:      ns.Name,
			MatchType: vlib.Exact,
			MatchItems: []mlib.MatchItem{
				{
					ID:       ns.CanonicalID,
					MatchStr: ns.Canonical,
				},
			},
		}
	}
	return nilResult
}

// matchVirus tries to match a name-string exactly to a virus name from the
// gnames database.
func (m matcher) matchVirus(ns nameString) *mlib.Match {
	isIn := m.exactMatcher.MatchNameStringID(ns.ID)
	if isIn {
		return &mlib.Match{
			ID:         ns.ID,
			Name:       ns.Name,
			VirusMatch: true,
			MatchType:  vlib.Exact,
			MatchItems: []mlib.MatchItem{
				{
					ID:       ns.ID,
					MatchStr: ns.Name,
				},
			},
		}
	}
	return emptyResult(ns)
}
