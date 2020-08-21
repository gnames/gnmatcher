// stems_db package operates on a key-value store that contains stems and
// canonical forms that correspond to these stems. It create such a key-value
// store if necessary.
package stemskv

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"

	"github.com/dgraph-io/badger/v2"
	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/sys"
	log "github.com/sirupsen/logrus"
)

type CanonicalKV struct {
	ID   string
	Name string
}

// NewStemsKV creates key-value store for stems and their canonical forms.
func NewStemsKV(path string, d dbase.Dbase) {
	var err error
	db := d.NewDB()
	err = sys.MakeDir(path)
	if err != nil {
		log.Fatalf("Cannot create %s: %s", path, err)
	}

	if keyValExists(path) {
		log.Info("Stems key-value store already exists, skipping.")
		return
	}
	kv := ConnectKeyVal(path)
	defer kv.Close()

	q := `SELECT name_stem, name, id
          FROM canonicals
          WHERE name_stem IS NOT NULL
            AND name_stem !~ '\.'
          ORDER BY name_stem`
	rows, err := db.Query(q)
	if err != nil {
		log.Fatalf("Cannot get stems from DB: %s.", err)
	}

	kvTxn := kv.NewTransaction(true)
	var stemRes []CanonicalKV
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
			if count > 100_000 {
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
		stemRes = append(stemRes, CanonicalKV{ID: id, Name: name})
	}
	err = kvTxn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// ConnectKeyVal connects to a key-value store
func ConnectKeyVal(path string) *badger.DB {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	bdb, err := badger.Open(options)
	if err != nil {
		log.Fatal("Cannot connect to key-value store: %s.", err)
	}
	return bdb
}

// GetValue takes a string and a connection to a key-value store and checks if
// there is such stem key. It returns a list of canonicals that correspond to
// that key.
func GetValue(kv *badger.DB, key string) []byte {
	txn := kv.NewTransaction(false)
	defer func() {
		err := txn.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}()

	val, err := txn.Get([]byte(key))
	if err == badger.ErrKeyNotFound {
		return []byte("")
	} else if err != nil {
		log.Fatal(err)
	}
	var res []byte
	res, err = val.ValueCopy(res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func keyValExists(path string) bool {
	files, err := ioutil.ReadDir(path)
	return (err == nil && len(files) > 0)
}
