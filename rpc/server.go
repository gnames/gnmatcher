package rpc

import (
	"context"
	"fmt"
	"net"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnmatcher"
	"github.com/gnames/gnmatcher/protob"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// gnmatcherServer implements gRPC server out of protob package.
type gnmatcherServer struct {
	matcher *gnmatcher.GNMatcher
}

// Run starts the gRPC service.
func Run(port int, gnm *gnmatcher.GNMatcher) {
	defer gnm.KV.Close()
	log.Info(fmt.Sprintf("Starting gnmatcher gRPC server on port %d.", port))
	gnms := gnmatcherServer{
		matcher: gnm,
	}
	srv := grpc.NewServer()
	protob.RegisterGNMatcherServer(srv, gnms)
	portVal := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portVal)
	if err != nil {
		log.Fatalf("Could not listen on port %s: %s.", portVal, err)
	}
	log.Printf("Matching runs on %d parallel jobs", gnm.Matcher.Config.JobsNum)
	log.Fatal(srv.Serve(l))
}

// Ping checks if connection to gRPC service is working.
func (gnmatcherServer) Ping(ctx context.Context,
	void *protob.Void) (*protob.Pong, error) {
	pong := protob.Pong{Value: "pong"}
	return &pong, nil
}

// Ver sends back the version of the _gnmatcher_.
func (gnmatcherServer) Ver(ctx context.Context,
	void *protob.Void) (*protob.Version, error) {
	ver := protob.Version{Version: gnfinder.Version}
	return &ver, nil
}

// MatchAry is the main function of the service. It takes a list of
// name-strings and returns back matched canonical forms for each name-string.
func (gnms gnmatcherServer) MatchAry(ctx context.Context,
	names *protob.Names) (*protob.Output, error) {
	res := gnms.matcher.MatchNames(names.Names)
	output := &protob.Output{
		Results: res,
	}
	return output, nil
}
