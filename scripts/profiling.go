// The purpose of this script is to find out how fast algorithms can go through
// a list of 100_000 names.
package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	mlib "github.com/gnames/gnlib/domain/entity/matcher"
	"github.com/gnames/gnlib/encode"
)

const Batch = 10_000

const url = "http://:8080/"

func main() {
	var wgRes sync.WaitGroup
	chNames := make(chan []string)
	go namesToChannel(chNames)
	wgRes.Add(1)
	go processData(chNames, &wgRes)
	wgRes.Wait()
}

func processData(chNames <-chan []string, wg *sync.WaitGroup) {
	enc := encode.GNjson{}
	defer wg.Done()
	w := csv.NewWriter(os.Stdout)
	defer func() {
		w.Flush()
	}()
	count := 0
	timeStart := time.Now().UnixNano()
	for request := range chNames {
		count += 1
		total := count * Batch
		req, err := enc.Encode(&request)
		if err != nil {
			log.Fatalf("Cannot marshall input: %v", err)
		}
		resp, err := http.Post(url+"match", "application/json", bytes.NewReader(req))
		if err != nil {
			log.Fatalf("Cannot send request: %v", err)
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Cannot get data: %v", err)
		}
		var response []mlib.Match
		enc.Decode(respBytes, &response)

		var name, match, matchType string
		var editDist, editDistStem int
		for _, res := range response {
			name = res.Name
			matchType = res.MatchType.String()
			if len(res.MatchItems) == 0 {
				err = w.Write([]string{name, matchType, "", "", ""})
				if err != nil {
					log.Fatalf("Cannot write CSV file: %s", err)
				}
			}
			for _, v := range res.MatchItems {
				match = v.MatchStr
				editDist = int(v.EditDistance)
				editDistStem = int(v.EditDistanceStem)
				err = w.Write([]string{
					name, matchType, match,
					strconv.Itoa(editDist), strconv.Itoa(editDistStem)})
				if err != nil {
					log.Fatalf("Cannot write CSV file: %s", err)
				}
			}
		}

		if total%10_000 == 0 {
			timeSpent := float64(time.Now().UnixNano()-timeStart) / 1_000_000_000
			speed := int64(float64(total) / timeSpent)
			log.Printf("Verified %s names, %s names/sec",
				humanize.Comma(int64(total)), humanize.Comma(speed))
		}
	}
}

func namesToChannel(chNames chan<- []string) {
	path := filepath.Join("..", "testdata", "testdata.csv")
	names := make([]string, 0, Batch)
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open %s: %s", path, err)
	}
	defer f.Close()
	r := csv.NewReader(f)

	// skip header
	_, err = r.Read()
	if err != nil {
		log.Fatalf("Cannot read from %s: %s", path, err)
	}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read body from %s: %s", path, err)
		}
		names = append(names, row[0])
		if len(names) > Batch-1 {
			chNames <- names
			names = make([]string, 0, Batch)
		}
	}
	chNames <- names
	close(chNames)
}
