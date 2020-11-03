package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Run creates and runs the HTTP API service of gnmatcher.
func Run(m MatcherService) {
	log.Printf("Starting the HTTP API server on port %d.", m.Port())
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

	addr := fmt.Sprintf(":%d", m.Port())

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func rootHTTP(resp http.ResponseWriter, _ *http.Request) {
	log.Debug("Pong from root")
	resp.Write([]byte("OK"))
}

func pingHTTP(resp http.ResponseWriter, _ *http.Request,
	m MatcherService) {
	result := m.Ping()
	if response, err := m.Encode(result); err == nil {
		resp.Write(response)
	} else {
		log.Warnf("pingHTTP: cannot encode response : %v", err)
	}
}

func getVersionHTTP(resp http.ResponseWriter, _ *http.Request,
	m MatcherService) {
	result := m.GetVersion()
	if out, err := m.Encode(result); err == nil {
		resp.Write(out)
	} else {
		log.Warnf("getVersionHTTP: cannot encode response : %v", err)
	}
}

func matchAryHTTP(resp http.ResponseWriter, req *http.Request,
	m MatcherService) {
	var names []string
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Warnf("matchAryHTTP: cannot read message from request : %v", err)
		return
	}
	err = m.Decode(body, &names)
	if err != nil {
		log.Warnf("matchAryHTTP: cannot decode request : %v", err)
		return
	}

	matches := m.MatchAry(names)

	if out, err := m.Encode(matches); err == nil {
		resp.Write(out)
	} else {
		log.Warnf("MatchAry: Cannot encode response : %v", err)
	}
}
