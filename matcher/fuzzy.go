package matcher

import (
	"bytes"

	"github.com/dgraph-io/badger/v2"
	"github.com/dvirsky/levenshtein"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/stemskv"
	log "github.com/sirupsen/logrus"
)

type fuzzyMatcher struct {
	Trie    *levenshtein.MinTree
	KeyVal  *badger.DB
	Encoder encode.Encoder
}

// Creates new a struct compatible with FuzzyMatcher interface.
func NewFuzzyMatcher(t *levenshtein.MinTree, kv *badger.DB) fuzzyMatcher {
	return fuzzyMatcher{Trie: t, KeyVal: kv, Encoder: encode.GNgob{}}
}

func (fm fuzzyMatcher) MatchStem(stem string, maxEditDistance int) []string {
	return fm.Trie.FuzzyMatches(stem, maxEditDistance)
}

func (fm fuzzyMatcher) StemToMatchItems(stem string) []mlib.MatchItem {
	var res []mlib.MatchItem
	misGob := bytes.NewBuffer(stemskv.GetValue(fm.KeyVal, stem))
	err := fm.Encoder.Decode(misGob.Bytes(), &res)
	if err != nil {
		log.Warnf("Decode in StemToMatchItems for '%s' failed: %s", stem, err)
	}
	return res
}

// MatchFuzzy tries to get fuzzy matching of a stemmed name-string to canonical
// forms from the gnames database.
func (m Matcher) matchFuzzy(name, stem string,
	ns nameString) *mlib.Match {
	cnf := m.Config
	stemMatches := m.MatchStem(stem, cnf.MaxEditDist)
	if len(stemMatches) == 0 {
		return nilResult
	}

	res := &mlib.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  vlib.Fuzzy,
		MatchItems: make([]mlib.MatchItem, 0, len(stemMatches)*2),
	}
	for _, stemMatch := range stemMatches {
		editDistanceStem := fuzzy.EditDistance(stemMatch, stem)
		if editDistanceStem == -1 {
			continue
		}
		matchItems := m.FuzzyMatcher.StemToMatchItems(stemMatch)
		for _, matchItem := range matchItems {
			matchItem.EditDistanceStem = editDistanceStem
			// runs edit distance with checks, returns -1 if checks failed.
			matchItem.EditDistance = fuzzy.EditDistance(name, matchItem.MatchStr)
			// skip matches that failed edit distance checks.
			if matchItem.EditDistance == -1 {
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
