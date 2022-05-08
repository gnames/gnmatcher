package gnmatcher

import (
	"github.com/gnames/gnlib/ent/gnvers"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/ent/exact"
	"github.com/gnames/gnmatcher/ent/fuzzy"
	"github.com/gnames/gnmatcher/ent/matcher"
	"github.com/gnames/gnmatcher/ent/virus"
)

// gnmatcher implements GNmatcher interface.
type gnmatcher struct {
	cfg     config.Config
	matcher matcher.Matcher
}

// New is a constructor for GNmatcher interface. It takes two
// interfaces ExactMatcher and FuzzyMatcher.
func New(
	em exact.ExactMatcher,
	fm fuzzy.FuzzyMatcher,
	vm virus.VirusMatcher,
	cfg config.Config,
) GNmatcher {
	gnm := gnmatcher{cfg: cfg}
	gnm.matcher = matcher.NewMatcher(em, fm, vm, cfg)
	gnm.matcher.Init()
	return gnm
}

func (gnm gnmatcher) MatchNames(names []string, opts ...config.Option) mlib.Output {
	return gnm.matcher.MatchNames(names, opts...)
}

func (gnm gnmatcher) GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}

func (gnm gnmatcher) GetConfig() config.Config {
	return gnm.cfg
}
