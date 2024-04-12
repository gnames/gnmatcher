package matcher

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	mlib "github.com/gnames/gnlib/ent/matcher"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnmatcher/internal/ent/exact"
	"github.com/gnames/gnmatcher/internal/ent/fuzzy"
	"github.com/gnames/gnmatcher/internal/ent/virus"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/parsed"
	"golang.org/x/sync/errgroup"
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
	cfg          config.Config
}

// NewMatcher returns Matcher object. It takes interfaces to ExactMatcher
// and FuzzyMatcher.
func NewMatcher(
	em exact.ExactMatcher,
	fm fuzzy.FuzzyMatcher,
	vm virus.VirusMatcher,
	cfg config.Config) Matcher {
	return matcher{
		exactMatcher: em,
		fuzzyMatcher: fm,
		virusMatcher: vm,
		cfg:          cfg,
	}
}

func (m matcher) Init() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return m.exactMatcher.Init()
	})

	g.Go(func() error {
		return m.fuzzyMatcher.Init()
	})

	g.Go(func() error {
		return m.virusMatcher.Init()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

type nameIn struct {
	index int
	name  string
}

type matchOut struct {
	index int
	match mlib.Match
}

func (m matcher) MatchNames(
	names []string,
	opts ...config.Option,
) mlib.Output {
	chIn := make(chan nameIn)
	chOut := make(chan matchOut)
	var wgIn sync.WaitGroup
	var wgOut sync.WaitGroup
	wgIn.Add(m.cfg.JobsNum)
	wgOut.Add(1)

	for _, opt := range opts {
		opt(&m.cfg)
	}

	m.exactMatcher.SetConfig(m.cfg)
	m.fuzzyMatcher.SetConfig(m.cfg)
	m.virusMatcher.SetConfig(m.cfg)

	maxNum := MaxNamesNum

	names = truncateNamesToMaxNumber(names, maxNum)
	res := make([]mlib.Match, len(names))

	go loadNames(chIn, names)
	for range m.cfg.JobsNum {
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

	return m.prepareOutput(res)
}

func (m matcher) prepareOutput(ms []mlib.Match) mlib.Output {
	res := mlib.Output{
		Meta: mlib.Meta{
			NamesNum:                len(ms),
			WithSpeciesGroup:        m.cfg.WithSpeciesGroup,
			WithUninomialFuzzyMatch: m.cfg.WithUninomialFuzzyMatch,
			DataSources:             m.cfg.DataSources,
		},
	}
	for i := range ms {
		for ii := range ms[i].MatchItems {
			ms[i].MatchItems[ii].DataSources =
				m.convertDataSources(ms[i].MatchItems[ii])
		}
	}
	res.Matches = ms
	return res
}

func (m matcher) convertDataSources(mi mlib.MatchItem) []int {
	if len(m.cfg.DataSources) == 0 {
		res := make([]int, len(mi.DataSourcesMap))
		var i int
		for k := range mi.DataSourcesMap {
			res[i] = k
			i++
		}

		slices.Sort(res)
		return res
	}

	var res []int
	for _, i := range m.cfg.DataSources {
		if _, ok := mi.DataSourcesMap[i]; ok {
			res = append(res, i)
		}
	}
	slices.Sort(res)
	return res
}

// matchWorker takes name-strings from chIn channel, matches them
// and sends results to chOut channel.
func (m matcher) matchWorker(
	chIn <-chan nameIn,
	chOut chan<- matchOut,
	wg *sync.WaitGroup,
) error {
	var err error
	gnpCfg := gnparser.NewConfig()
	parser := gnparser.New(gnpCfg)
	defer wg.Done()

	for tsk := range chIn {
		var matchResult *mlib.Match
		ns, prsd := newNameString(parser, tsk.name)

		var nsSpGr *nameString
		if m.cfg.WithSpeciesGroup {
			nsSpGr = ns.spGroupString(parser)
		}

		if prsd.Parsed {
			if abbrResult := detectAbbreviated(prsd); abbrResult != nil {
				chOut <- matchOut{index: tsk.index, match: *abbrResult}
				continue
			}
			matchResult, err = m.matchStem(ns)
			if err != nil {
				return err
			}

			// if we are matching a whole species group, add group's
			// data to the match.
			if nsSpGr != nil {
				spGrResult, err := m.matchStem(*nsSpGr)
				if err != nil {
					return err
				}
				ns.fixSpGrResult(spGrResult)
				if matchResult == nil {
					matchResult = spGrResult
				} else if spGrResult != nil {
					matchResult.MatchItems = append(
						matchResult.MatchItems,
						spGrResult.MatchItems...,
					)
				}
			}

			if ns.Cardinality < 2 && !m.cfg.WithUninomialFuzzyMatch {
				if matchResult == nil {
					matchResult = emptyResult(ns)
				}
				chOut <- matchOut{index: tsk.index, match: *matchResult}
				continue
			}
		} else if ns.IsVirus {
			matchResult, err = m.matchVirus(ns)
			if err != nil {
				return err
			}
		}
		if matchResult == nil {
			matchResult, err = m.matchFuzzy(ns.Canonical, ns.CanonicalStem, ns)
			if err != nil {
				return err
			}
		}
		if matchResult == nil {
			matchResult, err = m.matchPartial(ns, parser)
			if err != nil {
				return err
			}
		}
		chOut <- matchOut{index: tsk.index, match: *matchResult}
	}
	return nil
}

func loadNames(chIn chan<- nameIn, names []string) {
	for i, name := range names {
		ni := nameIn{index: i, name: name}
		chIn <- ni
	}
	close(chIn)
}

func truncateNamesToMaxNumber(names []string, maxNum int) []string {
	if l := len(names); l > maxNum {
		slog.Warn(
			"Too many names, truncating the list.",
			"names-number", l, "max-number", maxNum)
		names = names[0:maxNum]
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

func (m matcher) exactStemMatches(
	stemUUID, stem string,
) ([]mlib.MatchItem, error) {
	if !m.exactMatcher.MatchCanonicalID(stemUUID) {
		return nil, nil
	}
	if m.fuzzyMatcher.MatchStemExact(stem) {
		res, err := m.fuzzyMatcher.StemToMatchItems(stem)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, nil
}

func emptyResult(ns nameString) *mlib.Match {
	return &mlib.Match{
		ID:        ns.ID,
		Name:      ns.Name,
		MatchType: vlib.NoMatch,
	}
}

func (m matcher) filterDataSources(mis []mlib.MatchItem) []mlib.MatchItem {
	if len(mis) == 0 || len(m.cfg.DataSources) == 0 {
		return mis
	}

	var res []mlib.MatchItem
	for i := range mis {
		dataSourcesMap := m.matchDataSources(mis[i])
		if len(dataSourcesMap) > 0 {
			mis[i].DataSourcesMap = dataSourcesMap
			res = append(res, mis[i])

		}
	}
	return res
}

func (m matcher) matchDataSources(mi mlib.MatchItem) map[int]struct{} {
	res := make(map[int]struct{})
	var ok bool
	for _, dsID := range m.cfg.DataSources {
		if _, ok = mi.DataSourcesMap[dsID]; ok {
			res[dsID] = struct{}{}
		}
	}
	return res
}
