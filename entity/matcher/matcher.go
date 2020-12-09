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
)

var nilResult *mlib.Match

type matcher struct {
	exactMatcher exact.ExactMatcher
	fuzzyMatcher fuzzy.FuzzyMatcher
	jobsNum      int
}

// NewMatcher returns Matcher object. It takes interfaces to ExactMatcher
// and FuzzyMatcher.
func NewMatcher(em exact.ExactMatcher, fm fuzzy.FuzzyMatcher, j int) Matcher {
	return matcher{exactMatcher: em, fuzzyMatcher: fm, jobsNum: j}
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
	wgIn.Add(m.jobsNum)
	wgOut.Add(1)

	names = truncateNamesToMaxNumber(names)
	log.Infof("Processing %d names.", len(names))
	res := make([]*mlib.Match, len(names))

	go loadNames(chIn, names)
	for i := 0; i < m.jobsNum; i++ {
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
func (m matcher) matchWorker(
	chIn <-chan nameIn,
	chOut chan<- matchOut,
	wg *sync.WaitGroup,
) {
	parser := gnparser.NewGNparser()
	defer wg.Done()

	for tsk := range chIn {
		var matchResult *mlib.Match
		ns, parsed := newNameString(parser, tsk.name)
		if parsed.Parsed {
			if abbrResult := detectAbbreviated(parsed); abbrResult != nil {
				chOut <- matchOut{index: tsk.index, match: abbrResult}
				continue
			}
			matchResult = m.match(ns)
			if ns.Cardinality < 2 {
				if matchResult == nil {
					matchResult = emptyResult(ns)
				}
				chOut <- matchOut{index: tsk.index, match: matchResult}
				continue
			}
		} else if ns.IsVirus {
			matchResult = m.matchVirus(ns)
		}
		if matchResult == nil {
			matchResult = m.matchFuzzy(ns.Canonical, ns.CanonicalStem, ns)
		}
		if matchResult == nil {
			matchResult = m.matchPartial(ns, parser)
		}
		chOut <- matchOut{index: tsk.index, match: matchResult}
	}
}

func loadNames(chIn chan<- nameIn, names []string) {
	for i, name := range names {
		ni := nameIn{index: i, name: name}
		chIn <- ni
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

func (m matcher) isExactMatch(uuid, stem string) bool {
	return m.exactMatcher.MatchCanonicalID(uuid) &&
		m.fuzzyMatcher.MatchStemExact(stem)
}

func emptyResult(ns nameString) *mlib.Match {
	return &mlib.Match{
		ID:        ns.ID,
		Name:      ns.Name,
		MatchType: vlib.NoMatch,
	}
}
