package gnmatcher_test

import (
	"fmt"
	"regexp"
	"testing"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnmatcher"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/io/bloom"
	"github.com/gnames/gnmatcher/io/trie"
	"github.com/stretchr/testify/assert"
)

type mockExactMatcher struct{}

func (em mockExactMatcher) Init() {}

func (em mockExactMatcher) MatchCanonicalID(s string) bool {
	return false
}

func (em mockExactMatcher) MatchNameStringID(s string) bool {
	return false
}

type mockFuzzyMatcher struct{}

func (fm mockFuzzyMatcher) Init() {}

func (fm mockFuzzyMatcher) MatchStem(s string) []string {
	var res []string
	return res
}

func (fm mockFuzzyMatcher) MatchStemExact(s string) bool {
	return true
}

func (fm mockFuzzyMatcher) StemToMatchItems(s string) []mlib.MatchItem {
	var res []mlib.MatchItem
	return res
}

func TestVersion(t *testing.T) {
	em := mockExactMatcher{}
	fm := mockFuzzyMatcher{}
	gnm := gnmatcher.NewGNMatcher(em, fm)
	ver := gnm.GetVersion()
	verRegex := regexp.MustCompile(`^v[\d]+\.[\d]+\.[\d]+$`)
	assert.Regexp(t, verRegex, ver.Version)
	assert.Equal(t, ver.Build, "n/a")
}

func Example() {
	// Note that it takes several minutes to initialize lookup data structures.
	// Requirement for initialization: Postgresql database with loaded
	// http://opendata.globalnames.org/dumps/gnames-latest.sql.gz
	//
	// If data are imported already, it still takes several seconds to
	// load lookup data into memory.
	cfg := config.NewConfig()
	em := bloom.NewExactMatcher(cfg)
	fm := trie.NewFuzzyMatcher(cfg)
	gnm := gnmatcher.NewGNMatcher(em, fm)
	res := gnm.MatchNames([]string{"Pomatomus saltator", "Pardosa moesta"})
	for _, match := range res {
		fmt.Println(match.Name)
		fmt.Println(match.MatchType)
		for _, item := range match.MatchItems {
			fmt.Println(item.MatchStr)
			fmt.Println(item.EditDistance)
		}
	}
}
