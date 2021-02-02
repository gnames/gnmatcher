// package bloom creates and serves bloom filters for canonical names, and names of viruses. The
// filters are persistent throughout the life of the program. The filters are
// used to find exact matches to the database data fast.
package bloom

import (
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/ent/exact"
	"github.com/gnames/gnsys"
	log "github.com/sirupsen/logrus"
)

type exactMatcher struct {
	cfg     config.Config
	filters *bloomFilters
}

// NewENewExactMatcher takes configuration object and returns ExactMatcher.
func NewExactMatcher(cfg config.Config) exact.ExactMatcher {
	return &exactMatcher{cfg: cfg}
}

func (em *exactMatcher) Init() {
	em.prepareDir()
	log.Println("Initializing bloom filters.")
	em.getFilters()
}

func (em *exactMatcher) MatchCanonicalID(uuid string) bool {
	em.filters.mux.Lock()
	isIn := em.filters.canonical.Check([]byte(uuid))
	em.filters.mux.Unlock()
	return isIn
}

func (em *exactMatcher) MatchNameStringID(uuid string) bool {
	em.filters.mux.Lock()
	isIn := em.filters.virus.Check([]byte(uuid))
	em.filters.mux.Unlock()
	return isIn
}

func (em exactMatcher) prepareDir() {
	log.Println("Preparing dir for bloom filters.")
	bloomDir := em.cfg.FiltersDir()
	err := gnsys.MakeDir(em.cfg.FiltersDir())
	if err != nil {
		log.Fatalf("Cannot create directory %s: %s.", bloomDir, err)
	}
}
