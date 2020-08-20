package gnmatcher

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/badger/v2"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/protob"
	"github.com/gnames/gnmatcher/stemskv"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser"
	"gitlab.com/gogna/gnparser/stemmer"
)

const MaxNamesNumber = 10_000

var nilResult *protob.Result

type NameString struct {
	ID              string
	Name            string
	Canonical       string
	CanonicalID     string
	CanonicalFull   string
	CanonicalFullID string
	CanonicalStem   string
}

func (gnm GNMatcher) NewNameString(parser gnparser.GNparser, name string) (NameString, bool) {
	parsed := parser.ParseToObject(name)
	if parsed.Parsed {
		ns := NameString{
			ID:              parsed.Id,
			Name:            name,
			Canonical:       parsed.Canonical.Simple,
			CanonicalID:     uuid.NewV5(gnm.GNUUID, parsed.Canonical.Simple).String(),
			CanonicalFull:   parsed.Canonical.Full,
			CanonicalFullID: uuid.NewV5(gnm.GNUUID, parsed.Canonical.Full).String(),
			CanonicalStem:   parsed.Canonical.Stem,
		}

		// We do not fuzzy matching uninomials, however there might be cases when
		// a binomial lost empty space during OCR. We increase probability to
		// match such binomials, if we stem them. It happens because we use trie
		// of stemmed canonicals.
		// For example we will be able to match 'Pardosamoestus' to 'Pardosa moesta'
		if parsed.Cardinality == 1 {
			ns.CanonicalStem = stemmer.Stem(ns.Canonical).Stem
		}
		return ns, parsed.Parsed
	}

	return NameString{ID: parsed.Id, Name: name}, parsed.Parsed
}

func (gnm GNMatcher) MatchNames(names []string) []*protob.Result {
	kv := stemskv.ConnectKeyVal(gnm.StemsDir())
	defer kv.Close()
	res := make([]*protob.Result, len(names))
	var matchResult *protob.Result
	parser := gnparser.NewGNparser()
	if len(names) <= MaxNamesNumber {
		log.Printf("Processing %d names.", len(names))
	} else {
		log.Warn(fmt.Sprintf("Too many names, truncating list to %d entries.",
			MaxNamesNumber))
		names = names[0:MaxNamesNumber]
	}

	for i, name := range names {
		ns, parsed := gnm.NewNameString(parser, name)
		if parsed {
			matchResult = gnm.Match(ns)
		} else {
			matchResult = gnm.MatchVirus(ns)
		}
		if matchResult == nil {
			matchResult = gnm.MatchFuzzy(ns, kv)
		}
		res[i] = matchResult
	}
	return res
}

func (gnm GNMatcher) Match(ns NameString) *protob.Result {
	if gnm.Filters.CanonicalFull.Check([]byte(ns.CanonicalFullID)) {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_CANONICAL_FULL,
			MatchData: []*protob.MatchItem{
				{
					Id:       ns.CanonicalFullID,
					MatchStr: ns.CanonicalFull,
				},
			},
		}
	}

	if gnm.Filters.Canonical.Check([]byte(ns.CanonicalID)) {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_CANONICAL,
			MatchData: []*protob.MatchItem{
				{
					Id:       ns.CanonicalID,
					MatchStr: ns.Canonical,
				},
			},
		}
	}

	return nilResult
}

func (gnm GNMatcher) MatchVirus(ns NameString) *protob.Result {
	if gnm.Filters.Virus.Check([]byte(ns.ID)) {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_VIRUS,
			MatchData: []*protob.MatchItem{
				{
					Id:       ns.ID,
					MatchStr: ns.Name,
				},
			},
		}
	}
	return nilResult
}

func (gnm GNMatcher) MatchFuzzy(ns NameString, kv *badger.DB) *protob.Result {
	stems := gnm.Trie.FuzzyMatches(ns.CanonicalStem, 1)
	if len(stems) == 0 {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_NONE,
		}
	}

	res := &protob.Result{
		Id:        ns.ID,
		Name:      ns.Name,
		MatchType: protob.MatchType_FUZZY,
		MatchData: make([]*protob.MatchItem, 0, len(stems)*2),
	}
	for _, v := range stems {
		editDistanceStem := int32(fuzzy.ComputeDistance(v, ns.CanonicalStem))
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
	res = calculateEditDistance(res)
	log.Debug(fmt.Sprintf("%+v", res))
	return res
}

func calculateEditDistance(res *protob.Result) *protob.Result {
	for _, v := range res.MatchData {
		v.EditDistance = int32(fuzzy.ComputeNameDistance(res.Name, v.MatchStr))
	}
	return res
}
