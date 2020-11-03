package gnmatcher

import (
	"fmt"
	"sync"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnmatcher/matcher"
	log "github.com/sirupsen/logrus"
)

// MaxMaxNamesNumber is the upper limit of the number of name-strings the
// MatchNames function can process. If the number is higher, the list of
// name-strings will be truncated.
const MaxNamesNumber = 10_000

// GNMatcher contains high level methods for scientific name matching.
type GNMatcher struct {
	Matcher matcher.Matcher
}

// NewGNMatcher is a constructor for GNMatcher instance
func NewGNMatcher(m matcher.Matcher) GNMatcher {
	return GNMatcher{Matcher: m}
}

// MatchNames takes a list of name-strings and matches them against known
// by names aggregated in gnames database.
func (gnm GNMatcher) MatchNames(names []string) []*mlib.Match {
	m := gnm.Matcher
	cnf := m.Config

	chIn := make(chan matcher.MatchTask)
	chOut := make(chan matcher.MatchResult)
	var wgIn sync.WaitGroup
	var wgOut sync.WaitGroup
	wgIn.Add(cnf.JobsNum)
	wgOut.Add(1)

	names = truncateNamesToMaxNumber(names)
	log.Printf("Processing %d names.", len(names))
	res := make([]*mlib.Match, len(names))

	go loadNames(chIn, names)
	for i := 0; i < cnf.JobsNum; i++ {
		go m.MatchWorker(chIn, chOut, &wgIn)
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
