// package trie implements FuzzyMatcher interface that is responsible for
// fuzzy-matching strings to canonical forms of scientific names.
package trie

import (
	"bytes"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
	"github.com/dvirsky/levenshtein"
	"github.com/gnames/gnfmt"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/internal/ent/fuzzy"
	"github.com/gnames/gnmatcher/internal/io/dbase"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/gnames/gnsys"
)

const trieFile = "stem.trie"

type fuzzyMatcher struct {
	cfg     config.Config
	trie    *levenshtein.MinTree
	kvStems *badger.DB
	encoder gnfmt.Encoder
}

// New takes configuration and returns back FuzzyMatcher object
// responsible for fuzzy-matching strings to canonical forms of scientific
// names.
func New(cfg config.Config) fuzzy.FuzzyMatcher {
	fm := fuzzyMatcher{cfg: cfg, encoder: gnfmt.GNgob{}}
	return &fm
}

func (fm *fuzzyMatcher) Init() error {
	var err error
	fm.prepareDirs()
	db, err := dbase.NewDB(fm.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	fm.trie, err = getTrie(fm.cfg.TrieDir(), db)
	if err != nil {
		return err
	}

	err = initStemsKV(fm.cfg.StemsDir(), db)
	if err != nil {
		return err
	}

	fm.kvStems, err = connectKeyVal(fm.cfg.StemsDir())
	if err != nil {
		return err
	}

	return nil
}

// SetConfig updates configuration of the matcher.
func (fm *fuzzyMatcher) SetConfig(cfg config.Config) {
	fm.cfg = cfg
}

func (fm *fuzzyMatcher) MatchStem(stem string) []string {
	return fm.trie.FuzzyMatches(stem, fm.cfg.MaxEditDist)
}

func (fm *fuzzyMatcher) MatchStemExact(stem string) bool {
	matches := fm.trie.FuzzyMatches(stem, 0)
	return len(matches) > 0
}

func (fm *fuzzyMatcher) StemToMatchItems(
	stem string,
) ([]mlib.MatchItem, error) {
	var res []mlib.MatchItem

	bs, err := getValue(fm.kvStems, stem)
	if err != nil {
		return nil, err
	}
	misGob := bytes.NewBuffer(bs)
	err = fm.encoder.Decode(misGob.Bytes(), &res)
	if err != nil {
		slog.Error("Decode failed", "stem", stem, "error", err)
		return res, err
	}
	return res, nil
}

// getTrie generates an in-memory trie for levenshtein automata. Such tree
// can either be constructed from database or from a dump file. The tree
// consists stemmed canonical forms of _gnames_ database.
func getTrie(triePath string, db *sql.DB) (*levenshtein.MinTree, error) {
	var trie *levenshtein.MinTree
	trie, err := getCachedTrie(triePath)
	if err == nil {
		slog.Info("Trie data is rebuilt from cache")
		return trie, nil
	}

	trie, err = populateAndSaveTrie(db, triePath)
	if err != nil {
		slog.Error("Cannot build trie from db", "error", err)
		return nil, err
	}
	return trie, nil
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
	slog.Info("Getting trie data from database")
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
	slog.Info("Building trie and saving it to disk")
	path := filepath.Join(triePath, trieFile)
	w, err := os.Create(path)
	if err != nil {
		return trie, err
	}
	trie, err = levenshtein.NewMinTreeWrite(names, w)
	if err != nil {
		return trie, err
	}
	slog.Info("Trie is created")
	return trie, nil
}

func (fm fuzzyMatcher) prepareDirs() {
	slog.Info("Preparing dirs for trie and stems key-value store")
	dirs := []string{fm.cfg.TrieDir(), fm.cfg.StemsDir()}
	for _, dir := range dirs {
		err := gnsys.MakeDir(dir)
		if err != nil {
			slog.Error("Cannot create directory", "path", dir, "error", err)
		}
	}
}
