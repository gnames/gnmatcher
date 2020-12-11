package gnmatcher

import (
	"github.com/gnames/gnlib/domain/entity/gn"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnmatcher/entity/exact"
	"github.com/gnames/gnmatcher/entity/fuzzy"
	"github.com/gnames/gnmatcher/entity/matcher"
)

// gnmatcher implements GNMatcher interface.
type gnmatcher struct {
	matcher matcher.Matcher
}

// NewGNMatcher is a constructor for GNMatcher interface. It takes two
// interfaces ExactMatcher and FuzzyMatcher.
func NewGNMatcher(em exact.ExactMatcher, fm fuzzy.FuzzyMatcher, j int) GNMatcher {
	gnm := gnmatcher{}
	gnm.matcher = matcher.NewMatcher(em, fm, j)
	gnm.matcher.Init()
	return gnm
}

func (gnm gnmatcher) MatchNames(names []string) []mlib.Match {
	return gnm.matcher.MatchNames(names)
}

func (gnm gnmatcher) GetVersion() gn.Version {
	return gn.Version{
		Version: Version,
		Build:   Build,
	}
}
