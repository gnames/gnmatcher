package matcher

import (
	"sync"

	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/ent/exact"
	"github.com/gnames/gnmatcher/ent/fuzzy"
	"github.com/gnames/gnmatcher/ent/virus"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/parsed"
	"github.com/rs/zerolog/log"
)

const (
	// MaxMaxNamesNum is the largest number of names that can be processed
	// per request. If input contains more names, it will be truncated.
	MaxNamesNum = 10_000
)

type matcher struct {
	exactMatcher exact.ExactMatcher
	fuzzyMatcher fuzzy.FuzzyMatcher
	virusMatcher virus.VirusMatcher
	jobsNum      int
}

// NewMatcher returns Matcher object. It takes interfaces to ExactMatcher
// and FuzzyMatcher.
func NewMatcher(
	em exact.ExactMatcher,
	fm fuzzy.FuzzyMatcher,
	vm virus.VirusMatcher,
	j int) Matcher {
	return matcher{
		exactMatcher: em,
		fuzzyMatcher: fm,
		virusMatcher: vm,
		jobsNum:      j}
}

func (m matcher) Init() {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		m.exactMatcher.Init()
	}()
	go func() {
		defer wg.Done()
		m.fuzzyMatcher.Init()
	}()
	go func() {
		defer wg.Done()
		m.virusMatcher.Init()
	}()
	wg.Wait()
}

type nameIn struct {
	index int
	name  string
}

type matchOut struct {
	index int
	match mlib.Match
}

func (m matcher) MatchNames(names []string) []mlib.Match {
	names = truncateNamesToMaxNumber(names)
	chIn := make(chan nameIn)
	chOut := make(chan matchOut)
	var wgIn sync.WaitGroup
	var wgOut sync.WaitGroup
	wgIn.Add(m.jobsNum)
	wgOut.Add(1)

	names = truncateNamesToMaxNumber(names)
	res := make([]mlib.Match, len(names))

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
	cfg := gnparser.NewConfig()
	parser := gnparser.New(cfg)
	defer wg.Done()

	for tsk := range chIn {
		var matchResult *mlib.Match
		ns, prsd := newNameString(parser, tsk.name)
		if prsd.Parsed {
			if abbrResult := detectAbbreviated(prsd); abbrResult != nil {
				chOut <- matchOut{index: tsk.index, match: *abbrResult}
				continue
			}
			matchResult = m.matchStem(ns)
			if ns.Cardinality < 2 {
				if matchResult == nil {
					matchResult = emptyResult(ns)
				}
				chOut <- matchOut{index: tsk.index, match: *matchResult}
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
		chOut <- matchOut{index: tsk.index, match: *matchResult}
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
	if l := len(names); l > MaxNamesNum {
		log.Warn().Int("namesNum", l).
			Str("example", names[0]).
			Msgf("Too many names, truncating list to %d entries", MaxNamesNum)
		names = names[0:MaxNamesNum]
	}
	return names
}

// detectAbbreviated checks if parsed name is abbreviated. If name is not
// abbreviated the function returns nil. If it is abbreviated, it returns
// result with the MatchType 'NONE'.
func detectAbbreviated(prsd *parsed.Parsed) *mlib.Match {
	// Abbreviations belong to ParseQuality 4
	if prsd.ParseQuality != 4 {
		return nil
	}
	for _, v := range prsd.QualityWarnings {
		if v.Warning == parsed.GenusAbbrWarn {
			return &mlib.Match{
				ID:        prsd.VerbatimID,
				Name:      prsd.Verbatim,
				MatchType: vlib.NoMatch,
			}
		}
	}
	return nil
}

func (m matcher) exactStemMatches(stemUUID, stem string) []mlib.MatchItem {
	if !m.exactMatcher.MatchCanonicalID(stemUUID) {
		return nil
	}
	if m.fuzzyMatcher.MatchStemExact(stem) {
		res := m.fuzzyMatcher.StemToMatchItems(stem)
		return res
	}
	return nil
}

func emptyResult(ns nameString) *mlib.Match {
	return &mlib.Match{
		ID:        ns.ID,
		Name:      ns.Name,
		MatchType: vlib.NoMatch,
	}
}
