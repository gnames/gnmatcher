package matcher

import (
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/gnuuid"
	"gitlab.com/gogna/gnparser/stemmer"
)

// matchPartial tries to match all patial variants of a name-string. The
// process stops as soon as a match was found.
func (m matcher) matchPartial(ns nameString) *mlib.Match {
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

func (m matcher) processPartialGenus(ns nameString) *mlib.Match {
	gID := gnuuid.New(ns.Partial.Genus).String()
	isIn := m.exactMatcher.MatchCanonicalID(gID)
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

func (m matcher) processPartial(p multinomial, ns nameString) *mlib.Match {
	names := []string{p.Tail, p.Head}
	for _, name := range names {
		id := gnuuid.New(name).String()
		isIn := m.exactMatcher.MatchCanonicalID(id)
		if isIn {
			res := &mlib.Match{
				ID:         ns.ID,
				Name:       ns.Name,
				MatchType:  vlib.PartialExact,
				MatchItems: []mlib.MatchItem{{ID: id, MatchStr: ns.Partial.Genus}},
			}
			return res
		}
	}

	// if exact partial failed, try fuzzy
	for _, name := range names {
		stem := stemmer.Stem(name).Stem
		if res := m.matchFuzzy(name, stem, ns); res != nil {
			res.MatchType = vlib.PartialFuzzy
			return res
		}
	}

	return nilResult
}
