package trie

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log/slog"
	"os"

	"github.com/dgraph-io/badger/v2"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnsys"
)

// initStemsKV creates key-value store for stems and their canonical forms.
func initStemsKV(path string, db *sql.DB) error {
	var err error
	err = gnsys.MakeDir(path)
	if err != nil {
		slog.Error("Cannot create dir", "path", path, "error", err)
		return err
	}

	if keyValExists(path) {
		slog.Info("Stems key-value store already exists, skipping")
		return nil
	}
	kv, err := connectKeyVal(path)
	if err != nil {
		return err
	}
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
		slog.Error("Cannot get stems from DB", "error", err)
		return err
	}
	slog.Info("Setting Stems Key-Value store")
	kvTxn := kv.NewTransaction(true)
	var stemRes []mlib.MatchItem
	var dsMap map[int]struct{}
	var dsID int
	var currentStem, currentID, currentName, stem, name, id string
	count := 0
	for rows.Next() {
		if err = rows.Scan(&stem, &name, &id, &dsID); err != nil {
			slog.Error("Cannot read stem data from query", "error", err)
			return err
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
					ID:             currentID,
					MatchStr:       currentName,
					DataSourcesMap: dsMap,
				})
			setKeyVal(kvTxn, currentStem, stemRes)
			if count > 10_000 {
				err = kvTxn.Commit()
				if err != nil {
					slog.Error("Transaction commit faied", "error", err)
					return err
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
					ID:             currentID,
					MatchStr:       currentName,
					DataSourcesMap: dsMap,
				})
			currentID = id
			currentName = name
			dsMap = make(map[int]struct{})
		}
		dsMap[dsID] = struct{}{}

	}
	stemRes = append(stemRes,
		mlib.MatchItem{
			ID:             id,
			MatchStr:       name,
			DataSourcesMap: dsMap,
		})
	setKeyVal(kvTxn, currentStem, stemRes)
	err = kvTxn.Commit()
	if err != nil {
		slog.Error("Cannot commit kay-value transaction", "error", err)
		return err
	}
	return nil
}

func setKeyVal(kvTxn *badger.Txn,
	stem string,
	stemRes []mlib.MatchItem,
) error {
	var err error
	key := []byte(stem)
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err = enc.Encode(stemRes); err != nil {
		slog.Error("Cannot marshal canonicals", "error", err)
		return err
	}
	val := b.Bytes()
	if err = kvTxn.Set(key, val); err != nil {
		slog.Error("Transaction failed to set key", "error", err)
		return err
	}
	return nil
}

// connectKeyVal connects to a key-value store
func connectKeyVal(path string) (*badger.DB, error) {
	options := badger.DefaultOptions(path)
	// running in mem: options := badger.DefaultOptions("").WithInMemory(true)
	options.Logger = nil
	bdb, err := badger.Open(options)
	if err != nil {
		slog.Error("Cannot connect to key-value store", "error", err)
		return nil, err
	}
	return bdb, nil
}

// getValue takes a string and a connection to a key-value store and checks if
// there is such stem key. It returns a list of canonicals that correspond to
// that key.
func getValue(kv *badger.DB, key string) ([]byte, error) {
	var res []byte
	err := kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			slog.Error(
				"Cannot retrieve value from key-value store",
				"key", key, "error", err,
			)
			return err
		}

		return item.Value(func(val []byte) error {
			res = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		slog.Error("Cannot get value from key-value store", "error", err)
		return nil, err
	}
	return res, nil
}

// keyValExists checks if key-value store is set.
func keyValExists(path string) bool {
	files, err := os.ReadDir(path)
	return (err == nil && len(files) > 0)
}
