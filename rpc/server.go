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

type gnmatcherServer struct {
	matcher *gnmatcher.GNMatcher
}

func Run(port int, gnm *gnmatcher.GNMatcher) {
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
	log.Fatal(srv.Serve(l))
}

func (gnmatcherServer) Ping(ctx context.Context,
	void *protob.Void) (*protob.Pong, error) {
	pong := protob.Pong{Value: "pong"}
	return &pong, nil
}

func (gnmatcherServer) Ver(ctx context.Context,
	void *protob.Void) (*protob.Version, error) {
	ver := protob.Version{Version: gnfinder.Version}
	return &ver, nil
}

func (gnms gnmatcherServer) MatchAry(ctx context.Context,
	names *protob.Names) (*protob.Output, error) {
	res := gnms.matcher.MatchNames(names.Names)
	output := &protob.Output{
		Results: res,
	}

	return output, nil
}

func (gnmatcherServer) MatchStream(stream protob.GNMatcher_MatchStreamServer) error {
	return nil
}
