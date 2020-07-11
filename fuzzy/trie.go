package fuzzy

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dvirsky/levenshtein"
	"github.com/gnames/gnmatcher/dbase"
)

const trieFile = "stem.trie"

func GetTrie(triePath string, d dbase.Dbase) (*levenshtein.MinTree, error) {
	var trie *levenshtein.MinTree
	trie, err := getCachedTrie(triePath)
	if err == nil {
		log.Println("Trie data is rebuilt from cache")
		return trie, nil
	}

	db := d.NewDB()
	trie, err = populateTrie(db, triePath)
	if err != nil {
		return trie, err
	}
	return trie, nil
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

func populateTrie(db *sql.DB, triePath string) (*levenshtein.MinTree, error) {
	log.Println("Getting data from database")
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
	log.Println("Building trie and saving it to disk")
	path := filepath.Join(triePath, trieFile)
	w, err := os.Create(path)
	if err != nil {
		return trie, err
	}
	trie, err = levenshtein.NewMinTreeWrite(names, w)
	if err != nil {
		return trie, err
	}
	log.Println("Trie is created")
	return trie, nil
}
