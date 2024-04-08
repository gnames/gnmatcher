// package virus contains an interface for matching strings to names of
// viruses, plasmids, prions and other non-cellular entities.
package virus

import (
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/pkg/config"
)

type VirusMatcher interface {
	// Init loads cached data into memory, or creates cache, if it does not
	// exist yet.
	Init()

	// SetConfig updates configuration of the matcher.
	SetConfig(cfg config.Config)

	// MatchVirus takes a virus name and returns back matched items for
	// the name. In case if there were too many returned results, returns an
	// error. Matching is successful if entered name matches the start of the
	// virus name string from the database. If there are too many matches,
	// the result is truncated. "Curated" databases have a priority in
	// returned results.
	MatchVirus(s string) []mlib.MatchItem

	// NameToBytes normalizes a virus name by removing all extra spaces,
	// converting all runes to lower case, adding '\x00' to the start and
	// returning result as bytes.
	NameToBytes(s string) []byte
}
