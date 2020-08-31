package model

// MatcherService describes remote service of gnmatchter.
type MatcherService interface {
	// GetPort returns port of the service.
	GetPort() int

	// Ping checks connection to the service.
	Ping() Pong

	// GetVersion sends current version and build timestamp of
	// gnmatcher.
	GetVersion() Version

	// MatchAry takes a list of strings and matches each of them
	// to known scientific names.
	MatchAry([]string) []*Match
}
