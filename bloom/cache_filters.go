package bloom

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gnmatcher/sys"
	log "github.com/sirupsen/logrus"
)

// filtersFromcache unmarchals data from a file, and uses the data for the
// filter creations.
func filtersFromCache(path string) error {
	cPath := filepath.Join(path, canonicalFile)
	cfPath := filepath.Join(path, canonicalFullFile)
	vPath := filepath.Join(path, virusFile)
	sizesPath := filepath.Join(path, sizesFile)
	if sys.FileExists(cPath) && sys.FileExists(cfPath) {
		if err := getFiltersFromCache(cPath, cfPath, vPath, sizesPath); err != nil {
			return err
		}
	}
	return nil
}

func getFiltersFromCache(cPath, cfPath, vPath, sizesPath string) error {
	log.Println("Geting lookup data from a cache on disk.")
	cCfg, cfCfg, vCfg := restoreConfigs(sizesPath)
	cFilter := baseBloomfilter.New(cCfg)
	cfFilter := baseBloomfilter.New(cfCfg)
	vFilter := baseBloomfilter.New(vCfg)

	cBytes, err := ioutil.ReadFile(cPath)
	if err != nil {
		return err
	}

	cfBytes, err := ioutil.ReadFile(cfPath)
	if err != nil {
		return err
	}

	vBytes, err := ioutil.ReadFile(vPath)
	if err != nil {
		return err
	}

	if err = cFilter.UnmarshalBinary(cBytes); err != nil {
		return err
	}
	if err = cfFilter.UnmarshalBinary(cfBytes); err != nil {
		return err
	}
	if err = vFilter.UnmarshalBinary(vBytes); err != nil {
		return err
	}

	filters = &Filters{
		Canonical:     cFilter,
		CanonicalFull: cfFilter,
		Virus:         vFilter,
	}
	return nil
}

func restoreConfigs(sizeFile string) (bloomfilter.Config, bloomfilter.Config,
	bloomfilter.Config) {
	var cCfg, cfCfg, vCfg bloomfilter.Config
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
		case "CanonicalFullSize":
			cfCfg = bloomfilter.Config{
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
	return cCfg, cfCfg, vCfg
}