// package rest provides http REST interface to gnmatcher functionality.
package rest

import (
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnmatcher"
)

// MatcherService describes remote service of gnmatchter.
type MatcherService interface {
	// Port returns the port of the service.
	Port() int

	// Ping checks connection to the service.
	Ping() string

	// GNMatcher is the main use-case of the gnmatcher project.
	gnmatcher.GNMatcher

	// Encoder provides serialization/deserialization interface.
	encode.Encoder
}
