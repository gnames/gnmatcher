package matcher

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/gnames/gnmatcher/protob"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/gogna/gnparser/stemmer"
)

func (m Matcher) MatchPartial(ns NameString, kv *badger.DB) *protob.Result {
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

func (m Matcher) processPartialGenus(ns NameString) *protob.Result {
	var isIn bool
	gID := uuid.NewV5(GNUUID, ns.Partial.Genus).String()
	m.Filters.Mux.Lock()
	isIn = m.Filters.Canonical.Check([]byte(gID))
	m.Filters.Mux.Unlock()
	if isIn {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_PARTIAL,
			MatchData: []*protob.MatchItem{{Id: gID, MatchStr: ns.Partial.Genus}},
		}
	}
	return emptyResult(ns)
}

func (m Matcher) processPartial(p Multinomial, ns NameString,
	kv *badger.DB) *protob.Result {
	names := []string{p.Tail, p.Head}
	for _, name := range names {
		id := uuid.NewV5(GNUUID, name).String()
		m.Filters.Mux.Lock()
		isIn := m.Filters.Canonical.Check([]byte(id))
		m.Filters.Mux.Unlock()
		if isIn {
			return &protob.Result{
				Id:        ns.ID,
				Name:      ns.Name,
				MatchType: protob.MatchType_PARTIAL,
				MatchData: []*protob.MatchItem{{Id: id, MatchStr: ns.Partial.Genus}},
			}
		}

		stem := stemmer.Stem(name).Stem
		if res := m.MatchFuzzy(name, stem, ns, kv); res != nil {
			res.MatchType = protob.MatchType_PARTIAL_FUZZY
			return res
		}
	}
	return nilResult
}
