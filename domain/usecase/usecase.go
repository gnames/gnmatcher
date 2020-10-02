// Package usecase provide core usage of gnmatcher.
package usecase

import "github.com/gnames/gnmatcher/domain/entity"

// Matcher describes methods required for matching name-strings to names.
type Matcher interface {
	// Version sends current version and build timestamp of
	// gnmatcher.
	Version() entity.Version

	// MatchAry takes a list of strings and matches each of them
	// to known scientific names.
	MatchAry(names []string) []*entity.Match
}

// FuzzyMatcher describes methods needed for fuzzy matching
type FuzzyMatcher interface {
	// MatchStem takes a stemmed scientific name and max edit distance.
	// The search stops if current edit distance becomes bigger than edit
	// distance. The method returns 0 or more stems that did match the
	// input stem within the edit distance constraint.
	MatchStem(stem string, maxEditDistance int) []string
	// StemToCanonicals takes a stem and returns back canonicals
	// that correspond to that stem.
	StemToMatchItems(stem string) []entity.MatchItem
}
