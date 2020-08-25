// The purpose of this script is to find out how fast algorithms can go through
// a list of 100_000 names.
package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gnames/gnmatcher/protob"
	"google.golang.org/grpc"
)

const Batch = 1000

const HostRPC = ":8778"

var (
	conn *grpc.ClientConn
)

func main() {
	var wgRes sync.WaitGroup
	chNames := make(chan protob.Names)
	go namesToChannel(chNames)
	wgRes.Add(1)
	go grpcData(chNames, &wgRes)
	wgRes.Wait()
}

func grpcData(chNames <-chan protob.Names, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	path := filepath.Join("..", "testdata", "profiling-res.csv")
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Cannot create %s: %s", path, err)
	}
	w := csv.NewWriter(file)
	defer func() {
		w.Flush()
		file.Sync()
		file.Close()
	}()
	conn, err = grpc.Dial(HostRPC, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Cannot connect to gRPC: %s", err)
	}
	client := protob.NewGNMatcherClient(conn)
	count := 0
	timeStart := time.Now().UnixNano()
	for names := range chNames {
		count += 1
		total := count * Batch
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		out, err := client.MatchAry(ctx, &names)
		if err != nil {
			log.Fatalf("Cannot match data: %s", err)
		}
		var name, match, matchType string
		var editDist, editDistStem int
		for _, res := range out.Results {
			name = res.Name
			matchType = res.MatchType.String()
			if len(res.MatchData) == 0 {
				err = w.Write([]string{name, matchType, "", "", ""})
				if err != nil {
					log.Fatalf("Cannot write to CSV file %s: %s", path, err)
				}
			}
			for _, v := range res.MatchData {
				match = v.MatchStr
				editDist = int(v.EditDistance)
				editDistStem = int(v.EditDistanceStem)
				err = w.Write([]string{
					name, matchType, match,
					strconv.Itoa(editDist), strconv.Itoa(editDistStem)})
				if err != nil {
					log.Fatalf("Cannot write to CSV file %s: %s", path, err)
				}
			}
		}

		if total%10_000 == 0 {
			timeSpent := float64(time.Now().UnixNano()-timeStart) / 1_000_000_000
			speed := int64(float64(total) / timeSpent)
			log.Printf("Verified %s names, %s names/sec",
				humanize.Comma(int64(total)), humanize.Comma(speed))
		}
		cancel()
	}
}

func namesToChannel(chNames chan<- protob.Names) {
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
			chNames <- protob.Names{Names: names}
			names = make([]string, 0, Batch)
		}
	}
	chNames <- protob.Names{Names: names}
	close(chNames)
}
