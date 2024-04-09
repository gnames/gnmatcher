// The purpose of this script is to find out how fast algorithms can go through
// a list of 100_000 names.
package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gnames/gnfmt"
	mlib "github.com/gnames/gnlib/ent/matcher"
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
	for names := range chNames {
		count++
		total := count * batch
		request := mlib.Input{Names: names}
		req, err := enc.Encode(&request)
		if err != nil {
			slog.Error("Cannot marshall input", "error", err)
			os.Exit(1)
		}
		resp, err := http.Post(url+"matches", "application/json", bytes.NewReader(req))
		if err != nil {
			slog.Error("Cannot send request", "error", err)
			os.Exit(1)
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Cannot get data", "error", err)
			os.Exit(1)
		}
		var response mlib.Output
		_ = enc.Decode(respBytes, &response)

		var name, match, matchType string
		var editDist, editDistStem int
		for _, res := range response.Matches {
			name = res.Name
			matchType = res.MatchType.String()
			if len(res.MatchItems) == 0 {
				err = w.Write([]string{name, matchType, "", "", ""})
				if err != nil {
					slog.Error("Cannot write CSV file", "error", err)
					os.Exit(1)
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
					slog.Error("Cannot write CSV file", "error", err)
					os.Exit(1)
				}
			}
		}

		if total%10_000 == 0 {
			timeSpent := float64(time.Now().UnixNano()-timeStart) / 1_000_000_000
			speed := int64(float64(total) / timeSpent)
			slog.Info("Verified",
				"names", humanize.Comma(int64(total)),
				"names/sec", humanize.Comma(speed),
			)
		}
	}
}

func namesToChannel(chNames chan<- []string) {
	path := filepath.Join("..", "testdata", "testdata.csv")
	names := make([]string, 0, batch)
	f, err := os.Open(path)
	if err != nil {
		slog.Error("Cannot open file", "path", path)
		os.Exit(1)
	}
	defer f.Close()
	r := csv.NewReader(f)

	// skip header
	_, err = r.Read()
	if err != nil {
		slog.Error("Cannot read from file", "path", path)
		os.Exit(1)
	}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Cannot read body from file", "path", path)
			os.Exit(1)
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
