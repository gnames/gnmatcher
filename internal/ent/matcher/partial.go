package matcher

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/internal/ent/fuzzy"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/stemmer"
	"github.com/gnames/gnuuid"
)

// matchPartial tries to match all patial variants of a name-string. The
// process stops as soon as a match was found.
func (m matcher) matchPartial(
	ns nameString,
	parser gnparser.GNparser,
) (*mlib.Match, error) {
	var res *mlib.Match
	var err error
	if ns.Partial == nil {
		return emptyResult(ns), nil
	}

	for _, partial := range ns.Partial.Multinomials {
		if res, err = m.processPartial(partial, ns, parser); res != nil {
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
	res, err = m.processPartialGenus(ns)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (m matcher) processPartialGenus(ns nameString) (*mlib.Match, error) {
	gID := gnuuid.New(ns.Partial.Genus).String()
	matchItems, err := m.exactStemMatches(gID, ns.Partial.Genus)
	if err != nil {
		return nil, err
	}

	matchItems = m.filterDataSources(matchItems)

	if len(matchItems) == 0 && m.cfg.WithUninomialFuzzyMatch {
		gen := ns.Partial.Genus
		res, err := m.matchFuzzy(gen, gen, ns)
		if err != nil {
			return nil, err
		}
		if len(res.MatchItems) == 0 {
			return nil, nil
		}
		res.MatchType = vlib.PartialFuzzy
		for i := range res.MatchItems {
			res.MatchItems[i].MatchType = vlib.PartialFuzzy
		}
		return res, nil
	}

	if len(matchItems) == 0 {
		return emptyResult(ns), nil
	}

	for i := range matchItems {
		matchItems[i].InputStr = ns.Partial.Genus
		matchItems[i].MatchType = vlib.PartialExact
	}
	res := &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  vlib.PartialExact,
		MatchItems: matchItems,
	}
	return res, nil
}

func (m matcher) processPartial(p multinomial, ns nameString,
	parser gnparser.GNparser) (*mlib.Match, error) {
	names := []string{p.Tail, p.Head}
	relax := m.cfg.WithRelaxedFuzzyMatch
	for _, name := range names {
		nsPart, parsed := newNameString(parser, name)
		if !parsed.Parsed {
			continue
		}
		matchType := vlib.PartialFuzzy
		matches, err := m.exactStemMatches(nsPart.CanonicalStemID, nsPart.CanonicalStem)
		if err != nil {
			return nil, err
		}
		if len(matches) > 0 {
			matchItems := make([]mlib.MatchItem, 0, len(matches))
			for _, v := range matches {
				v.InputStr = nsPart.Canonical
				if v.MatchStr == v.InputStr {
					matchType = vlib.PartialExact
					v.MatchType = matchType
				} else {
					editDistance := fuzzy.EditDistance(v.InputStr, v.MatchStr, relax)
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
				return nil, nil
			}

			res := &mlib.Match{
				ID:         ns.ID,
				Name:       ns.Name,
				MatchType:  matchType,
				MatchItems: matchItems,
			}
			return res, nil
		}
	}

	// if exact partial failed, try fuzzy
	for _, name := range names {
		stem := stemmer.Stem(name).Stem
		if res, err := m.matchFuzzy(name, stem, ns); res != nil {
			if err != nil {
				return nil, err
			}
			res.MatchItems = m.filterDataSources(res.MatchItems)
			if len(res.MatchItems) == 0 {
				return nil, nil
			}

			for i := range res.MatchItems {
				res.MatchItems[i].MatchType = vlib.PartialFuzzy
			}

			res.MatchType = vlib.PartialFuzzy
			return res, nil
		}
	}

	return nil, nil
}
