package matcher

import "github.com/gnames/gnmatcher/protob"

func (m Matcher) Match(ns NameString) *protob.Result {
	var isIn bool
	m.Filters.Mux.Lock()
	isIn = m.Filters.CanonicalFull.Check([]byte(ns.CanonicalFullID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_CANONICAL_FULL,
			MatchData: []*protob.MatchItem{
				{
					Id:       ns.CanonicalFullID,
					MatchStr: ns.CanonicalFull,
				},
			},
		}
	}
	m.Filters.Mux.Lock()
	isIn = m.Filters.Canonical.Check([]byte(ns.CanonicalID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_CANONICAL,
			MatchData: []*protob.MatchItem{
				{
					Id:       ns.CanonicalID,
					MatchStr: ns.Canonical,
				},
			},
		}
	}
	return nilResult
}

func (m Matcher) MatchVirus(ns NameString) *protob.Result {
	var isIn bool
	m.Filters.Mux.Lock()
	isIn = m.Filters.Virus.Check([]byte(ns.ID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_VIRUS,
			MatchData: []*protob.MatchItem{
				{
					Id:       ns.ID,
					MatchStr: ns.Name,
				},
			},
		}
	}
	return emptyResult(ns)
}
