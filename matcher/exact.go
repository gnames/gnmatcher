package matcher

import (
	gn "github.com/gnames/gnames/domain/entity"
	"github.com/gnames/gnmatcher/domain/entity"
)

// Match tries to match a canonical form of a name-string exactly to canonical
// from from gnames database.
func (m Matcher) Match(ns NameString) *entity.Match {
	var isIn bool
	m.Filters.Mux.Lock()
	isIn = m.Filters.Canonical.Check([]byte(ns.CanonicalID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &entity.Match{
			ID:        ns.ID,
			Name:      ns.Name,
			MatchType: gn.Exact,
			MatchItems: []entity.MatchItem{
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
func (m Matcher) MatchVirus(ns NameString) *entity.Match {
	var isIn bool
	m.Filters.Mux.Lock()
	isIn = m.Filters.Virus.Check([]byte(ns.ID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &entity.Match{
			ID:         ns.ID,
			Name:       ns.Name,
			VirusMatch: true,
			MatchType:  gn.Exact,
			MatchItems: []entity.MatchItem{
				{
					ID:       ns.ID,
					MatchStr: ns.Name,
				},
			},
		}
	}
	return emptyResult(ns)
}
