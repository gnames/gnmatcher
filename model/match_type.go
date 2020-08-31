package model

// MatchType descrbes categories that might happen during the matching
// process.
type MatchType int

const (
	// None means no match was found.
	None MatchType = iota
	// Canonical means the match occured to a canonical form of a name-string.
	Canonical
	// CanonicalFull means the match happened to a exended canonical form with
	// ranks.
	CanonicalFull
	// Virus means there was exact match to a verbatim virus name.
	Virus
	// Fuzzy means that the match happened to canonical form, but it is not exact.
	Fuzzy
	// Partial means the complete name-string did not match, but an exact match
	// happened to its truncated form without some epithets.
	Partial
	// PartialFuzzy means that fuzzy match was found to a truncated name.
	PartialFuzzy
)

var matchTypeStrings = map[int]string{
	0: "NONE",
	1: "CANONICAL",
	2: "CANONICAL_FULL",
	3: "VIRUS",
	4: "FUZZY",
	5: "PARTIAL",
	6: "PARTIAL_FUZZY",
}

// String provides a string form for a MatchType.
func (mt MatchType) String() string {
	return matchTypeStrings[int(mt)]
}

