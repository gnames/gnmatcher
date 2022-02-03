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
)

// GNmatcher is a public API to the project functionality.
type GNmatcher interface {
	NameMatcher
	WebLogger
}

type NameMatcher interface {
	// MatchNames takes a slice of scientific name-strings and returns back
	// matches to canonical forms of known scientific names. The following
	// matches are attempted:
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
	MatchNames(names []string) []mlib.Match

	// GetVersion returns version number and build timestamp.
	GetVersion() gnvers.Version
}

type WebLogger interface {
	// WithWebLogs returns true if web logs are enabled.
	WithWebLogs() bool

	// WebLogsNsqdTCP returns an address to a NSQ messaging TCP service or
	// an empty string.
	WebLogsNsqdTCP() string
}
