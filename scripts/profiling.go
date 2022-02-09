// The purpose of this script is to find out how fast algorithms can go through
// a list of 100_000 names.
package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gnames/gnfmt"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/rs/zerolog/log"
)

const batch = 10_000

const url = "http://:8080/api/v1/"

func main() {
	var wgRes sync.WaitGroup
	chNames := make(chan []string)
	go namesToChannel(chNames)
	wgRes.Add(1)
	go processData(chNames, &wgRes)
	wgRes.Wait()
}

func processData(chNames <-chan []string, wg *sync.WaitGroup) {
	enc := gnfmt.GNjson{}
	defer wg.Done()
	w := csv.NewWriter(os.Stdout)
	defer func() {
		w.Flush()
	}()
	count := 0
	timeStart := time.Now().UnixNano()
	for request := range chNames {
		count++
		total := count * batch
		req, err := enc.Encode(&request)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot marshall input")
		}
		resp, err := http.Post(url+"matches", "application/json", bytes.NewReader(req))
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot send request")
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot get data")
		}
		var response []mlib.Match
		_ = enc.Decode(respBytes, &response)

		var name, match, matchType string
		var editDist, editDistStem int
		for _, res := range response {
			name = res.Name
			matchType = res.MatchType.String()
			if len(res.MatchItems) == 0 {
				err = w.Write([]string{name, matchType, "", "", ""})
				if err != nil {
					log.Fatal().Err(err).Msg("Cannot write CSV file")
				}
			}
			for _, v := range res.MatchItems {
				match = v.MatchStr
				editDist = v.EditDistance
				editDistStem = v.EditDistanceStem
				err = w.Write([]string{
					name, matchType, match,
					strconv.Itoa(editDist), strconv.Itoa(editDistStem)})
				if err != nil {
					log.Fatal().Err(err).Msg("Cannot write CSV file")
				}
			}
		}

		if total%10_000 == 0 {
			timeSpent := float64(time.Now().UnixNano()-timeStart) / 1_000_000_000
			speed := int64(float64(total) / timeSpent)
			log.Info().
				Str("names", humanize.Comma(int64(total))).
				Str("names/sec", humanize.Comma(speed)).
				Msg("Verified")
		}
	}
}

func namesToChannel(chNames chan<- []string) {
	path := filepath.Join("..", "testdata", "testdata.csv")
	names := make([]string, 0, batch)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot open %s", path)
	}
	defer f.Close()
	r := csv.NewReader(f)

	// skip header
	_, err = r.Read()
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot read from %s", path)
	}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot read body from %s", path)
		}
		names = append(names, row[0])
		if len(names) > batch-1 {
			chNames <- names
			names = make([]string, 0, batch)
		}
	}
	chNames <- names
	close(chNames)
}
