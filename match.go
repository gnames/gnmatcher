package gnmatcher

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v2"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/protob"
	"github.com/gnames/gnmatcher/stemskv"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser"
	"gitlab.com/gogna/gnparser/pb"
	"gitlab.com/gogna/gnparser/stemmer"
)

const MaxNamesNumber = 10_000

var (
	GNUUID    = uuid.NewV5(uuid.NamespaceDNS, "globalnames.org")
	nilResult *protob.Result
)

type NameString struct {
	ID              string
	Name            string
	Cardinality     int
	Canonical       string
	CanonicalID     string
	CanonicalFull   string
	CanonicalFullID string
	CanonicalStem   string
	Partial         *Partial
}

type Partial struct {
	Genus string
	Parts []Part
}

type Part struct {
	// Tail is genus + the last epithet.
	Tail string
	// Head is the name without the last epithet.
	Head string
}

func NewNameString(parser gnparser.GNparser,
	name string) (NameString, *pb.Parsed) {
	parsed := parser.ParseToObject(name)
	if parsed.Parsed {
		ns := NameString{
			ID:              parsed.Id,
			Name:            name,
			Cardinality:     int(parsed.Cardinality),
			Canonical:       parsed.Canonical.Simple,
			CanonicalID:     uuid.NewV5(GNUUID, parsed.Canonical.Simple).String(),
			CanonicalFull:   parsed.Canonical.Full,
			CanonicalFullID: uuid.NewV5(GNUUID, parsed.Canonical.Full).String(),
			CanonicalStem:   parsed.Canonical.Stem,
		}

		ns.NewPartial(parsed)
		// We do not fuzzy-match uninomials, however there are cases when
		// a binomial lost empty space during OCR. We increase probability to
		// match such binomials, if we stem them. It happens because we use trie
		// of stemmed canonicals.
		// For example we will be able to match 'Pardosamoestus' to 'Pardosa moesta'
		if parsed.Cardinality == 1 {
			ns.CanonicalStem = stemmer.Stem(ns.Canonical).Stem
		}
		return ns, parsed
	}

	return NameString{ID: parsed.Id, Name: name}, parsed
}

func (ns *NameString) NewPartial(parsed *pb.Parsed) {
	if parsed.Cardinality < 2 {
		return
	}
	canAry := strings.Split(ns.Canonical, " ")

	ns.Partial = &Partial{Genus: canAry[0]}
	partialNum := parsed.Cardinality - 2

	// In case of binomial we return only genus
	if partialNum < 1 {
		return
	}

	ns.Partial.Parts = make([]Part, partialNum)
	for i := range ns.Partial.Parts {
		lastLen := len(canAry) - i - 1
		tail := []string{ns.Partial.Genus, canAry[lastLen]}
		head := canAry[0:lastLen]

		ns.Partial.Parts[i] = Part{
			Tail: strings.Join(tail, " "),
			Head: strings.Join(head, " "),
		}
	}
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
		ns, parsed := NewNameString(parser, name)
		if parsed.Parsed {
			if abbreviated, matchResult := detectAbbreviated(parsed); abbreviated {
				res[i] = matchResult
				continue
			}
			matchResult = gnm.Match(ns)
		} else {
			matchResult = gnm.MatchVirus(ns)
		}
		if matchResult == nil {
			matchResult = gnm.MatchFuzzy(ns.Canonical, ns.CanonicalStem, ns, kv)
		}
		if matchResult == nil {
			matchResult = gnm.MatchPartial(ns, kv)
		}
		res[i] = matchResult
	}
	return res
}

func detectAbbreviated(parsed *pb.Parsed) (bool, *protob.Result) {
	if parsed.Quality != int32(3) {
		return false, nilResult
	}
	for _, v := range parsed.QualityWarning {
		if strings.HasPrefix(v.Message, "Abbreviated") {
			return true, &protob.Result{
				Id:        parsed.Id,
				Name:      parsed.Verbatim,
				MatchType: protob.MatchType_NONE,
			}
		}
	}
	return false, nilResult
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

func (gnm GNMatcher) MatchFuzzy(name, stem string,
	ns NameString, kv *badger.DB) *protob.Result {
	stems := gnm.Trie.FuzzyMatches(stem, gnm.MaxEditDist)
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

func (gnm GNMatcher) MatchPartial(ns NameString, kv *badger.DB) *protob.Result {
	if ns.Partial == nil {
		return emptyResult(ns)
	}

	for _, partial := range ns.Partial.Parts {
		if res := gnm.processPartial(partial, ns, kv); res != nil {
			return res
		}
	}

	return gnm.processPartialGenus(ns)
}

func (gnm GNMatcher) processPartialGenus(ns NameString) *protob.Result {
	gID := uuid.NewV5(GNUUID, ns.Partial.Genus).String()
	if gnm.Filters.Canonical.Check([]byte(gID)) {
		return &protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_PARTIAL,
			MatchData: []*protob.MatchItem{{Id: gID, MatchStr: ns.Partial.Genus}},
		}
	}
	return emptyResult(ns)
}

func (gnm GNMatcher) processPartial(p Part, ns NameString,
	kv *badger.DB) *protob.Result {
	names := []string{p.Tail, p.Head}
	for _, name := range names {
		id := uuid.NewV5(GNUUID, name).String()
		if gnm.Filters.Canonical.Check([]byte(id)) {
			return &protob.Result{
				Id:        ns.ID,
				Name:      ns.Name,
				MatchType: protob.MatchType_PARTIAL,
				MatchData: []*protob.MatchItem{{Id: id, MatchStr: ns.Partial.Genus}},
			}
		}

		stem := stemmer.Stem(name).Stem
		if res := gnm.MatchFuzzy(name, stem, ns, kv); res != nil {
			res.MatchType = protob.MatchType_PARTIAL_FUZZY
			return res
		}
	}
	return nilResult
}

func calculateEditDistance(name string, res *protob.Result) *protob.Result {
	for _, v := range res.MatchData {
		v.EditDistance = int32(fuzzy.ComputeNameDistance(name, v.MatchStr))
	}
	return res
}

func emptyResult(ns NameString) *protob.Result {
	return &protob.Result{
		Id:        ns.ID,
		Name:      ns.Name,
		MatchType: protob.MatchType_NONE,
	}
}
