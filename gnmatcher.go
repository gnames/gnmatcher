package gnmatcher

import (
	"fmt"

	"github.com/gnames/gnmatcher/matcher"
	"github.com/gnames/gnmatcher/protob"
	"github.com/gnames/gnmatcher/stemskv"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser"
)

// MaxMaxNamesNumber is the upper limit of the number of name-strings the
// MatchNames function can process. If the number is higher, the list of
// name-strings will be truncated.
const MaxNamesNumber = 10_000

// GNMatcher contains high level methods for scientific name matching.
type GNMatcher struct {
	Matcher matcher.Matcher
}

// NewGNMatcher is a constructor for GNMatcher instance
func NewGNMatcher(m matcher.Matcher) GNMatcher {
	return GNMatcher{Matcher: m}
}

// MatchNames takes a list of name-strings and matches them against known
// by names aggregated in gnames database.
func (gnm GNMatcher) MatchNames(names []string) []*protob.Result {
	m := gnm.Matcher
	cnf := m.Config
	kv := stemskv.ConnectKeyVal(cnf.StemsDir())
	defer kv.Close()

	res := make([]*protob.Result, len(names))
	var matchResult *protob.Result
	parser := gnparser.NewGNparser()

	names = truncateNamesToMaxNumber(names)

	log.Printf("Processing %d names.", len(names))
	for i, name := range names {
		ns, parsed := matcher.NewNameString(parser, name)
		if parsed.Parsed {
			if abbrResult := matcher.DetectAbbreviated(parsed); abbrResult != nil {
				res[i] = abbrResult
				continue
			}
			matchResult = m.Match(ns)
		} else {
			matchResult = m.MatchVirus(ns)
		}
		if matchResult == nil {
			matchResult = m.MatchFuzzy(ns.Canonical, ns.CanonicalStem, ns, kv)
		}
		if matchResult == nil {
			matchResult = m.MatchPartial(ns, kv)
		}
		res[i] = matchResult
	}
	return res
}

func truncateNamesToMaxNumber(names []string) []string {
	if len(names) > MaxNamesNumber {
		log.Warn(fmt.Sprintf("Too many names, truncating list to %d entries.",
			MaxNamesNumber))
		names = names[0:MaxNamesNumber]
	}
	return names
}
