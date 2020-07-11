package gnmatcher

import (
	"log"
	"path/filepath"

	"github.com/dvirsky/levenshtein"
	"github.com/gnames/gnmatcher/bloom"
	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/sys"
)

// GNmatcher keeps most general configuration settings and high level
// methods for scientific name matching.
type GNmatcher struct {
	WorkDir string
	JobsNum int
	dbase.Dbase
	Filters *bloom.Filters
	Trie    *levenshtein.MinTree
}

// NewGNmatcher is a constructor for GNmatcher instance
func NewGNmatcher(opts ...Option) (GNmatcher, error) {
	gnm := GNmatcher{
		WorkDir: "/tmp/gnmatcher",
		JobsNum: 4,
		Dbase:   dbase.NewDbase(),
	}
	for _, opt := range opts {
		opt(&gnm)
	}
	err := gnm.CreateWorkDirs()
	if err != nil {
		return gnm, err
	}

	log.Println("Initializing bloom filters...")
	filters, err := bloom.GetFilters(gnm.FiltersDir(), gnm.Dbase)
	if err != nil {
		return gnm, err
	}
	gnm.Filters = filters
	log.Println("Initializing levenshtein trie...")
	trie, err := fuzzy.GetTrie(gnm.TrieDir(), gnm.Dbase)
	if err != nil {
		return gnm, err
	}
	gnm.Trie = trie
	return gnm, nil
}

func (gnm GNmatcher) TrieDir() string {
	return filepath.Join(gnm.WorkDir, "levenshein")
}

func (gnm GNmatcher) FiltersDir() string {
	return filepath.Join(gnm.WorkDir, "bloom")
}

func (gnm GNmatcher) CreateWorkDirs() error {
	err := sys.MakeDir(gnm.FiltersDir())
	if err != nil {
		return err
	}
	return sys.MakeDir(gnm.TrieDir())
}

// Option is a type of all options for GNmatcher.
type Option func(gnm *GNmatcher)

// OptWorkDir sets a directory for key-value stores and temporary files.
func OptWorkDir(d string) Option {
	return func(gnm *GNmatcher) {
		gnm.WorkDir = d
	}
}

// OptJobsNum sets number of concurrent jobs to run for parallel tasks.
func OptJobsNum(i int) Option {
	return func(gnm *GNmatcher) {
		gnm.JobsNum = i
	}
}

// OptPgHost sets the host of gnames database
func OptPgHost(h string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgHost = h
	}
}

// OptPgUser sets the user of gnnames database
func OptPgUser(u string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgUser = u
	}
}

// OptPgPass sets the password to access gnnames database
func OptPgPass(p string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgPass = p
	}
}

// OptPgPort sets the port for gnames database
func OptPgPort(p int) Option {
	return func(gnm *GNmatcher) {
		gnm.PgPort = p
	}
}

// OptPgDB sets the name of gnames database
func OptPgDB(n string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgDB = n
	}
}
