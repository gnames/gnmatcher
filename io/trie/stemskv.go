package trie

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"os"

	"github.com/dgraph-io/badger/v2"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnsys"
	log "github.com/sirupsen/logrus"
)

// initStemsKV creates key-value store for stems and their canonical forms.
func initStemsKV(path string, db *sql.DB) {
	var err error
	err = gnsys.MakeDir(path)
	if err != nil {
		log.Fatalf("Cannot create %s: %s", path, err)
	}

	if keyValExists(path) {
		log.Info("Stems key-value store already exists, skipping.")
		return
	}
	kv := connectKeyVal(path)
	defer kv.Close()

	q := `SELECT s.name as name_stem, c.name, c.id
          FROM canonical_stems s
            JOIN name_strings ns
              ON ns.canonical_stem_id = s.id
            JOIN canonicals c
              ON ns.canonical_id = c.id
        GROUP BY c.name, c.id, s.name
          ORDER BY name`

	rows, err := db.Query(q)
	if err != nil {
		log.Fatalf("Cannot get stems from DB: %s.", err)
	}

	kvTxn := kv.NewTransaction(true)
	var stemRes []mlib.MatchItem
	var currentStem, stem, name, id string
	count := 0
	for rows.Next() {
		if err := rows.Scan(&stem, &name, &id); err != nil {
			log.Fatalf("Cannot read stem data from query: %s.", err)
		}
		if currentStem == "" {
			currentStem = stem
		}
		if stem != currentStem {
			count += 1

			key := []byte(currentStem)
			var b bytes.Buffer
			enc := gob.NewEncoder(&b)
			if err = enc.Encode(stemRes); err != nil {
				log.Fatalf("Cannot marshal canonicals: %s.", err)
			}
			val := b.Bytes()
			if err = kvTxn.Set(key, val); err != nil {
				log.Fatalf("Transaction failed to set key: %s.", err)
			}
			if count > 10_000 {
				err = kvTxn.Commit()
				if err != nil {
					log.Fatalf("Transaction commit faied: %s.", err)
				}
				count = 0
				kvTxn = kv.NewTransaction(true)
			}
			currentStem = stem
			stemRes = nil
		}
		stemRes = append(stemRes, mlib.MatchItem{ID: id, MatchStr: name})
	}
	err = kvTxn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// connectKeyVal connects to a key-value store
func connectKeyVal(path string) *badger.DB {
	options := badger.DefaultOptions(path)
	// running in mem: options := badger.DefaultOptions("").WithInMemory(true)
	options.Logger = nil
	bdb, err := badger.Open(options)
	if err != nil {
		log.Fatalf("Cannot connect to key-value store: %s.", err)
	}
	return bdb
}

// getValue takes a string and a connection to a key-value store and checks if
// there is such stem key. It returns a list of canonicals that correspond to
// that key.
func getValue(kv *badger.DB, key string) []byte {
	var res []byte
	err := kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			log.Fatal(err)
		}

		return item.Value(func(val []byte) error {
			res = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}
	return res
}

// keyValExists checks if key-value store is set.
func keyValExists(path string) bool {
	files, err := os.ReadDir(path)
	return (err == nil && len(files) > 0)
}
