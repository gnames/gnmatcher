package gnmatcher

import (
	"path/filepath"

	"github.com/dvirsky/levenshtein"
	"github.com/gnames/gnmatcher/bloom"
	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/stemskv"
	"github.com/gnames/gnmatcher/sys"
	log "github.com/sirupsen/logrus"
)

// GNMatcher keeps most general configuration settings and high level
// methods for scientific name matching.
type GNMatcher struct {
	WorkDir     string
	NatsURI     string
	JobsNum     int
	MaxEditDist int
	GNamesDB    dbase.Dbase
	Filters     *bloom.Filters
	Trie        *levenshtein.MinTree
}

// NewGNMatcher is a constructor for GNMatcher instance
func NewGNMatcher(cnf Config) (GNMatcher, error) {
	gnm := GNMatcher{
		WorkDir:     cnf.WorkDir,
		NatsURI:     cnf.NatsURI,
		JobsNum:     cnf.JobsNum,
		MaxEditDist: cnf.MaxEditDist,
		GNamesDB:    cnf.GNamesDB,
	}
	err := gnm.CreateWorkDirs()
	if err != nil {
		return gnm, err
	}

	log.Println("Initializing bloom filters.")
	filters, err := bloom.GetFilters(gnm.FiltersDir(), gnm.GNamesDB)
	if err != nil {
		return gnm, err
	}
	gnm.Filters = filters
	log.Println("Initializing levenshtein trie.")
	trie, err := fuzzy.GetTrie(gnm.TrieDir(), gnm.GNamesDB)
	if err != nil {
		return gnm, err
	}
	gnm.Trie = trie

	log.Println("Initializing key-value store for stems.")
	stemskv.NewStemsKV(gnm.StemsDir(), gnm.GNamesDB)

	return gnm, nil
}

func (gnm GNMatcher) TrieDir() string {
	return filepath.Join(gnm.WorkDir, "levenshein")
}

func (gnm GNMatcher) FiltersDir() string {
	return filepath.Join(gnm.WorkDir, "bloom")
}

func (gnm GNMatcher) StemsDir() string {
	return filepath.Join(gnm.WorkDir, "stems-kv")
}

func (gnm GNMatcher) CreateWorkDirs() error {
	err := sys.MakeDir(gnm.FiltersDir())
	if err != nil {
		return err
	}
	return sys.MakeDir(gnm.TrieDir())
}
