package matcher

import (
	"bytes"
	"fmt"

	"github.com/dgraph-io/badger/v2"
	"github.com/dvirsky/levenshtein"
	gn "github.com/gnames/gnames/domain/entity"
	"github.com/gnames/gnames/lib/encode"
	"github.com/gnames/gnmatcher/domain/entity"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/stemskv"
	log "github.com/sirupsen/logrus"
)

type FuzzyMatcherTrie struct {
	Trie    *levenshtein.MinTree
	KeyVal  *badger.DB
	Encoder encode.Encoder
}

func NewFuzzyMatcherTrie(t *levenshtein.MinTree, kv *badger.DB) FuzzyMatcherTrie {
	return FuzzyMatcherTrie{Trie: t, KeyVal: kv, Encoder: encode.GNgob{}}
}

func (fm FuzzyMatcherTrie) MatchStem(stem string, maxEditDistance int) []string {
	return fm.Trie.FuzzyMatches(stem, maxEditDistance)
}

func (fm FuzzyMatcherTrie) StemToMatchItems(stem string) []entity.MatchItem {
	var res []entity.MatchItem
	misGob := bytes.NewBuffer(stemskv.GetValue(fm.KeyVal, stem))
	err := fm.Encoder.Decode(misGob.Bytes(), &res)
	if err != nil {
		log.Warnf("Decode in StemToMatchItems for '%s' failed: %s", stem, err)
	}
	return res
}

// MatchFuzzy tries to get fuzzy matching of a stemmed name-string to canonical
// forms from the gnames database.
func (m Matcher) MatchFuzzy(name, stem string,
	ns NameString) *entity.Match {
	cnf := m.Config
	stemMatches := m.MatchStem(stem, cnf.MaxEditDist)
	if len(stemMatches) == 0 {
		return nilResult
	}

	res := &entity.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  gn.Fuzzy,
		MatchItems: make([]entity.MatchItem, 0, len(stemMatches)*2),
	}
	for _, stemMatch := range stemMatches {
		editDistanceStem := fuzzy.ComputeDistance(stemMatch, stem)
		matchItems := m.FuzzyMatcher.StemToMatchItems(stemMatch)
		for _, matchItem := range matchItems {
			matchItem.EditDistanceStem = editDistanceStem
			matchItem.EditDistance = fuzzy.ComputeDistance(name, matchItem.MatchStr)
			// skip matches with too large edit distance
			if matchItem.EditDistance > 2 {
				continue
			}
			res.MatchItems = append(res.MatchItems, matchItem)
		}
	}
	if len(res.MatchItems) == 0 {
		return nil
	}
	return res
}
