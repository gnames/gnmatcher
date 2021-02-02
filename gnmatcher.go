package gnmatcher

import (
	"github.com/gnames/gnlib/ent/gnvers"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/ent/exact"
	"github.com/gnames/gnmatcher/ent/fuzzy"
	"github.com/gnames/gnmatcher/ent/matcher"
)

// gnmatcher implements GNmatcher interface.
type gnmatcher struct {
	matcher matcher.Matcher
}

// NewGNmatcher is a constructor for GNmatcher interface. It takes two
// interfaces ExactMatcher and FuzzyMatcher.
func NewGNmatcher(em exact.ExactMatcher, fm fuzzy.FuzzyMatcher, j int) GNmatcher {
	gnm := gnmatcher{}
	gnm.matcher = matcher.NewMatcher(em, fm, j)
	gnm.matcher.Init()
	return gnm
}

func (gnm gnmatcher) MatchNames(names []string) []mlib.Match {
	return gnm.matcher.MatchNames(names)
}

func (gnm gnmatcher) GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}
