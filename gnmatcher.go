package gnmatcher

import (
	"github.com/gnames/gnlib/domain/entity/gn"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnmatcher/entity/exact"
	"github.com/gnames/gnmatcher/entity/fuzzy"
	"github.com/gnames/gnmatcher/entity/matcher"
)

// MaxMaxNamesNumber is the upper limit of the number of name-strings the
// MatchNames function can process. If the number is higher, the list of
// name-strings will be truncated.
const MaxNamesNumber = 10_000

// gnmatcher implements GNMatcher interface.
type gnmatcher struct {
	matcher matcher.Matcher
}

// NewGNMatcher is a constructor for GNMatcher interface
func NewGNMatcher(em exact.ExactMatcher, fm fuzzy.FuzzyMatcher) gnmatcher {
	gnm := gnmatcher{}
	gnm.matcher = matcher.NewMatcher(em, fm)
	gnm.matcher.Init()
	return gnm
}

// MatchNames takes a list of name-strings and matches them against
// known names aggregated in gnames database.
func (gnm gnmatcher) MatchNames(names []string) []*mlib.Match {
	return gnm.matcher.MatchNames(names)
}

func (gnm gnmatcher) GetVersion() gn.Version {
	return gn.Version{
		Version: Version,
		Build:   Build,
	}
}
