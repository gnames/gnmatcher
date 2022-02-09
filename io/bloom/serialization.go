package bloom

import (
	"fmt"
	"os"
	"path/filepath"

	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	"github.com/rs/zerolog/log"
)

func saveFilters(path string, filters *bloomFilters) {
	var err error
	var nilFilter *baseBloomfilter.Bloomfilter
	files := map[string]*baseBloomfilter.Bloomfilter{
		canonicalFile: filters.canonical,
		sizesFile:     nilFilter,
	}

	for f, filter := range files {
		var file *os.File
		filePath := filepath.Join(path, f)

		file, err = createFile(filePath)
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot create %s", filePath)
		}
		if f == sizesFile {
			err = saveSizesFile(file, filters)
			if err != nil {
				log.Fatal().Err(err).Msg("Cannot create sizesFile")
			}
			continue
		}

		err = saveFilterFile(filePath, file, filter)
		if err != nil {
			log.Fatal().Err(err).Msgf("Cannot create %s", filePath)
		}
	}

	if err == nil {
		log.Info().Msg("Saved cached filters to disk")
	}
}

func createFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Warn().Err(err).Msgf("Could not create file %s", filePath)
	}
	return file, err
}

func saveFilterFile(
	filePath string,
	file *os.File,
	filter *baseBloomfilter.Bloomfilter,
) error {
	var bin []byte
	var err error
	bin, err = filter.MarshalBinary()
	if err != nil {
		log.Warn().Err(err).Msgf("Could not serialize for %s", filePath)
	}
	_, err = file.Write(bin)
	if err != nil {
		log.Warn().Err(err).Msgf("Could not save %s", filePath)
	}
	return err
}

func saveSizesFile(file *os.File, filters *bloomFilters) error {
	var err error
	sizes := fmt.Sprintf("CanonicalSize,%d\n",
		filters.canonicalSize)
	if _, err = file.WriteString(sizes); err != nil {
		log.Warn().Err(err).Msg("Could not save filter sizes to disk")
	}
	return err
}
