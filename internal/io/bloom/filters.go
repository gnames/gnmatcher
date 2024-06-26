// package bloom creates and serves bloom filters for stemmed canonical names,
// and names of viruses. The filters are persistent throughout the life of the
// program. The filters are used to find exact matches to the database data
// fast.
package bloom

import (
	"log/slog"
	"sync"

	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
)

// Names of the files to create cache of bloom filters.
const (
	canonicalStemFile = "canonical_stems.bf"
	virusFile         = "viruses.bf"
	sizesFile         = "canonical_sizes.csv"
)

// bloomFilters contain bloom filters data we use for matching.
type bloomFilters struct {
	// canonicalStem is a filter for matching with canonicalStem names.
	canonicalStem *baseBloomfilter.Bloomfilter

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
func (em *exactMatcher) getFilters() error {
	path := em.cfg.FiltersDir()
	var err error

	if em.filters != nil {
		return nil
	}

	err = em.filtersFromCache(path)
	if err != nil {
		slog.Error("Cannot create filters from cache", "path", path, "error", err)
		return err
	}

	if em.filters != nil {
		return nil
	}

	err = em.filtersFromDB(path)
	if err != nil {
		slog.Error(
			"Cannot create filters from database",
			"path", path,
			"error", err,
		)
		return err
	}
	return nil
}
