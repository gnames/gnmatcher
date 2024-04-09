package bloom

import (
	"encoding/csv"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gnsys"
)

// filtersFromcache unmarchals data from a file, and uses the data for the
// filter creations.
func (em *exactMatcher) filtersFromCache(path string) error {
	cPath := filepath.Join(path, canonicalStemFile)
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
	slog.Info("Geting bloom lookup data from a cache on disk")
	cCfg, err := restoreConfigs(sizesPath)
	if err != nil {
		return err
	}

	cFilter := baseBloomfilter.New(cCfg)

	cBytes, err := os.ReadFile(cPath)
	if err != nil {
		return err
	}

	if err = cFilter.UnmarshalBinary(cBytes); err != nil {
		return err
	}

	em.filters = &bloomFilters{
		canonicalStem: cFilter,
	}
	return nil
}

func restoreConfigs(sizeFile string) (bloomfilter.Config, error) {
	var cCfg bloomfilter.Config
	f, err := os.Open(sizeFile)
	if err != nil {
		slog.Error("Could not open file", "file", sizeFile, "error", err)
		return cCfg, err
	}
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		slog.Error("Could not read data from file", "file", sizeFile, "error", err)
		return cCfg, err
	}
	for _, v := range rows {
		size, err := strconv.Atoi(v[1])
		if err != nil {
			slog.Error("Could not convert data size to int",
				"data", v[1],
				"error", err,
			)
			return cCfg, err
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
	return cCfg, nil
}
