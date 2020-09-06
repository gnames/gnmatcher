package matcher

import (
	"github.com/dgraph-io/badger/v2"
	gn "github.com/gnames/gnames/model"
	"github.com/gnames/gnmatcher/model"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/gogna/gnparser/stemmer"
)

// MatchPartial tries to match all patial variants of a name-string. The
// process stops as soon as a match was found.
func (m Matcher) MatchPartial(ns NameString, kv *badger.DB) *model.Match {
	if ns.Partial == nil {
		return emptyResult(ns)
	}

	for _, partial := range ns.Partial.Multinomials {
		if res := m.processPartial(partial, ns, kv); res != nil {
			return res
		}
	}

	return m.processPartialGenus(ns)
}

func (m Matcher) processPartialGenus(ns NameString) *model.Match {
	var isIn bool
	gID := uuid.NewV5(GNUUID, ns.Partial.Genus).String()
	m.Filters.Mux.Lock()
	isIn = m.Filters.Canonical.Check([]byte(gID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &model.Match{
			ID:         ns.ID,
			Name:       ns.Name,
			MatchType:  gn.PartialExact,
			MatchItems: []model.MatchItem{{ID: gID, MatchStr: ns.Partial.Genus}},
		}
	}
	return emptyResult(ns)
}

func (m Matcher) processPartial(p Multinomial, ns NameString,
	kv *badger.DB) *model.Match {
	names := []string{p.Tail, p.Head}
	for _, name := range names {
		id := uuid.NewV5(GNUUID, name).String()
		m.Filters.Mux.Lock()
		isIn := m.Filters.Canonical.Check([]byte(id))
		m.Filters.Mux.Unlock()
		if isIn {
			return &model.Match{
				ID:         ns.ID,
				Name:       ns.Name,
				MatchType:  gn.PartialExact,
				MatchItems: []model.MatchItem{{ID: id, MatchStr: ns.Partial.Genus}},
			}
		}

		stem := stemmer.Stem(name).Stem
		if res := m.MatchFuzzy(name, stem, ns, kv); res != nil {
			res.MatchType = gn.PartialFuzzy
			return res
		}
	}
	return nilResult
}
