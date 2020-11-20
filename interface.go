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
	"github.com/gnames/gnlib/domain/entity/gn"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
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

	gn.Versioner
}
