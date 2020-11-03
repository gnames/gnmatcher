package matcher

import (
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
)

// Match tries to match a canonical form of a name-string exactly to canonical
// from from gnames database.
func (m Matcher) Match(ns NameString) *mlib.Match {
	var isIn bool
	m.Filters.Mux.Lock()
	isIn = m.Filters.Canonical.Check([]byte(ns.CanonicalID))
	m.Filters.Mux.Unlock()
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

// MatchVirus tries to match a name-string exactly to a virus name from the
// gnames database.
func (m Matcher) MatchVirus(ns NameString) *mlib.Match {
	var isIn bool
	m.Filters.Mux.Lock()
	isIn = m.Filters.Virus.Check([]byte(ns.ID))
	m.Filters.Mux.Unlock()
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
