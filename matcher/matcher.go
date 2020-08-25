package matcher

import (
	"strings"

	"github.com/dvirsky/levenshtein"
	"github.com/gnames/gnmatcher/bloom"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/protob"
	"github.com/gnames/gnmatcher/stemskv"
	"github.com/gnames/gnmatcher/sys"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser/pb"
)

var (
	GNUUID    = uuid.NewV5(uuid.NamespaceDNS, "globalnames.org")
	nilResult *protob.Result
)

type Matcher struct {
	Config  config.Config
	Filters *bloom.Filters
	Trie    *levenshtein.MinTree
}

func NewMatcher(cnf config.Config) Matcher {
	m := Matcher{Config: cnf}

	db := dbase.NewDB(cnf)
	defer db.Close()

	log.Println("Preparing dirs for bloom filters, trie, and stems key-value store.")
	m.prepareWorkDirs()

	log.Println("Initializing bloom filters.")
	filters := bloom.GetFilters(cnf.FiltersDir(), db)
	m.Filters = filters

	log.Println("Initializing levenshtein trie.")
	trie := fuzzy.GetTrie(cnf.TrieDir(), db)
	m.Trie = trie

	log.Println("Initializing key-value store for stems.")
	stemskv.NewStemsKV(cnf.StemsDir(), db)

	return m
}

// DetectAbbreviated checks if parsed name is abbreviated. If name is not
// abbreviated the function returns nil. If it is abbreviated, it returns
// result with the MatchType 'NONE'.
func DetectAbbreviated(parsed *pb.Parsed) *protob.Result {
	if parsed.Quality != int32(3) {
		return nilResult
	}
	for _, v := range parsed.QualityWarning {
		if strings.HasPrefix(v.Message, "Abbreviated") {
			return &protob.Result{
				Id:        parsed.Id,
				Name:      parsed.Verbatim,
				MatchType: protob.MatchType_NONE,
			}
		}
	}
	return nilResult
}

func (m Matcher) prepareWorkDirs() {
	cnf := m.Config
	dirs := []string{cnf.FiltersDir(), cnf.TrieDir(), cnf.StemsDir()}
	for _, dir := range dirs {
		err := sys.MakeDir(dir)
		if err != nil {
			log.Fatalf("Cannot create directory %s: %s.", dir, err)
		}
	}
}

func emptyResult(ns NameString) *protob.Result {
	return &protob.Result{
		Id:        ns.ID,
		Name:      ns.Name,
		MatchType: protob.MatchType_NONE,
	}
}
