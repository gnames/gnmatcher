package bloom

import (
	"database/sql"
	"log/slog"

	"github.com/devopsfaith/bloomfilter"
	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/gnames/gnmatcher/internal/io/dbase"
)

func (em *exactMatcher) filtersFromDB(path string) error {
	db, err := dbase.NewDB(em.cfg)
	if err != nil {
		return err
	}

	slog.Info("Importing lookup data for stemmed canonicals")
	cFilter, cSize, err := createFilter(db)
	if err != nil {
		return err
	}
	em.filters = &bloomFilters{
		canonicalStem: cFilter,
		canonicalSize: cSize,
	}
	saveFilters(path, em.filters)
	return db.Close()
}

func createFilter(
	db *sql.DB,
) (*baseBloomfilter.Bloomfilter, uint, error) {
	var err error
	var nilFilter *baseBloomfilter.Bloomfilter

	size, err := getFilterSize(db)
	if err != nil {
		return nilFilter, 0, err
	}
	return newFilter(db, size)
}

func getFilterSize(db *sql.DB) (uint, error) {
	q := "SELECT COUNT(*) FROM canonical_stems " +
		"WHERE NOT all_canonicals_misspelled"
	var num uint
	row := db.QueryRow(q)
	if err := row.Scan(&num); err != nil {
		return 0, err
	}
	return num, nil
}

func newFilter(
	db *sql.DB,
	filterSize uint,
) (*baseBloomfilter.Bloomfilter, uint, error) {
	var uuid string
	cfg := bloomfilter.Config{
		N:        filterSize,
		P:        0.00001,
		HashName: bloomfilter.HASHER_OPTIMAL,
	}
	bf := baseBloomfilter.New(cfg)

	// Stems flagged by gndb are hidden from matching, otherwise a misspelling
	// matches itself exactly, and fuzzy matching to the correct name never
	// runs.
	rows, err := db.Query(
		"SELECT id FROM canonical_stems WHERE NOT all_canonicals_misspelled",
	)
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
