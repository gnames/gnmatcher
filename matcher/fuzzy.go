package matcher

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v2"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/protob"
	"github.com/gnames/gnmatcher/stemskv"
)

// MatchFuzzy tries to do fuzzy matchin of a stemmed name-string to canonical
// forms from the gnames database.
func (m Matcher) MatchFuzzy(name, stem string,
	ns NameString, kv *badger.DB) *protob.Result {
	cnf := m.Config
	stems := m.Trie.FuzzyMatches(stem, cnf.MaxEditDist)
	if len(stems) == 0 {
		return nilResult
	}

	res := &protob.Result{
		Id:        ns.ID,
		Name:      ns.Name,
		MatchType: protob.MatchType_FUZZY,
		MatchData: make([]*protob.MatchItem, 0, len(stems)*2),
	}
	for _, v := range stems {

		editDistanceStem := int32(fuzzy.ComputeDistance(v, stem))
		var cans []stemskv.CanonicalKV
		cansGob := bytes.NewBuffer(stemskv.GetValue(kv, v))
		dec := gob.NewDecoder(cansGob)
		dec.Decode(&cans)
		for _, v := range cans {
			res.MatchData = append(
				res.MatchData,
				&protob.MatchItem{
					Id:               v.ID,
					MatchStr:         v.Name,
					EditDistanceStem: editDistanceStem,
				})
		}
	}
	res = calculateEditDistance(name, res)
	return res
}

// calculateEditDistance finds the difference between the canonical form of
// a name-string and canonical forms that fuzzy-matched its stemed version.
func calculateEditDistance(name string, res *protob.Result) *protob.Result {
	for _, v := range res.MatchData {
		v.EditDistance = int32(fuzzy.ComputeDistance(name, v.MatchStr))
	}
	return res
}
