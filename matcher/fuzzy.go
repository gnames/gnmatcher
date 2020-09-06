package matcher

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v2"
	gn "github.com/gnames/gnames/model"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/model"
	"github.com/gnames/gnmatcher/stemskv"
)

// MatchFuzzy tries to do fuzzy matchin of a stemmed name-string to canonical
// forms from the gnames database.
func (m Matcher) MatchFuzzy(name, stem string,
	ns NameString, kv *badger.DB) *model.Match {
	cnf := m.Config
	stems := m.Trie.FuzzyMatches(stem, cnf.MaxEditDist)
	if len(stems) == 0 {
		return nilResult
	}

	res := &model.Match{
		ID:         ns.ID,
		Name:       ns.Name,
		MatchType:  gn.Fuzzy,
		MatchItems: make([]model.MatchItem, 0, len(stems)*2),
	}
	for _, v := range stems {

		editDistanceStem := fuzzy.ComputeDistance(v, stem)
		var cans []stemskv.CanonicalKV
		cansGob := bytes.NewBuffer(stemskv.GetValue(kv, v))
		dec := gob.NewDecoder(cansGob)
		dec.Decode(&cans)
		for _, v := range cans {
			res.MatchItems = append(
				res.MatchItems,
				model.MatchItem{
					ID:               v.ID,
					MatchStr:         v.Name,
					EditDistanceStem: editDistanceStem,
				})
		}
	}
	calculateEditDistance(name, res)
	return res
}

// calculateEditDistance finds the difference between the canonical form of
// a name-string and canonical forms that fuzzy-matched its stemed version.
func calculateEditDistance(name string, res *model.Match) {
	for i, v := range res.MatchItems {
		res.MatchItems[i].EditDistance = fuzzy.ComputeDistance(name, v.MatchStr)
	}
}
