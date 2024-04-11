package gnmatcher_test

import (
	"fmt"
	"regexp"
	"testing"

	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/internal/io/bloom"
	"github.com/gnames/gnmatcher/internal/io/trie"
	"github.com/gnames/gnmatcher/internal/io/virusio"
	gnmatcher "github.com/gnames/gnmatcher/pkg"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/stretchr/testify/assert"
)

type mockExactMatcher struct{}

func (em mockExactMatcher) Init() error { return nil }

func (em mockExactMatcher) SetConfig(cfg config.Config) {}

func (em mockExactMatcher) MatchCanonicalID(s string) bool {
	return false
}

func (em mockExactMatcher) MatchNameStringID(s string) bool {
	return false
}

type mockFuzzyMatcher struct{}

func (fm mockFuzzyMatcher) Init() error { return nil }

func (fm mockFuzzyMatcher) SetConfig(cfg config.Config) {}

func (fm mockFuzzyMatcher) MatchStem(s string) []string {
	var res []string
	return res
}

func (fm mockFuzzyMatcher) MatchStemExact(s string) bool {
	return true
}

func (fm mockFuzzyMatcher) StemToMatchItems(
	s string,
) ([]mlib.MatchItem, error) {
	var res []mlib.MatchItem
	return res, nil
}

type mockVirusMatcher struct{}

func (vm mockVirusMatcher) Init() error { return nil }

func (vm mockVirusMatcher) SetConfig(cfg config.Config) {}

func (vm mockVirusMatcher) MatchVirus(s string) ([]mlib.MatchItem, error) {
	return nil, nil
}

func (vm mockVirusMatcher) NameToBytes(s string) []byte { return nil }

func TestVersion(t *testing.T) {
	cfg := config.New()
	em := mockExactMatcher{}
	fm := mockFuzzyMatcher{}
	vm := mockVirusMatcher{}
	gnm, err := gnmatcher.New(em, fm, vm, cfg)
	assert.Nil(t, err)
	ver := gnm.GetVersion()
	verRegex := regexp.MustCompile(`^v[\d]+\.[\d]+\.[\d]+\+?`)
	assert.Regexp(t, verRegex, ver.Version)
	assert.Equal(t, "n/a", ver.Build)
}

func Example() {
	// Note that it takes several minutes to initialize lookup data structures.
	// Requirement for initialization: Postgresql database with loaded
	// http://opendata.globalnames.org/dumps/gnames-latest.sql.gz
	//
	// If data are imported already, it still takes several seconds to
	// load lookup data into memory.
	cfg := config.New()
	em := bloom.New(cfg)
	fm := trie.New(cfg)
	vm := virusio.New(cfg)
	gnm, err := gnmatcher.New(em, fm, vm, cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	res := gnm.MatchNames([]string{"Pomatomus saltator", "Pardosa moesta"})
	for _, match := range res.Matches {
		fmt.Println(match.Name)
		fmt.Println(match.MatchType)
		for _, item := range match.MatchItems {
			fmt.Println(item.MatchStr)
			fmt.Println(item.EditDistance)
		}
	}
}
