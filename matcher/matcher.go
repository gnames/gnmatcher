package matcher

import (
	"strings"
	"sync"

	"github.com/dgraph-io/badger/v2"
	gn "github.com/gnames/gnames/domain/entity"
	"github.com/gnames/gnlib/sys"
	"github.com/gnames/gnmatcher/bloom"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/domain/entity"
	"github.com/gnames/gnmatcher/domain/usecase"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/stemskv"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser"
	"gitlab.com/gogna/gnparser/pb"
)

var (
	// GNUUID is a UUID seed made from 'globalnames.org' domain to generate
	// UUIDv5 identifiers.
	GNUUID    = uuid.NewV5(uuid.NamespaceDNS, "globalnames.org")
	nilResult *entity.Match
)

// Matcher contains data and functions necessary for exact, fuzzy and partial
// matching of scientific names.
type Matcher struct {
	Config  config.Config
	Filters *bloom.Filters
	KeyVal  *badger.DB
	usecase.FuzzyMatcher
}

// MatchTask contains a name to be matched and an index where it should be
// located in an array.
type MatchTask struct {
	Index int
	Name  string
}

type MatchResult struct {
	Index int
	Match *entity.Match
}

// NewMatcher creates a new instance of Matcher struct.
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

	log.Println("Initializing key-value store for stems.")
	stemskv.NewStemsKV(cnf.StemsDir(), db)

	path := m.Config.StemsDir()
	kv := stemskv.ConnectKeyVal(path)
	m.FuzzyMatcher = NewFuzzyMatcherTrie(trie, kv)

	return m
}

// MatchWorker takes name-strings from chIn channel, matches them
// and sends results to chOut channel.
func (m Matcher) MatchWorker(chIn <-chan MatchTask,
	chOut chan<- MatchResult, wg *sync.WaitGroup) {
	parser := gnparser.NewGNparser()
	defer wg.Done()
	var matchResult *entity.Match

	for tsk := range chIn {
		ns, parsed := NewNameString(parser, tsk.Name)
		if parsed.Parsed {
			if abbrResult := DetectAbbreviated(parsed); abbrResult != nil {
				chOut <- MatchResult{Index: tsk.Index, Match: abbrResult}
				continue
			}
			matchResult = m.Match(ns)
		} else {
			matchResult = m.MatchVirus(ns)
		}
		if matchResult == nil {
			matchResult = m.MatchFuzzy(ns.Canonical, ns.CanonicalStem, ns)
		}
		if matchResult == nil {
			matchResult = m.MatchPartial(ns)
		}
		chOut <- MatchResult{Index: tsk.Index, Match: matchResult}
	}
}

// DetectAbbreviated checks if parsed name is abbreviated. If name is not
// abbreviated the function returns nil. If it is abbreviated, it returns
// result with the MatchType 'NONE'.
func DetectAbbreviated(parsed *pb.Parsed) *entity.Match {
	if parsed.Quality != int32(3) {
		return nilResult
	}
	for _, v := range parsed.QualityWarning {
		if strings.HasPrefix(v.Message, "Abbreviated") {
			return &entity.Match{
				ID:        parsed.Id,
				Name:      parsed.Verbatim,
				MatchType: gn.NoMatch,
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

func emptyResult(ns NameString) *entity.Match {
	return &entity.Match{
		ID:        ns.ID,
		Name:      ns.Name,
		MatchType: gn.NoMatch,
	}
}
