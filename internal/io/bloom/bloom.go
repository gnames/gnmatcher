// package bloom creates and serves bloom filters for stemmed canonical names,
// and names of viruses. The filters are persistent throughout the life of the
// program. The filters are used to find exact matches to the database data
// fast.
package bloom

import (
	"log/slog"

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
	em := &exactMatcher{cfg: cfg}
	return em
}

func (em *exactMatcher) Init() error {
	err := em.prepareDir()
	if err != nil {
		return err
	}

	slog.Info("Initializing bloom filters")
	err = em.getFilters()
	if err != nil {
		return err
	}
	return nil
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

func (em exactMatcher) prepareDir() error {
	slog.Info("Preparing dir for bloom filters")
	bloomDir := em.cfg.FiltersDir()
	err := gnsys.MakeDir(em.cfg.FiltersDir())
	if err != nil {
		slog.Error("Cannot create directory", "path", bloomDir, "error", err)
		return err
	}
	return nil
}
