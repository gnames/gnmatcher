// package bloom creates and serves bloom filters for stemmed canonical names,
// and names of viruses. The filters are persistent throughout the life of the
// program. The filters are used to find exact matches to the database data
// fast.
package bloom

import (
	"log/slog"
	"os"

	"github.com/gnames/gnmatcher/internal/ent/exact"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/gnames/gnsys"
)

type exactMatcher struct {
	cfg     config.Config
	filters *bloomFilters
}

// New takes configuration object and returns ExactMatcher.
func New(cfg config.Config) exact.ExactMatcher {
	return &exactMatcher{cfg: cfg}
}

func (em *exactMatcher) Init() {
	em.prepareDir()
	slog.Info("Initializing bloom filters")
	em.getFilters()
}

// SetConfig updates configuration of the matcher.
func (em *exactMatcher) SetConfig(cfg config.Config) {
	em.cfg = cfg
}

func (em *exactMatcher) MatchCanonicalID(uuid string) bool {
	em.filters.mux.Lock()
	isIn := em.filters.canonicalStem.Check([]byte(uuid))
	em.filters.mux.Unlock()
	return isIn
}

func (em exactMatcher) prepareDir() {
	slog.Info("Preparing dir for bloom filters")
	bloomDir := em.cfg.FiltersDir()
	err := gnsys.MakeDir(em.cfg.FiltersDir())
	if err != nil {
		slog.Error("Cannot create directory", "path", bloomDir, "error", err)
		os.Exit(1)
	}
}
