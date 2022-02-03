package bloom

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gnsys"
	log "github.com/sirupsen/logrus"
)

// filtersFromcache unmarchals data from a file, and uses the data for the
// filter creations.
func (em *exactMatcher) filtersFromCache(path string) error {
	cPath := filepath.Join(path, canonicalFile)
	sizesPath := filepath.Join(path, sizesFile)
	cPathExists, err := gnsys.FileExists(cPath)
	if err != nil {
		return err
	}
	if cPathExists {
		if err := em.getFiltersFromCache(cPath, sizesPath); err != nil {
			return err
		}
	}
	return nil
}

func (em *exactMatcher) getFiltersFromCache(cPath, sizesPath string) error {
	log.Println("Geting bloom lookup data from a cache on disk.")
	cCfg := restoreConfigs(sizesPath)
	cFilter := baseBloomfilter.New(cCfg)

	cBytes, err := os.ReadFile(cPath)
	if err != nil {
		return err
	}

	if err = cFilter.UnmarshalBinary(cBytes); err != nil {
		return err
	}

	em.filters = &bloomFilters{
		canonical: cFilter,
	}
	return nil
}

func restoreConfigs(sizeFile string) bloomfilter.Config {
	var cCfg bloomfilter.Config
	f, err := os.Open(sizeFile)
	if err != nil {
		log.Warning(fmt.Sprintf("Could not open file %s: %s.", sizeFile, err))
	}
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Warning(fmt.Sprintf("Could not read data from file %s: %s.", sizeFile, err))
	}
	for _, v := range rows {
		size, err := strconv.Atoi(v[1])
		if err != nil {
			log.Warning(fmt.Sprintf("Could not convert data size to int: %s.", err))
		}
		switch v[0] {
		case "CanonicalSize":
			cCfg = bloomfilter.Config{
				N:        uint(size),
				P:        0.00001,
				HashName: bloomfilter.HASHER_OPTIMAL,
			}
		}
	}
	return cCfg
}
