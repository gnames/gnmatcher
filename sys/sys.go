package sys

import (
	"fmt"
	"log"
	"os"
)

func MakeDir(dir string) error {
	path, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	if path.Mode().IsRegular() {
		return fmt.Errorf("'%s' is a file, not a directory", dir)
	}
	return nil
}

func FileExists(f string) bool {
	path, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	if !path.Mode().IsRegular() {
		log.Fatal(fmt.Errorf("'%s' is not a regular file, "+
			"delete or move it and try again", f))
	}
	return true
}
