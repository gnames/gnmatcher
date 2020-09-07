/* Pacage model contains main structures and interfaces that describe input,
output and methods that provide functionality for gnmatcher service.
*/
package model

import (
	gn "github.com/gnames/gnames/model"
)

// Pong is output from Ping method. Suppose to return "pong".
type Pong string

// Version is output from GetVersion method.
type Version struct {
	// Version of gnmatcher.
	Version string
	// Build timestamp of gnmatcher.
	Build string
}

// Match is output of MatchAry method.
type Match struct {
	// ID is UUIDv5 generated from verbatim input name-string.
	ID string
	// Name is verbatim input name-string.
	Name string
	// VirusMatch is true if matching
	VirusMatch bool
	// MatchType describe what kind of match happened.
	MatchType gn.MatchType
	// MatchItems provide all matched data. It will be empty if no matches
	// occured.
	MatchItems []MatchItem
}

// MatchItem describes one matched string and its properties.
type MatchItem struct {
	// ID is a UUIDv5 generated out of MatchStr.
	ID string
	// MatchStr is the string that matched a particular input. More often than
	// not it is a canonical form of a name. However for viruses it
	// can be matched string from the database.
	MatchStr string
	// EditDistance is a Levenshtein edit distance between normalized
	// input and MatchStr.
	EditDistance int
	// EditDistanceStem is a Levenshtein edit distance between stemmed input and
	// stemmed MatchStr.
	EditDistanceStem int
}
