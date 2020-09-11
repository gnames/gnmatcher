package rest

import "github.com/gnames/gnmatcher/domain/usecase"
import "github.com/gnames/gnames/lib/encode"

// MatcherService describes remote service of gnmatchter.
type MatcherService interface {
	// Port returns the port of the service.
	Port() int

	// Ping checks connection to the service.
	Ping() string

	usecase.Matcher

	encode.Encoder
}
