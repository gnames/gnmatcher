package fuzzy

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dvirsky/levenshtein"
	log "github.com/sirupsen/logrus"
)

const trieFile = "stem.trie"

// GetTrie generates an in-memory trie for levenshtein automata. Such tree
// can either be constructed from database or from a dump file. The tree
// consists stemmed canonical forms of _gnames_ database.
func GetTrie(triePath string, db *sql.DB) *levenshtein.MinTree {
	var trie *levenshtein.MinTree
	trie, err := getCachedTrie(triePath)
	if err == nil {
		log.Println("Trie data is rebuilt from cache.")
		return trie
	}

	trie, err = populateAndSaveTrie(db, triePath)
	if err != nil {
		log.Fatalf("Cannot build trie from db: %s", err)
	}
	return trie
}

func getTrieSize(db *sql.DB) (int, error) {
	q := fmt.Sprint("SELECT count(*) from canonical_stems")
	var num int
	row := db.QueryRow(q)
	if err := row.Scan(&num); err != nil {
		return 0, err
	}
	return num, nil
}

func getCachedTrie(triePath string) (*levenshtein.MinTree, error) {
	var trie *levenshtein.MinTree
	path := filepath.Join(triePath, trieFile)
	trieFile, err := os.Open(path)
	if err != nil {
		return trie, err
	}
	return levenshtein.LoadMinTree(trieFile)
}

func populateAndSaveTrie(db *sql.DB, triePath string) (*levenshtein.MinTree, error) {
	log.Println("Getting trie data from database.")
	var trie *levenshtein.MinTree
	size, err := getTrieSize(db)
	if err != nil {
		return trie, err
	}
	names := make([]string, size)

	var name string
	q := fmt.Sprint("SELECT name FROM canonical_stems order by name")
	rows, err := db.Query(q)
	if err != nil {
		return trie, err
	}

	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return trie, err
		}
		names = append(names, name)
	}
	log.Println("Building trie and saving it to disk.")
	path := filepath.Join(triePath, trieFile)
	w, err := os.Create(path)
	if err != nil {
		return trie, err
	}
	trie, err = levenshtein.NewMinTreeWrite(names, w)
	if err != nil {
		return trie, err
	}
	log.Println("Trie is created.")
	return trie, nil
}
