// package bloom creates and serves bloom filters for canonical names, and names of viruses. The
// filters are persistent throughout the life of the program. The filters are
// used to find exact matches to the database data fast.
package bloom

import (
	"database/sql"
	"log"

	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
)

// Names of the files to create cache of bloom filters.
const (
	canonicalFile     = "canonicals.bf"
	canonicalFullFile = "canonical_fulls.bf"
	virusFile         = "viruses.bf"
	sizesFile         = "canonical_sizes.csv"
)

// filters is a package variable and once it is created it is reused.
var filters *Filters

// Filters contain bloom filters data we use for matching.
type Filters struct {
	// Canonical is a filter for matching with canonical names.
	Canonical *baseBloomfilter.Bloomfilter
	// CanonicalSize is number of entries in 'simple' canonical filter. It is
	// used as an option during Canonical filter creation.
	CanonicalSize uint
	// CanonicalFull is a filter for matching with full canonical names.
	CanonicalFull *baseBloomfilter.Bloomfilter
	// CanonicalFullSize is number of entries in 'full' canonical filter. It is
	// used as an option during CanonicalFull filter creation.
	CanonicalFullSize uint
	// Virus is a filter for matching with viruses names.
	Virus *baseBloomfilter.Bloomfilter
	// VirusesSize is a number of entries if 'viruses' filter.
	VirusSize uint
}

// GetFilters returns bloom filters for name-string matching.
// If filters had been already created before, it just returns them.
// Otherwise it creates filters from either database, or from cached files.
// Creating filters from cache is significantly faster.
func GetFilters(path string, db *sql.DB) *Filters {
	var err error

	if filters != nil {
		return filters
	}

	err = filtersFromCache(path)
	if err != nil {
		log.Fatalf("Cannot create filters at %s from cache: %s.", path, err)
	}

	if filters != nil {
		return filters
	}

	err = filtersFromDB(path, db)
	if err != nil {
		log.Fatalf("Cannot create filters at %s from database: %s.", path, err)
	}

	return filters
}
