// package bloom creates and serves bloom filters for stemmed canonical names,
// and names of viruses. The filters are persistent throughout the life of the
// program. The filters are used to find exact matches to the database data
// fast.
package bloom

import (
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/ent/exact"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Initializing bloom filters")
	em.getFilters()
}

func (em *exactMatcher) MatchCanonicalID(uuid string) bool {
	em.filters.mux.Lock()
	isIn := em.filters.canonicalStem.Check([]byte(uuid))
	em.filters.mux.Unlock()
	return isIn
}

func (em exactMatcher) prepareDir() {
	log.Info().Msg("Preparing dir for bloom filters")
	bloomDir := em.cfg.FiltersDir()
	err := gnsys.MakeDir(em.cfg.FiltersDir())
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create directory %s", bloomDir)
	}
}
