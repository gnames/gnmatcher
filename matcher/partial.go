package matcher

import (
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/gogna/gnparser/stemmer"
)

// MatchPartial tries to match all patial variants of a name-string. The
// process stops as soon as a match was found.
func (m Matcher) MatchPartial(ns NameString) *mlib.Match {
	if ns.Partial == nil {
		return emptyResult(ns)
	}

	for _, partial := range ns.Partial.Multinomials {
		if res := m.processPartial(partial, ns); res != nil {
			return res
		}
	}

	return m.processPartialGenus(ns)
}

func (m Matcher) processPartialGenus(ns NameString) *mlib.Match {
	var isIn bool
	gID := uuid.NewV5(GNUUID, ns.Partial.Genus).String()
	m.Filters.Mux.Lock()
	isIn = m.Filters.Canonical.Check([]byte(gID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &mlib.Match{
			ID:         ns.ID,
			Name:       ns.Name,
			MatchType:  vlib.PartialExact,
			MatchItems: []mlib.MatchItem{{ID: gID, MatchStr: ns.Partial.Genus}},
		}
	}
	return emptyResult(ns)
}

func (m Matcher) processPartial(p Multinomial, ns NameString) *mlib.Match {
	names := []string{p.Tail, p.Head}
	for _, name := range names {
		id := uuid.NewV5(GNUUID, name).String()
		m.Filters.Mux.Lock()
		isIn := m.Filters.Canonical.Check([]byte(id))
		m.Filters.Mux.Unlock()
		if isIn {
			return &mlib.Match{
				ID:         ns.ID,
				Name:       ns.Name,
				MatchType:  vlib.PartialExact,
				MatchItems: []mlib.MatchItem{{ID: id, MatchStr: ns.Partial.Genus}},
			}
		}

		stem := stemmer.Stem(name).Stem
		if res := m.MatchFuzzy(name, stem, ns); res != nil {
			res.MatchType = vlib.PartialFuzzy
			return res
		}
	}
	return nilResult
}
