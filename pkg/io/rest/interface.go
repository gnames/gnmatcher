// package rest provides http REST interface to gnmatcher functionality.
package rest

import (
	"github.com/gnames/gnfmt"
	gnmatcher "github.com/gnames/gnmatcher/pkg"
)

// MatcherService describes remote service of gnmatchter.
type MatcherService interface {
	// Port returns the port of the service.
	Port() int

	// Ping checks connection to the service.
	Ping() string

	// GNmatcher is the main use-case of the gnmatcher project.
	gnmatcher.GNmatcher

	// Encoder provides serialization/deserialization interface.
	gnfmt.Encoder
}
