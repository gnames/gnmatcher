// package trie implements FuzzyMatcher interface that is responsible for
// fuzzy-matching strings to canonical forms of scientific names.
package trie

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
	"github.com/dvirsky/levenshtein"
	"github.com/gnames/gnfmt"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/ent/fuzzy"
	"github.com/gnames/gnmatcher/io/dbase"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
)

const trieFile = "stem.trie"

type fuzzyMatcher struct {
	cfg     config.Config
	trie    *levenshtein.MinTree
	keyVal  *badger.DB
	encoder gnfmt.Encoder
}

// New takes configuration and returns back FuzzyMatcher object
// responsible for fuzzy-matching strings to canonical forms of scientific
// names.
func New(cfg config.Config) fuzzy.FuzzyMatcher {
	fm := fuzzyMatcher{cfg: cfg, encoder: gnfmt.GNgob{}}
	return &fm
}

func (fm *fuzzyMatcher) Init() {
	fm.prepareDirs()
	db := dbase.NewDB(fm.cfg)
	defer db.Close()
	fm.trie = getTrie(fm.cfg.TrieDir(), db)
	initStemsKV(fm.cfg.StemsDir(), db)
	fm.keyVal = connectKeyVal(fm.cfg.StemsDir())
}

func (fm *fuzzyMatcher) MatchStem(stem string) []string {
	return fm.trie.FuzzyMatches(stem, fm.cfg.MaxEditDist)
}

func (fm *fuzzyMatcher) MatchStemExact(stem string) bool {
	matches := fm.trie.FuzzyMatches(stem, 0)
	return len(matches) > 0
}

func (fm *fuzzyMatcher) StemToMatchItems(stem string) []mlib.MatchItem {
	var res []mlib.MatchItem
	misGob := bytes.NewBuffer(getValue(fm.keyVal, stem))
	err := fm.encoder.Decode(misGob.Bytes(), &res)
	if err != nil {
		log.Warn().Err(err).
			Msgf("Decode in StemToMatchItems for '%s' failed", stem)
	}
	return res
}

// getTrie generates an in-memory trie for levenshtein automata. Such tree
// can either be constructed from database or from a dump file. The tree
// consists stemmed canonical forms of _gnames_ database.
func getTrie(triePath string, db *sql.DB) *levenshtein.MinTree {
	var trie *levenshtein.MinTree
	trie, err := getCachedTrie(triePath)
	if err == nil {
		log.Info().Msg("Trie data is rebuilt from cache")
		return trie
	}

	trie, err = populateAndSaveTrie(db, triePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot build trie from db")
	}
	return trie
}

func getTrieSize(db *sql.DB) (int, error) {
	q := "SELECT count(*) from canonical_stems"
	var num int
	row := db.QueryRow(q)
	if err := row.Scan(&num); err != nil {
		return 0, err
	}
	return num, nil
}

func getCachedTrie(triePath string) (*levenshtein.MinTree, error) {
	var trie *levenshtein.MinTree
	path := filepath.Join(triePath, trieFile)
	trieFile, err := os.Open(path)
	if err != nil {
		return trie, err
	}
	return levenshtein.LoadMinTree(trieFile)
}

func populateAndSaveTrie(db *sql.DB, triePath string) (*levenshtein.MinTree, error) {
	log.Info().Msg("Getting trie data from database")
	var trie *levenshtein.MinTree
	size, err := getTrieSize(db)
	if err != nil {
		return trie, err
	}
	names := make([]string, size)

	var name string
	q := "SELECT name FROM canonical_stems order by name"
	rows, err := db.Query(q)
	if err != nil {
		return trie, err
	}

	for rows.Next() {
		if err = rows.Scan(&name); err != nil {
			return trie, err
		}
		names = append(names, name)
	}
	log.Info().Msg("Building trie and saving it to disk")
	path := filepath.Join(triePath, trieFile)
	w, err := os.Create(path)
	if err != nil {
		return trie, err
	}
	trie, err = levenshtein.NewMinTreeWrite(names, w)
	if err != nil {
		return trie, err
	}
	log.Info().Msg("Trie is created")
	return trie, nil
}

func (fm fuzzyMatcher) prepareDirs() {
	log.Info().Msg("Preparing dirs for trie and stems key-value store")
	dirs := []string{fm.cfg.TrieDir(), fm.cfg.StemsDir()}
	for _, dir := range dirs {
		err := gnsys.MakeDir(dir)
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot create directory %s", dir)
		}
	}
}
