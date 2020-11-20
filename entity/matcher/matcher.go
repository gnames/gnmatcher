package matcher

import (
	"strings"
	"sync"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	vlib "github.com/gnames/gnlib/domain/entity/verifier"
	"github.com/gnames/gnmatcher/entity/exact"
	"github.com/gnames/gnmatcher/entity/fuzzy"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gogna/gnparser"
	"gitlab.com/gogna/gnparser/pb"
)

const (
	// MaxNMaxNamesNum is the largest number of names that can be processed
	// per request. If input contains more names, it will be truncated.
	MaxNamesNum = 10_000

	// JobsNum determines number of go-routines to run during matching.
	JobsNum = 16
)

var nilResult *mlib.Match

type matcher struct {
	exactMatcher exact.ExactMatcher
	fuzzyMatcher fuzzy.FuzzyMatcher
}

// NewMatcher returns Matcher object. It takes interfaces to ExactMatcher
// and FuzzyMatcher.
func NewMatcher(em exact.ExactMatcher, fm fuzzy.FuzzyMatcher) matcher {
	return matcher{exactMatcher: em, fuzzyMatcher: fm}
}

func (m matcher) Init() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		m.exactMatcher.Init()
	}()
	go func() {
		defer wg.Done()
		m.fuzzyMatcher.Init()
	}()
	wg.Wait()
}

type nameIn struct {
	index int
	name  string
}

type matchOut struct {
	index int
	match *mlib.Match
}

func (m matcher) MatchNames(names []string) []*mlib.Match {
	names = truncateNamesToMaxNumber(names)
	chIn := make(chan nameIn)
	chOut := make(chan matchOut)
	var wgIn sync.WaitGroup
	var wgOut sync.WaitGroup
	wgIn.Add(JobsNum)
	wgOut.Add(1)

	names = truncateNamesToMaxNumber(names)
	log.Printf("Processing %d names.", len(names))
	res := make([]*mlib.Match, len(names))

	go loadNames(chIn, names)
	for i := 0; i < JobsNum; i++ {
		go m.matchWorker(chIn, chOut, &wgIn)
	}

	go func() {
		defer wgOut.Done()
		for r := range chOut {
			res[r.index] = r.match
		}
	}()

	wgIn.Wait()
	close(chOut)
	wgOut.Wait()
	return res
}

// matchWorker takes name-strings from chIn channel, matches them
// and sends results to chOut channel.
func (m matcher) matchWorker(chIn <-chan nameIn,
	chOut chan<- matchOut, wg *sync.WaitGroup) {
	parser := gnparser.NewGNparser()
	defer wg.Done()
	var matchResult *mlib.Match

	for tsk := range chIn {
		ns, parsed := newNameString(parser, tsk.name)
		if parsed.Parsed {
			if abbrResult := detectAbbreviated(parsed); abbrResult != nil {
				chOut <- matchOut{index: tsk.index, match: abbrResult}
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
		chOut <- matchOut{index: tsk.index, match: matchResult}
	}
}

func loadNames(chIn chan<- nameIn, names []string) {
	for i, name := range names {
		chIn <- nameIn{index: i, name: name}
	}
	close(chIn)
}

func truncateNamesToMaxNumber(names []string) []string {
	if len(names) > MaxNamesNum {
		log.Warnf("Too many names, truncating list to %d entries.", MaxNamesNum)
		names = names[0:MaxNamesNum]
	}
	return names
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

func emptyResult(ns nameString) *mlib.Match {
	return &mlib.Match{
		ID:        ns.ID,
		Name:      ns.Name,
		MatchType: vlib.NoMatch,
	}
}
