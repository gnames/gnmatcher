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
	vPath := filepath.Join(path, virusFile)
	sizesPath := filepath.Join(path, sizesFile)
	cPathExists, err := gnsys.FileExists(cPath)
	if err != nil {
		return err
	}
	vPathExists, err := gnsys.FileExists(vPath)
	if err != nil {
		return err
	}
	if cPathExists && vPathExists {
		if err := em.getFiltersFromCache(cPath, vPath, sizesPath); err != nil {
			return err
		}
	}
	return nil
}

func (em *exactMatcher) getFiltersFromCache(cPath, vPath, sizesPath string) error {
	log.Info().Msg("Geting lookup data from a cache on disk.")
	cCfg, vCfg := restoreConfigs(sizesPath)
	cFilter := baseBloomfilter.New(cCfg)
	vFilter := baseBloomfilter.New(vCfg)

	cBytes, err := os.ReadFile(cPath)
	if err != nil {
		return err
	}

	vBytes, err := os.ReadFile(vPath)
	if err != nil {
		return err
	}

	if err = cFilter.UnmarshalBinary(cBytes); err != nil {
		return err
	}
	if err = vFilter.UnmarshalBinary(vBytes); err != nil {
		return err
	}

	em.filters = &bloomFilters{
		canonical: cFilter,
		virus:     vFilter,
	}
	return nil
}

func restoreConfigs(sizeFile string) (bloomfilter.Config, bloomfilter.Config) {
	var cCfg, vCfg bloomfilter.Config
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
		case "VirusSize":
			vCfg = bloomfilter.Config{
				N:        uint(size),
				P:        0.00001,
				HashName: bloomfilter.HASHER_OPTIMAL,
			}
		}
	}
	return cCfg, vCfg
}
