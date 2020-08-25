package matcher

import "github.com/gnames/gnmatcher/protob"

func (m Matcher) Match(ns NameString) *protob.Result {
	if m.Filters.CanonicalFull.Check([]byte(ns.CanonicalFullID)) {
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

	if m.Filters.Canonical.Check([]byte(ns.CanonicalID)) {
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
	if m.Filters.Virus.Check([]byte(ns.ID)) {
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

