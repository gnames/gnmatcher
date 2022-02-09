package virusio

import (
	"index/suffixarray"
	"os"
	"path/filepath"

	"github.com/gnames/gnfmt"
	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/io/dbase"
	"github.com/rs/zerolog/log"
)

func (v *virusio) prepareData() {
	log.Info().Msg("Preparing virus data")
	path := v.cfg.VirusDir()
	var err error

	if v.sufary != nil {
		return
	}

	err = v.dataFromCache(path)
	if err != nil {
		log.Info().Msgf("Cache for viruses at '%s' is empty", path)
		log.Info().Msg("Virus data will be received from the database")
	}

	if v.sufary != nil {
		return
	}

	var data []mlib.MatchItem
	data, err = v.dataFromDB(path)
	if err != nil {
		log.Fatal().Err(err).
			Msgf("Cannot create filters at %s from database", path)
	}
	bs := v.processData(data)
	err = v.saveData(bs)
	if err != nil {
		log.Fatal().Err(err).
			Msgf("Cannot save virus data to disk at '%s'", path)
	}
}

func (v *virusio) saveData(bs []byte) error {
	var err error
	path := v.cfg.VirusDir()
	err = os.WriteFile(filepath.Join(path, "viruses"), bs, 0664)
	if err != nil {
		return err
	}

	encoded, err := gnfmt.GNgob{}.Encode(v.mapMatchItems)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "uuids"), encoded, 0664)
	return err
}

func (v *virusio) processData(data []mlib.MatchItem) []byte {
	uuids := make(map[string]struct{})
	mapMatchItems := make(map[int]mlib.MatchItem)
	names := make([]byte, 0, len(data))
	var start int

	for i := range data {
		if _, ok := uuids[data[i].ID]; ok {
			continue
		}
		name := v.NameToBytes(data[i].MatchStr)
		names = append(names, name...)
		uuids[data[i].ID] = struct{}{}
		mapMatchItems[start] = data[i]
		start += len(name)
	}
	v.mapMatchItems = mapMatchItems
	v.sufary = suffixarray.New(names)
	return names
}

func (v *virusio) dataFromDB(path string) ([]mlib.MatchItem, error) {
	var res []mlib.MatchItem
	db := dbase.NewDB(v.cfg)
	log.Info().Msg("Importing lookup data for viruses")

	q := `SELECT name_string_id, name
  FROM verification v
    JOIN data_sources ds ON ds.id = v.data_source_id
  WHERE virus='true'
  ORDER by ds.is_curated desc`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	var uuid, name string
	for rows.Next() {
		if err = rows.Scan(&uuid, &name); err != nil {
			return nil, err
		}
		res = append(res, mlib.MatchItem{ID: uuid, MatchStr: name})
	}
	return res, err
}

func (v *virusio) dataFromCache(path string) error {
	var bs []byte
	var err error

	bs, err = os.ReadFile(filepath.Join(path, "viruses"))
	if err != nil {
		return err
	}
	v.sufary = suffixarray.New(bs)

	bs, err = os.ReadFile(filepath.Join(path, "uuids"))
	if err != nil {
		return err
	}

	var mapMatchItems map[int]mlib.MatchItem
	err = gnfmt.GNgob{}.Decode(bs, &mapMatchItems)
	if err != nil {
		return err
	}

	v.mapMatchItems = mapMatchItems
	return nil
}
