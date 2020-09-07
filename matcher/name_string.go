package matcher

import (
	"strings"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/gogna/gnparser"
	"gitlab.com/gogna/gnparser/pb"
	"gitlab.com/gogna/gnparser/stemmer"
)

// NameString stores input data for doing exact, fuzzy, exact partial, and
// fuzzy partial matching. It is created by parsing a name-string and
// storing its semantic elements.
type NameString struct {
	// ID is UUID v5 generated from the verbatim name-string.
	ID string
	// Name is a verbatim name-string.
	Name string
	// Cardinality is the apparent number of elemenents in a name. Uninomial
	// corresponds to cardinality 1, bionmial to 2, trinomial to 3 etc.
	Cardinality int
	// Canonical is the simplest most common version of a canonical form of
	// a name string.
	Canonical string
	// CanonicalID is UUID v5 generated from the Canonical field.
	CanonicalID string
	// Canonical Stem is version of the Canonical field with suffixes removed
	// and characters substituted according to rules of Latin grammar.
	CanonicalStem string
	// Partial contains truncated versions of Canonical form. It is important
	// for matching names that could not be matched for all specific epithets.
	Partial *Partial
}

// Partial stores truncated version of a 'canonical' name-string.
type Partial struct {
	// Genus is a truncated canonical form with all specific epithets removed.
	Genus string
	// Multinomials are truncated canonical forms where one or more specific
	// epithets removed.
	Multinomials []Multinomial
}

// Multinomial contains multinomial names that were constructed from
// an 'infraspecific' name-string.
type Multinomial struct {
	// Tail is genus + the last epithet.
	Tail string
	// Head is the name without the last epithet.
	Head string
}

// NewNameString creates a new instance of NameString.
func NewNameString(parser gnparser.GNparser,
	name string) (NameString, *pb.Parsed) {
	parsed := parser.ParseToObject(name)
	if parsed.Parsed {
		ns := NameString{
			ID:            parsed.Id,
			Name:          name,
			Cardinality:   int(parsed.Cardinality),
			Canonical:     parsed.Canonical.Simple,
			CanonicalID:   uuid.NewV5(GNUUID, parsed.Canonical.Simple).String(),
			CanonicalStem: parsed.Canonical.Stem,
		}

		ns.newPartial(parsed)
		// We do not fuzzy-match uninomials, however there are cases when
		// a binomial lost empty space during OCR. We increase probability to
		// match such binomials, if we stem them. It happens because we use trie
		// of stemmed canonicals.
		// For example we will be able to match 'Pardosamoestus' to 'Pardosa moesta'
		if parsed.Cardinality == 1 {
			ns.CanonicalStem = stemmer.Stem(ns.Canonical).Stem
		}
		return ns, parsed
	}

	return NameString{ID: parsed.Id, Name: name}, parsed
}

func (ns *NameString) newPartial(parsed *pb.Parsed) {
	if parsed.Cardinality < 2 {
		return
	}
	canAry := strings.Split(ns.Canonical, " ")

	ns.Partial = &Partial{Genus: canAry[0]}
	partialNum := parsed.Cardinality - 2

	// In case of binomial we return only genus
	if partialNum < 1 {
		return
	}

	ns.Partial.Multinomials = make([]Multinomial, partialNum)
	for i := range ns.Partial.Multinomials {
		lastLen := len(canAry) - i - 1
		tail := []string{ns.Partial.Genus, canAry[lastLen]}
		head := canAry[0:lastLen]

		ns.Partial.Multinomials[i] = Multinomial{
			Tail: strings.Join(tail, " "),
			Head: strings.Join(head, " "),
		}
	}
}

