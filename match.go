package gnmatcher

import (
	"github.com/gnames/gnmatcher/protob"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser"
)

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
		return NameString{
			ID:              parsed.Id,
			Name:            name,
			Canonical:       parsed.Canonical.Simple,
			CanonicalID:     uuid.NewV5(gnm.GNUUID, parsed.Canonical.Simple).String(),
			CanonicalFull:   parsed.Canonical.Full,
			CanonicalFullID: uuid.NewV5(gnm.GNUUID, parsed.Canonical.Full).String(),
			CanonicalStem:   parsed.Canonical.Stem,
		}, parsed.Parsed
	}

	return NameString{ID: parsed.Id, Name: name}, parsed.Parsed
}

func (gnm GNMatcher) MatchNames(names []string) []*protob.Result {
	res := make([]*protob.Result, len(names))
	parser := gnparser.NewGNparser()
	log.Printf("Processing %d names", len(names))
	for i, name := range names {
		ns, parsed := gnm.NewNameString(parser, name)
		if parsed {
			match := gnm.Match(ns)
			res[i] = &match
		} else {
			res[i] = &protob.Result{Name: name}
		}
	}
	return res
}

func (gnm GNMatcher) Match(ns NameString) protob.Result {
	log.Debug(ns)
	if gnm.Filters.CanonicalFull.Check([]byte(ns.CanonicalFullID)) {
		return protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_CANONICAL_FULL,
			MatchData: []*protob.MatchItem{
				{
					CanonicalId: ns.CanonicalFullID,
					Canonical:   ns.CanonicalFull,
				},
			},
		}
	}

	if gnm.Filters.Canonical.Check([]byte(ns.CanonicalID)) {
		return protob.Result{
			Id:        ns.ID,
			Name:      ns.Name,
			MatchType: protob.MatchType_CANONICAL,
			MatchData: []*protob.MatchItem{
				{
					CanonicalId: ns.CanonicalID,
					Canonical:   ns.Canonical,
				},
			},
		}
	}

	return protob.Result{
		Id:        ns.ID,
		Name:      ns.Name,
		MatchType: protob.MatchType_NONE,
	}
}
