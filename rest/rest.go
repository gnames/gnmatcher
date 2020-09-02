package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gnames/gnmatcher/binary"
	"github.com/gnames/gnmatcher/model"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Run creates and runs the HTTP API service of gnmatcher.
func Run(m model.MatcherService) {
	log.Printf("Starting the HTTP API server on port %d.", m.GetPort())
	r := mux.NewRouter()

	r.HandleFunc("/",
		func(resp http.ResponseWriter, req *http.Request) {
			rootHTTP(resp, req)
		}).Methods("GET")

	r.HandleFunc("/ping",
		func(resp http.ResponseWriter, req *http.Request) {
			pingHTTP(resp, req, m)
		}).Methods("POST")

	r.HandleFunc("/version",
		func(resp http.ResponseWriter, req *http.Request) {
			getVersionHTTP(resp, req, m)
		}).Methods("POST")

	r.HandleFunc("/match",
		func(resp http.ResponseWriter, req *http.Request) {
			matchAryHTTP(resp, req, m)
		}).Methods("POST")

	addr := fmt.Sprintf(":%d", m.GetPort())

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func rootHTTP(resp http.ResponseWriter, _ *http.Request) {
	log.Print("Pong from root")
	resp.Write([]byte("OK"))
}

func pingHTTP(resp http.ResponseWriter, _ *http.Request,
	m model.MatcherService) {
	result := m.Ping()
	if response, err := binary.Encode(result); err == nil {
		resp.Write(response)
	} else {
		log.Warnf("pingHTTP: cannot encode response : %v", err)
	}
}

func getVersionHTTP(resp http.ResponseWriter, _ *http.Request,
	m model.MatcherService) {
	result := m.GetVersion()
	if out, err := binary.Encode(result); err == nil {
		resp.Write(out)
	} else {
		log.Warnf("getVersionHTTP: cannot encode response : %v", err)
	}
}

func matchAryHTTP(resp http.ResponseWriter, req *http.Request,
	m model.MatcherService) {
	var names []string
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warnf("matchAryHTTP: cannot read message from request : %v", err)
		return
	}
	binary.Decode(body, &names)

	matches := m.MatchAry(names)

	if out, err := binary.Encode(matches); err == nil {
		resp.Write(out)
	} else {
		log.Warnf("MatchAry: Cannot encode response : %v", err)
	}
}
