// package bloom creates and serves bloom filters for canonical names, and names of viruses. The
// filters are persistent throughout the life of the program. The filters are
// used to find exact matches to the database data fast.
package bloom

import (
	"sync"

	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/rs/zerolog/log"
)

// Names of the files to create cache of bloom filters.
const (
	canonicalFile = "canonicals.bf"
	virusFile     = "viruses.bf"
	sizesFile     = "canonical_sizes.csv"
)

// bloomFilters contain bloom filters data we use for matching.
type bloomFilters struct {
	// canonical is a filter for matching with canonical names.
	canonical *baseBloomfilter.Bloomfilter

	// canonicalSize is number of entries in 'simple' canonical filter. It is
	// used as an option during Canonical filter creation.
	canonicalSize uint

	// mux is a mutex for thread-safe operations
	mux sync.Mutex
}

// getFilters returns bloom filters for name-string matching.
// If filters had been already created before, it just returns them.
// Otherwise it creates filters from either database, or from cached files.
// Creating filters from cache is significantly faster.
func (em *exactMatcher) getFilters() {
	path := em.cfg.FiltersDir()
	var err error

	if em.filters != nil {
		return
	}

	err = em.filtersFromCache(path)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create filters at %s from cache", path)
	}

	if em.filters != nil {
		return
	}

	err = em.filtersFromDB(path)
	if err != nil {
		log.Fatal().Err(err).
			Msgf("Cannot create filters at %s from database", path)
	}
}
