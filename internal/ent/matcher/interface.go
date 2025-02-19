// package matcher is the central processing unit for matching name-strings
// to known scientific names.
package matcher

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/pkg/config"
)

// Matcher is the interface that enables matching strings to known scientific
// names.
type Matcher interface {
	// Init loads data from cache on disk, and, if cache is empty, populates it
	// from gnames database.
	Init() error
	// MatchNames takes a slice of strings and returns back metadata
	// of the request and the matches of the strings to known scientific names.
	MatchNames(names []string, opt ...config.Option) mlib.Output
}
