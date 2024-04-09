package bloom

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	baseBloomfilter "github.com/devopsfaith/bloomfilter/bloomfilter"
)

func saveFilters(path string, filters *bloomFilters) error {
	var err error
	var nilFilter *baseBloomfilter.Bloomfilter
	files := map[string]*baseBloomfilter.Bloomfilter{
		canonicalStemFile: filters.canonicalStem,
		sizesFile:         nilFilter,
	}

	for f, filter := range files {
		var file *os.File
		filePath := filepath.Join(path, f)

		file, err = createFile(filePath)
		if err != nil {
			slog.Error("Cannot create path", "path", filePath, "error", err)
			return err
		}
		if f == sizesFile {
			err = saveSizesFile(file, filters)
			if err != nil {
				slog.Error("Cannot create sizesFile", "error", err)
				return err
			}
			continue
		}

		err = saveFilterFile(filePath, file, filter)
		if err != nil {
			slog.Error("Cannot create file", "file", filePath, "error", err)
			return err
		}
	}

	slog.Info("Saved cached filters to disk")
	return nil
}

func createFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("Could not create file", "file", filePath, "error", err)
		return nil, err
	}

	return file, nil
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
		slog.Error("Could not serialize for file", "file", filePath, "error", err)
		return err
	}
	_, err = file.Write(bin)
	if err != nil {
		slog.Error("Could not save", "file", filePath, "error", err)
		return err
	}
	return nil
}

func saveSizesFile(file *os.File, filters *bloomFilters) error {
	var err error
	sizes := fmt.Sprintf("CanonicalSize,%d\n",
		filters.canonicalSize)
	if _, err = file.WriteString(sizes); err != nil {
		slog.Error("Could not save filter sizes to disk", "error", err)
		return err
	}
	return nil
}
