package gnmatcher_test

import (
	"fmt"
	"regexp"
	"testing"

	gnmatcher "github.com/gnames/gnmatcher/pkg"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	cfg := config.New()
	gnm := gnmatcher.New(cfg)
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
	gnm := gnmatcher.New(cfg)
	if err := gnm.Init(); err != nil {
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
