package matcher

import (
	"strings"

	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/parsed"
	"github.com/gnames/gnuuid"
)

// nameString stores input data for doing exact, fuzzy, exact partial, and
// fuzzy partial matching. It is created by parsing a name-string and
// storing its semantic elements.
type nameString struct {
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
	// IsVirus is true if parsed name seem to be a virus
	IsVirus bool
	// Partial contains truncated versions of Canonical form. It is important
	// for matching names that could not be matched for all specific epithets.
	Partial *partial
}

// partial stores truncated version of a 'canonical' name-string.
type partial struct {
	// Genus is a truncated canonical form with all specific epithets removed.
	Genus string
	// Multinomials are truncated canonical forms where one or more specific
	// epithets removed.
	Multinomials []multinomial
}

// multinomial contains multinomial names that were constructed from
// an 'infraspecific' name-string.
type multinomial struct {
	// Tail is genus + the last epithet.
	Tail string
	// Head is the name without the last epithet.
	Head string
}

// newNameString creates a new instance of NameString.
func newNameString(
	parser gnparser.GNparser,
	name string,
) (nameString, *parsed.Parsed) {
	prsd := parser.ParseName(name)
	if prsd.Parsed {
		ns := nameString{
			ID:            prsd.VerbatimID,
			Name:          name,
			Cardinality:   int(prsd.Cardinality),
			Canonical:     prsd.Canonical.Simple,
			CanonicalID:   gnuuid.New(prsd.Canonical.Simple).String(),
			CanonicalStem: prsd.Canonical.Stemmed,
		}

		ns.newPartial(prsd)
		return ns, &prsd
	}

	return nameString{
		ID:      prsd.VerbatimID,
		Name:    name,
		IsVirus: prsd.Virus,
	}, &prsd
}

func (ns *nameString) newPartial(prsd parsed.Parsed) {
	if prsd.Cardinality < 2 {
		return
	}
	canAry := strings.Split(ns.Canonical, " ")

	ns.Partial = &partial{Genus: canAry[0]}
	partialNum := prsd.Cardinality - 2

	// In case of binomial we return only genus
	if partialNum < 1 {
		return
	}

	ns.Partial.Multinomials = make([]multinomial, partialNum)
	for i := range ns.Partial.Multinomials {
		lastLen := len(canAry) - i - 1
		tail := []string{ns.Partial.Genus, canAry[lastLen]}
		head := canAry[0:lastLen]

		ns.Partial.Multinomials[i] = multinomial{
			Tail: strings.Join(tail, " "),
			Head: strings.Join(head, " "),
		}
	}
}
