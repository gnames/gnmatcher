// package exact contains interface for exact-matching strings to
// known scientific names.
package exact

import "github.com/gnames/gnmatcher/pkg/config"

// ExactMatcher is the interface for exact matching strings.
// It matches them using UUIDv5 strings generated from the strings.
type ExactMatcher interface {
	// Init loads cached data into memory, and creates cache, if it does not
	// exist yet.
	Init()

	// SetConfig updates configuration of the matcher.
	SetConfig(config config.Config)

	// MatchCanonicalID matches canonical forms of scientific names. It takes
	// UUIDv5 filter generated out of name-string and checks if the same
	// UUIDv5 exists in the cached data.
	MatchCanonicalID(uuid string) bool
}
