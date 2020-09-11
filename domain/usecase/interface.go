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
