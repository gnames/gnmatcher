// package bloom creates and serves bloom filters for canonical names. The
// filters are persistent throughout the life of the program. The filters are
// used to find exact matches to the database data fast.
package bloom

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/sys"
	log "github.com/sirupsen/logrus"
)

const (
	canonicalFile     = "canonicals.bf"
	canonicalFullFile = "canonical_fulls.bf"
	canonicalSizeFile = "canonical_sizes.csv"
)

var filters *Filters

// Filters contain bloom filters data we use for matching.
type Filters struct {
	// CanonicalSize is number of entries in 'simple' canonical filter. It is
	// used as an option during Canonical filter creation.
	CanonicalSize uint
	// CanonicalFullSize is number of entries in 'full' canonical filter. It is
	// used as an option during CanonicalFull filter creation.
	CanonicalFullSize uint
	// Canonical is a filter for matching with canonical names.
	Canonical *baseBloomfilter.Bloomfilter
	// CanonicalFull is a filter for matching with full canonical names.
	CanonicalFull *baseBloomfilter.Bloomfilter
}

// GetFilters creates filters from either database, or from cached files.
// Creating filters from cache is significantly faster.
func GetFilters(path string, d dbase.Dbase) (*Filters, error) {
	var err error
	if filters != nil {
		return filters, nil
	}
	if err = filtersFromCache(path); filters != nil {
		return filters, err
	}
	filters = &Filters{}
	return filters, createFilters(path, d)
}

// filtersFromcache unmarchals data from a file, and uses the data for the
// filter creations.
func filtersFromCache(path string) error {
	cPath := filepath.Join(path, canonicalFile)
	cfPath := filepath.Join(path, canonicalFullFile)
	sizesPath := filepath.Join(path, canonicalSizeFile)
	if sys.FileExists(cPath) && sys.FileExists(cfPath) {
		if err := getFiltersFromCache(cPath, cfPath, sizesPath); err != nil {
			return err
		}
	}
	return nil
}

func getFiltersFromCache(cPath, cfPath, sizesPath string) error {
	log.Println("Geting lookup data from a cache on disk")
	cCfg, cfCfg := restoreConfigs(sizesPath)
	cFilter := baseBloomfilter.New(cCfg)
	cfFilter := baseBloomfilter.New(cfCfg)
	cBytes, err := ioutil.ReadFile(cPath)

	if err != nil {
		return err
	}
	cfBytes, err := ioutil.ReadFile(cfPath)
	if err != nil {
		return err
	}

	if err = cFilter.UnmarshalBinary(cBytes); err != nil {
		return err
	}
	if err = cfFilter.UnmarshalBinary(cfBytes); err != nil {
		return err
	}

	filters = &Filters{
		Canonical:     cFilter,
		CanonicalFull: cfFilter,
	}
	log.Println("Success")
	return nil
}

func restoreConfigs(path string) (bloomfilter.Config, bloomfilter.Config) {
	var cCfg, cfCfg bloomfilter.Config
	f, err := os.Open(path)
	if err != nil {
		log.Warning(fmt.Sprintf("Could not open file %s: %s \n", path, err))
	}
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Warning(fmt.Sprintf("Could not read data from file %s: %s \n", path, err))
	}
	for _, v := range rows {
		size, err := strconv.Atoi(v[1])
		if err != nil {
			log.Warning(fmt.Sprintf("Could not convert data size to int: %s \n", err))
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
		}
	}
	return cCfg, cfCfg
}

func createFilters(path string, d dbase.Dbase) error {
	log.Println("Importing lookup data from remote database. " +
		"It will take a while.")
	db := d.NewDB()
	log.Println("Importing lookup data for simple canonicals.")
	cFilter, err := createCanonicalFilter(db)
	if err != nil {
		return err
	}
	log.Println("Lookup data part 1 is imported")

	log.Println("Importing lookup data for full canonicals.")
	cfFilter, err := createCanonicalFullFilter(db)
	if err != nil {
		return err
	}
	log.Println("Lookup data part 2 is imported")
	filters.Canonical = cFilter
	filters.CanonicalFull = cfFilter
	saveFilters(path, filters)
	return nil
}

