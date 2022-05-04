package trie

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"os"

	"github.com/dgraph-io/badger/v2"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
)

// initStemsKV creates key-value store for stems and their canonical forms.
func initStemsKV(path string, db *sql.DB) {
	var err error
	err = gnsys.MakeDir(path)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create %s", path)
	}

	if keyValExists(path) {
		log.Info().Msg("Stems key-value store already exists, skipping")
		return
	}
	kv := connectKeyVal(path)
	defer kv.Close()

	q := `SELECT s.name as name_stem, c.name, c.id, nsi.data_source_id
          FROM canonical_stems s
            JOIN name_strings ns
              ON ns.canonical_stem_id = s.id
            JOIN canonicals c
              ON ns.canonical_id = c.id
            JOIN name_string_indices nsi
              ON ns.id = nsi.name_string_id
        GROUP BY c.name, c.id, s.name, nsi.data_source_id
          ORDER BY name_stem`

	rows, err := db.Query(q)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get stems from DB")
	}
	log.Info().Msg("Setting Stems Key-Value store")
	kvTxn := kv.NewTransaction(true)
	var stemRes []mlib.MatchItem
	var dsMap map[int]struct{}
	var dsID int
	var currentStem, currentID, currentName, stem, name, id string
	count := 0
	for rows.Next() {
		if err = rows.Scan(&stem, &name, &id, &dsID); err != nil {
			log.Fatal().Err(err).Msg("Cannot read stem data from query")
		}
		if currentStem == "" {
			currentStem = stem
		}
		if currentID == "" {
			currentID = id
			currentName = name
			dsMap = make(map[int]struct{})
		}

		if stem != currentStem {
			count += 1

			stemRes = append(stemRes,
				mlib.MatchItem{
					ID:          currentID,
					MatchStr:    currentName,
					DataSources: dsMap,
				})
			setKeyVal(kvTxn, currentStem, stemRes)
			if count > 10_000 {
				err = kvTxn.Commit()
				if err != nil {
					log.Fatal().Err(err).Msg("Transaction commit faied")
				}
				count = 0
				kvTxn = kv.NewTransaction(true)
			}
			currentStem = stem
			currentID = id
			currentName = name
			stemRes = nil
			dsMap = make(map[int]struct{})
		}

		if id != currentID {
			stemRes = append(stemRes,
				mlib.MatchItem{
					ID:          currentID,
					MatchStr:    currentName,
					DataSources: dsMap,
				})
			currentID = id
			currentName = name
			dsMap = make(map[int]struct{})
		}
		dsMap[dsID] = struct{}{}

	}
	stemRes = append(stemRes,
		mlib.MatchItem{
			ID:          currentID,
			MatchStr:    currentName,
			DataSources: dsMap,
		})
	setKeyVal(kvTxn, currentStem, stemRes)
	err = kvTxn.Commit()
	if err != nil {
		log.Fatal().Err(err)
	}
}

func setKeyVal(kvTxn *badger.Txn,
	stem string,
	stemRes []mlib.MatchItem,
) {
	var err error
	key := []byte(stem)
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err = enc.Encode(stemRes); err != nil {
		log.Fatal().Err(err).Msg("Cannot marshal canonicals")
	}
	val := b.Bytes()
	if err = kvTxn.Set(key, val); err != nil {
		log.Fatal().Err(err).Msg("Transaction failed to set key")
	}
}

// connectKeyVal connects to a key-value store
func connectKeyVal(path string) *badger.DB {
	options := badger.DefaultOptions(path)
	// running in mem: options := badger.DefaultOptions("").WithInMemory(true)
	options.Logger = nil
	bdb, err := badger.Open(options)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to key-value store")
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
			log.Fatal().Err(err)
		}

		return item.Value(func(val []byte) error {
			res = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return res
}

// keyValExists checks if key-value store is set.
func keyValExists(path string) bool {
	files, err := os.ReadDir(path)
	return (err == nil && len(files) > 0)
}
