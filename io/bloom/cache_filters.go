package bloom

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Geting bloom lookup data from a cache on disk.")
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
		log.Warn().Err(err).Msgf("Could not open file %s", sizeFile)
	}
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Warn().Err(err).Msgf("Could not read data from file %s", sizeFile)
	}
	for _, v := range rows {
		size, err := strconv.Atoi(v[1])
		if err != nil {
			log.Warn().Err(err).Msg("Could not convert data size to int")
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
