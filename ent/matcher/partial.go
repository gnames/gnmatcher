package matcher

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/ent/fuzzy"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/stemmer"
	"github.com/gnames/gnuuid"
)

// matchPartial tries to match all patial variants of a name-string. The
// process stops as soon as a match was found.
func (m matcher) matchPartial(ns nameString, parser gnparser.GNparser) *mlib.Match {
	if ns.Partial == nil {
		return emptyResult(ns)
	}

	for _, partial := range ns.Partial.Multinomials {
		if res := m.processPartial(partial, ns, parser); res != nil {
			return res
		}
	}
	return m.processPartialGenus(ns)
}

func (m matcher) processPartialGenus(ns nameString) *mlib.Match {
	gID := gnuuid.New(ns.Partial.Genus).String()
	matchItems := m.exactStemMatches(gID, ns.Partial.Genus)

	matchItems = m.filterDataSources(matchItems)
	if len(matchItems) == 0 {
		return emptyResult(ns)
	}

	for i := range matchItems {
		matchItems[i].InputStr = ns.Partial.Genus
		matchItems[i].MatchType = vlib.PartialExact
	}
	return &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  vlib.PartialExact,
		MatchItems: matchItems,
	}
}

func (m matcher) processPartial(p multinomial, ns nameString,
	parser gnparser.GNparser) *mlib.Match {
	names := []string{p.Tail, p.Head}
	for _, name := range names {
		// TODO this is probably not efficient to use parser so many times
		nsPart, parsed := newNameString(parser, name)
		if !parsed.Parsed {
			continue
		}
		matchType := vlib.PartialFuzzy
		matches := m.exactStemMatches(nsPart.CanonicalStemID, nsPart.CanonicalStem)
		if len(matches) > 0 {
			matchItems := make([]mlib.MatchItem, 0, len(matches))
			for _, v := range matches {
				v.InputStr = nsPart.Canonical
				if v.MatchStr == v.InputStr {
					matchType = vlib.PartialExact
					v.MatchType = matchType
				} else {
					editDistance := fuzzy.EditDistance(v.InputStr, v.MatchStr)
					if editDistance == -1 {
						continue
					}
					v.EditDistance = editDistance
					v.MatchType = vlib.PartialFuzzy
				}
				matchItems = append(matchItems, v)
			}

			matchItems = m.filterDataSources(matchItems)
			if len(matchItems) == 0 {
				return nil
			}

			return &mlib.Match{
				ID:         ns.ID,
				Name:       ns.Name,
				MatchType:  matchType,
				MatchItems: matchItems,
			}
		}
	}

	// if exact partial failed, try fuzzy
	for _, name := range names {
		stem := stemmer.Stem(name).Stem
		if res := m.matchFuzzy(name, stem, ns); res != nil {

			res.MatchItems = m.filterDataSources(res.MatchItems)
			if len(res.MatchItems) == 0 {
				return nil
			}

			res.MatchType = vlib.PartialFuzzy
			return res
		}
	}

	return nil
}
