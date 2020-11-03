package rest

import mlib "github.com/gnames/gnlib/domain/entity/matcher"
import "github.com/gnames/gnlib/encode"

// MatcherService describes remote service of gnmatchter.
type MatcherService interface {
	// Port returns the port of the service.
	Port() int

	// Ping checks connection to the service.
	Ping() string

	mlib.Matcher

	encode.Encoder
}
