package matcher

import (
	"strings"
	"sync"

	"github.com/dgraph-io/badger/v2"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnlib/sys"
	"github.com/gnames/gnmatcher/bloom"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/fuzzy"
	"github.com/gnames/gnmatcher/stemskv"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser"
	"gitlab.com/gogna/gnparser/pb"
)

var (
	nilResult *mlib.Match
)

// Matcher contains data and functions necessary for exact, fuzzy and partial
// matching of scientific names.
type Matcher struct {
	Config  config.Config
	Filters *bloom.Filters
	KeyVal  *badger.DB
	mlib.FuzzyMatcher
}

// MatchTask contains a name to be matched and an index where it should be
// located in an array.
type MatchTask struct {
	Index int
	Name  string
}

type MatchResult struct {
	Index int
	Match *mlib.Match
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
	m.FuzzyMatcher = NewFuzzyMatcher(trie, kv)

	return m
}

// MatchWorker takes name-strings from chIn channel, matches them
// and sends results to chOut channel.
func (m Matcher) MatchWorker(chIn <-chan MatchTask,
	chOut chan<- MatchResult, wg *sync.WaitGroup) {
	parser := gnparser.NewGNparser()
	defer wg.Done()
	var matchResult *mlib.Match

	for tsk := range chIn {
		ns, parsed := newNameString(parser, tsk.Name)
		if parsed.Parsed {
			if abbrResult := detectAbbreviated(parsed); abbrResult != nil {
				chOut <- MatchResult{Index: tsk.Index, Match: abbrResult}
				continue
			}
			matchResult = m.match(ns)
		} else {
			matchResult = m.matchVirus(ns)
		}
		if matchResult == nil {
			matchResult = m.matchFuzzy(ns.Canonical, ns.CanonicalStem, ns)
		}
		if matchResult == nil {
			matchResult = m.matchPartial(ns)
		}
		chOut <- MatchResult{Index: tsk.Index, Match: matchResult}
	}
}

// detectAbbreviated checks if parsed name is abbreviated. If name is not
// abbreviated the function returns nil. If it is abbreviated, it returns
// result with the MatchType 'NONE'.
func detectAbbreviated(parsed *pb.Parsed) *mlib.Match {
	if parsed.Quality != int32(3) {
		return nilResult
	}
	for _, v := range parsed.QualityWarning {
		if strings.HasPrefix(v.Message, "Abbreviated") {
			return &mlib.Match{
				ID:        parsed.Id,
				Name:      parsed.Verbatim,
				MatchType: vlib.NoMatch,
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

func emptyResult(ns nameString) *mlib.Match {
	return &mlib.Match{
		ID:        ns.ID,
		Name:      ns.Name,
		MatchType: vlib.NoMatch,
	}
}
