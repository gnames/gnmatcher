// package bloom creates and serves bloom filters for canonical names, and names of viruses. The
// filters are persistent throughout the life of the program. The filters are
// used to find exact matches to the database data fast.
package bloom

import (
	"github.com/gnames/gnlib/sys"
	"github.com/gnames/gnmatcher/config"
	log "github.com/sirupsen/logrus"
)

type exactMatcher struct {
	cfg     config.Config
	filters *Filters
}

func NewExactMatcher(cfg config.Config) *exactMatcher {
	return &exactMatcher{cfg: cfg}
}

func (em *exactMatcher) Init() {
	em.prepareDir()
	log.Println("Initializing bloom filters.")
	em.getFilters()
}

func (em *exactMatcher) MatchCanonicalID(uuid string) bool {
	em.filters.Mux.Lock()
	isIn := em.filters.Canonical.Check([]byte(uuid))
	em.filters.Mux.Unlock()
	return isIn
}

func (em *exactMatcher) MatchNameStringID(uuid string) bool {
	em.filters.Mux.Lock()
	isIn := em.filters.Virus.Check([]byte(uuid))
	em.filters.Mux.Unlock()
	return isIn
}

func (em exactMatcher) prepareDir() {
	log.Println("Preparing dir for bloom filters.")
	bloomDir := em.cfg.FiltersDir()
	err := sys.MakeDir(em.cfg.FiltersDir())
	if err != nil {
		log.Fatalf("Cannot create directory %s: %s.", bloomDir, err)
	}
}
