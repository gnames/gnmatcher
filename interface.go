// package gnmatcher provides the main use-case of the project, which is
// matching of possible name-strings to scientific names registered in a
// variety of biodiversity databases.
//
// The goal of the project is to return matched canonical forms of
// scientific names by tens of thousands a second, making it possible to work
// with hundreds of millions/billions of name-string matching events.
//
// The package is intended to be used by long-running services because it
// takes a few seconds/minutes to initialize its lookup data structures.
package gnmatcher

import (
	"github.com/gnames/gnlib/ent/gnvers"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/config"
)

// GNmatcher is a public API to the project functionality.
type GNmatcher interface {
	// MatchNames takes a slice of scientific name-strings with options and
	// returns back matches to canonical forms of known scientific names. The
	// following matches are attempted:
	// - Exact string match for viruses
	// - Exact match of the name-string's canonical form
	// - Fuzzy match of the canonical form
	// - Partial match of the canonical form where the middle parts of the name
	//   or last elements of the name are removed.
	// - Partial fuzzy match of the canonical form.
	//
	// In case if a name is determined as a "virus" (a non-celular entity like
	// virus, prion, plasmid etc.), It is not matched, and returned back
	// to be found in a database.
	//
	// The resulting output does provide canonical forms, but not the sources
	// where they are registered.
	MatchNames(names []string, opts ...config.Option) []mlib.Output

	// GetConfig provides configuration object of GNmatcher.
	GetConfig() config.Config

	// GetVersion returns version number and build timestamp.
	GetVersion() gnvers.Version
}
