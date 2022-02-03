package bloom

import (
	"fmt"
	"os"
	"path/filepath"

	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
	log "github.com/sirupsen/logrus"
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
			log.Fatalf("Cannot create %s: %s", filePath, err)
		}
		if f == sizesFile {
			err = saveSizesFile(file, filters)
			if err != nil {
				log.Fatalf("Cannot create sizesFile: %s", err)
			}
			continue
		}

		err = saveFilterFile(filePath, file, filter)
		if err != nil {
			log.Fatalf("Cannot create %s: %s", filePath, err)
		}
	}

	if err == nil {
		log.Print("Saved cached filters to disk.")
	}
}

func createFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		warn := fmt.Sprintf("Could not create file %s: %s.", filePath, err)
		log.Warning(warn)
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
		warn := fmt.Sprintf("Could not serialize for %s: %s.", filePath, err)
		log.Warning(warn)
	}
	_, err = file.Write(bin)
	if err != nil {
		warn := fmt.Sprintf("Could not save %s: %s.", filePath, err)
		log.Warning(warn)
	}
	return err
}

func saveSizesFile(file *os.File, filters *bloomFilters) error {
	var err error
	sizes := fmt.Sprintf("CanonicalSize,%d\n",
		filters.canonicalSize)
	if _, err = file.WriteString(sizes); err != nil {
		warn := fmt.Sprintf("Could not save filter sizes to disk: %s.", err)
		log.Warn(warn)
	}
	return err
}
