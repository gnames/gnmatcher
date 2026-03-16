package gnmatcher

import (
	"github.com/gnames/gnlib/ent/gnvers"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/internal/ent/matcher"
	"github.com/gnames/gnmatcher/internal/io/bloom"
	"github.com/gnames/gnmatcher/internal/io/trie"
	"github.com/gnames/gnmatcher/internal/io/virusio"
	"github.com/gnames/gnmatcher/pkg/config"
)

// gnmatcher implements GNmatcher interface.
type gnmatcher struct {
	cfg     config.Config
	matcher matcher.Matcher
}

// New creates a GNmatcher from config. It wires internal components but
// performs no I/O. Call Init() to load caches and connect to the database.
func New(cfg config.Config) GNmatcher {
	em := bloom.New(cfg)
	fm := trie.New(cfg)
	vm := virusio.New(cfg)
	return gnmatcher{
		cfg:     cfg,
		matcher: matcher.NewMatcher(em, fm, vm, cfg),
	}
}

func (gnm gnmatcher) Init() error {
	return gnm.matcher.Init()
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
