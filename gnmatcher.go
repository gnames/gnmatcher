package gnmatcher

import (
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v2"
	"github.com/gnames/gnmatcher/matcher"
	"github.com/gnames/gnmatcher/model"
	"github.com/gnames/gnmatcher/stemskv"
	log "github.com/sirupsen/logrus"
)

// MaxMaxNamesNumber is the upper limit of the number of name-strings the
// MatchNames function can process. If the number is higher, the list of
// name-strings will be truncated.
const MaxNamesNumber = 10_000

// GNMatcher contains high level methods for scientific name matching.
type GNMatcher struct {
	Matcher matcher.Matcher
	KV      *badger.DB
}

// NewGNMatcher is a constructor for GNMatcher instance
func NewGNMatcher(m matcher.Matcher) GNMatcher {
	path := m.Config.StemsDir()
	kv := stemskv.ConnectKeyVal(path)
	return GNMatcher{Matcher: m, KV: kv}
}

// MatchNames takes a list of name-strings and matches them against known
// by names aggregated in gnames database.
func (gnm GNMatcher) MatchNames(names []string) []*model.Match {
	m := gnm.Matcher
	cnf := m.Config
	kv := gnm.KV

	chIn := make(chan matcher.MatchTask)
	chOut := make(chan matcher.MatchResult)
	var wgIn sync.WaitGroup
	var wgOut sync.WaitGroup
	wgIn.Add(cnf.JobsNum)
	wgOut.Add(1)

	names = truncateNamesToMaxNumber(names)
	log.Printf("Processing %d names.", len(names))
	res := make([]*model.Match, len(names))

	go loadNames(chIn, names)
	for i := 0; i < cnf.JobsNum; i++ {
		go m.MatchWorker(chIn, chOut, &wgIn, kv)
	}

	go func() {
		defer wgOut.Done()
		for r := range chOut {
			res[r.Index] = r.Match
		}
	}()

	wgIn.Wait()
	close(chOut)
	wgOut.Wait()
	return res
}

func loadNames(chIn chan<- matcher.MatchTask, names []string) {
	for i, name := range names {
		chIn <- matcher.MatchTask{Index: i, Name: name}
	}
	close(chIn)
}

func truncateNamesToMaxNumber(names []string) []string {
	if len(names) > MaxNamesNumber {
		log.Warn(fmt.Sprintf("Too many names, truncating list to %d entries.",
			MaxNamesNumber))
		names = names[0:MaxNamesNumber]
	}
	return names
}
