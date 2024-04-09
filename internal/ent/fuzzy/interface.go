// package fuzzy contains interfaces and code to facilitate fuzzy-matching
// of name-strings to scientific names collected in gnames database.
package fuzzy

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/pkg/config"
)

// FuzzyMatcher describes methods needed for fuzzy matching.
type FuzzyMatcher interface {
	// Initialize data for the matcher.
	Init() error

	// SetConfig updates configuration of the matcher.
	SetConfig(cfg config.Config)

	// MatchStem takes a stemmed scientific name and max edit distance.
	// The search stops if current edit distance becomes bigger than edit
	// distance. The method returns 0 or more stems that did match the
	// input stem within the edit distance constraint.
	MatchStem(stem string) []string

	// MatchStemExact takes a stem and returns true if the exact match of
	// the stem is found.
	MatchStemExact(stem string) bool

	// StemToCanonicals takes a stem and returns back canonicals
	// that correspond to that stem.
	StemToMatchItems(stem string) ([]mlib.MatchItem, error)
}