func saveFilters(path string, filters *Filters) {
	cPath := filepath.Join(path, canonicalFile)
	cfPath := filepath.Join(path, canonicalFullFile)
	sizesPath := filepath.Join(path, canonicalSizeFile)
	cFile, err := os.Create(cPath)
	if err != nil {
		log.Warning(fmt.Sprintf("Could not create file %s: %s \n", cPath, err))
	}
	cfFile, err := os.Create(cfPath)
	if err != nil {
		log.Warning(fmt.Sprintf("Could not create file %s: %s\n", cfPath, err))
	}
	sizesFile, err := os.Create(sizesPath)
	if err != nil {
		log.Warning(fmt.Sprintf("Could not create file %s: %s\n", sizesPath, err))
	}
	cBin, err := filters.Canonical.MarshalBinary()
	if err != nil {
		log.Warning(fmt.Sprintf("Could serialize lookup cache data1: %s\n", err))
	}
	cfBin, err := filters.CanonicalFull.MarshalBinary()
	if err != nil {
		log.Warning(fmt.Sprintf("Could serialize lookup cache data2 filter: %s\n", err))
	}
	_, err = cFile.Write(cBin)
	if err != nil {
		log.Warning(fmt.Sprintf("Could not save lookup cache data to disk: %s\n", err))
	}
	_, err = cfFile.Write(cfBin)
	if err != nil {
		log.Warning(fmt.Sprintf("Could not save lookup cache data2 to disk: %s\n", err))
	}
	sizes := fmt.Sprintf("CanonicalSize,%d\nCanonicalFullSize,%d\n",
		filters.CanonicalSize, filters.CanonicalFullSize)
	if _, err := sizesFile.WriteString(sizes); err != nil {
		log.Warning(fmt.Sprintf("Could not save sizes of data1 and data2 on disk: %s\n", err))
	}
	if err != nil {
		log.Warning("Failed to safe lookup data to disk")
	} else {
		log.Warning("Lookup cache data are saved to disk")
	}
}

func createCanonicalFilter(db *sql.DB) (*baseBloomfilter.Bloomfilter, error) {
	var bf *baseBloomfilter.Bloomfilter
	var err error
	table := "canonicals"
	filters.CanonicalSize, err = getFilterSize(db, table)

	if err != nil {
		return bf, err
	}
	return getFilter(db, table, filters.CanonicalSize)
}

func createCanonicalFullFilter(db *sql.DB) (*baseBloomfilter.Bloomfilter, error) {
	var err error
	var bf *baseBloomfilter.Bloomfilter
	table := "canonical_fulls"
	filters.CanonicalFullSize, err = getFilterSize(db, table)
	if err != nil {
		return bf, err
	}
	return getFilter(db, table, filters.CanonicalFullSize)
}

func getFilterSize(db *sql.DB, table string) (uint, error) {
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	var num uint
	row := db.QueryRow(q)
	if err := row.Scan(&num); err != nil {
		return 0, err
	}
	return num, nil
}

func getFilter(db *sql.DB, table string,
	filterSize uint) (*baseBloomfilter.Bloomfilter, error) {
	var uuid string
	cfg := bloomfilter.Config{
		N:        filterSize,
		P:        0.00001,
		HashName: bloomfilter.HASHER_OPTIMAL,
	}
	bf := baseBloomfilter.New(cfg)
	q := fmt.Sprintf("SELECT id FROM %s", table)
	rows, err := db.Query(q)
	if err != nil {
		return bf, err
	}
	for rows.Next() {
		if err := rows.Scan(&uuid); err != nil {
			return bf, err
		}
		bf.Add([]byte(uuid))
	}
	return bf, nil
}
