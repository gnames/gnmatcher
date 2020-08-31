package rest

import (
	"github.com/gnames/gnmatcher"
	"github.com/gnames/gnmatcher/model"
)

// MatchMatcherREST implements MatcherService interface.
type MatcherREST struct {
	gnm  *gnmatcher.GNMatcher
	port int
}

// NewMNewMatcherREST is a constructor for MatchREST.
func NewMatcherREST(gnm *gnmatcher.GNMatcher, port int) MatcherREST {
	return MatcherREST{
		gnm:  gnm,
		port: port,
	}
}

// GetPort returns port number to the service.
func (mr MatcherREST) GetPort() int {
	return mr.port
}

// Ping returns "pong" message if connection to the service did succed.
func (mr MatcherREST) Ping() model.Pong {
	return model.Pong("pong")
}

// GetVersion returns version number and build timestamp of gnmatcher.
func (mr MatcherREST) GetVersion() model.Version {
	return model.Version{
		Version: gnmatcher.Version,
		Build:   gnmatcher.Build,
	}
}

// MatchAry takes a list of strings and matches them to known scientific names.
func (mr MatcherREST) MatchAry(names []string) []*model.Match {
	return mr.gnm.MatchNames(names)
}
