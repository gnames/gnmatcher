// package gnmatcher provides the main use-case of the project, which is
// matching of possible name-strings to scientific names registered in a
// variety of biodiversity databases.
//
// The goal of the project is to return back matched canonical forms of
// scientific names by tens of thousands a second, making it possible to work
// with hundreds of millions/billions of name-string matching events.
//
// The package is intended to be used by long-running services, because it
// takes a few seconds to initialized its lookup data structures.
package gnmatcher

import (
	"fmt"

	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnmatcher/config"
	"github.com/gnames/gnmatcher/matcher"
)

// GNMatcher is a public API to the project functionality.
type GNMatcher interface {
	// MatchNames take a slice of scientific name-strings and return back
	// matches to canonical forms of known scientific names. The following
	// matches are attempted:
	// - Exact string match for viruses
	// - Exact match of the name-string's canonical form
	// - Fuzzy match of the canonical form
	// - Partial match of the canonical form where the middle parts of the name
	//   or last elements of the name are removed.
	// - Partial fuzzy match of the canonical form.
	//
	// The resulting output does provide canonical forms, but not the sources
	// where they are registered.
	//
	MatchNames(names []string) []*mlib.Match
}

func Example() {
	// Note that it takes several minutes to initialize lookup data structures.
	// Requirement for initialization: Postgresql database with loaded
	// http://opendata.globalnames.org/dumps/gnames-latest.sql.gz
	//
	// If data are imported already, it still takes several seconds to
	// load lookup data into memory.
	cnf := config.NewConfig()
	m := matcher.NewMatcher(cnf)
	gnm := NewGNMatcher(m)
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
