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

	r.HandleFunc("/", rootHTTP()).Methods("GET")
	r.HandleFunc("/ping", pingHTTP(m)).Methods("GET")
	r.HandleFunc("/version", versionHTTP(m)).Methods("GET")
	r.HandleFunc("/match", matchNamesHTTP(m)).Methods("POST")

	addr := fmt.Sprintf(":%d", m.Port())

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func rootHTTP() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Pong from root")
		w.Write([]byte("OK"))
	}
}

func pingHTTP(m MatcherService) func(http.ResponseWriter, *http.Request) {
	result := m.Ping()

	return func(w http.ResponseWriter, req *http.Request) {
		if response, err := m.Encode(result); err == nil {
			w.Write(response)
		} else {
			log.Warnf("pingHTTP: cannot encode response : %v", err)
		}
	}
}

func versionHTTP(m MatcherService) func(http.ResponseWriter, *http.Request) {
	result := m.GetVersion()

	return func(w http.ResponseWriter, req *http.Request) {
		if response, err := m.Encode(result); err == nil {
			w.Write(response)
		} else {
			log.Warnf("versionHTTP: cannot encode response : %v", err)
		}
	}
}

func matchNamesHTTP(m MatcherService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
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

		matches := m.MatchNames(names)

		if response, err := m.Encode(matches); err == nil {
			w.Write(response)
		} else {
			log.Warnf("versionHTTP: cannot encode response : %v", err)
		}
	}
}
