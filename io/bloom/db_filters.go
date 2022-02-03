package bloom

import (
	"database/sql"
	"fmt"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gnmatcher/io/dbase"
	log "github.com/sirupsen/logrus"
)

func (em *exactMatcher) filtersFromDB(path string) error {
	db := dbase.NewDB(em.cfg)
	log.Println("Importing lookup data for simple canonicals.")
	cFilter, cSize, err := createFilter(db, "canonicals")
	if err != nil {
		return err
	}
	em.filters = &bloomFilters{
		canonical:     cFilter,
		canonicalSize: cSize,
	}
	saveFilters(path, em.filters)
	return db.Close()
}

func createFilter(db *sql.DB,
	table string) (*baseBloomfilter.Bloomfilter, uint, error) {
	var err error
	var nilFilter *baseBloomfilter.Bloomfilter

	size, err := getFilterSize(db, table)
	if err != nil {
		return nilFilter, 0, err
	}
	return newFilter(db, table, size)
}

func getFilterSize(db *sql.DB, table string) (uint, error) {
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if table == "name_strings" {
		q = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE virus = TRUE", table)
	}
	var num uint
	row := db.QueryRow(q)
	if err := row.Scan(&num); err != nil {
		return 0, err
	}
	return num, nil
}

func newFilter(db *sql.DB, table string,
	filterSize uint) (*baseBloomfilter.Bloomfilter, uint, error) {
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
		return bf, filterSize, err
	}
	for rows.Next() {
		if err := rows.Scan(&uuid); err != nil {
			return bf, filterSize, err
		}
		bf.Add([]byte(uuid))
	}
	return bf, filterSize, nil
}
