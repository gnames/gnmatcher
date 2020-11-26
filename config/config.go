// package config contains information needed to run gnmatcher project.
package config

import (
	"fmt"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

// Config collects and stores external configuration data.
type Config struct {
	// WorkDir is the main directory for gnmatcher files. It contains
	// bloom filters levenshtein automata trees, key-value stores etc.
	WorkDir string
	// MaxEditDist is the maximal allowed edit distance for levenshtein automata.
	// The number cannot exceed 2, default number is 1. The speed of execution
	// slows down dramatically with the MaxEditDist > 1.
	MaxEditDist int
	// PgHost is a hostname for the PostgreSQL server.
	PgHost string
	// PgPort is the port of PostgreSQL server.
	PgPort int
	// PgUser is the user for the database.
	PgUser string
	// PgPass password to access PostgreSQL server.
	PgPass string
	// PgDB the database name where gnames data is located.
	PgDB string
}

// NewConfig is a Config constructor that takes external options to
// update default values to external ones.
func NewConfig(opts ...Option) Config {
	workDir := "~/.local/share/gnmatcher"
	cnf := Config{
		WorkDir:     ConvertTilda(workDir),
		MaxEditDist: 1,
		PgHost:      "localhost",
		PgPort:      5432,
		PgUser:      "postgres",
		PgPass:      "",
		PgDB:        "gnames",
	}
	for _, opt := range opts {
		opt(&cnf)
	}
	return cnf
}

// TrieDir returns path where to dump/restore
// serialized trie.
func (cnf Config) TrieDir() string {
	return filepath.Join(cnf.WorkDir, "trie")
}

// FiltersDir returns path where to dump/restore
// serialized bloom filters.
func (cnf Config) FiltersDir() string {
	return filepath.Join(cnf.WorkDir, "bloom")
}

// StemsDir returns path where stems key-value store
// is located
func (cnf Config) StemsDir() string {
	return filepath.Join(cnf.WorkDir, "stems-kv")
}

// Option is a type of all options for Config.
type Option func(cnf *Config)

// OptWorkDir sets a directory for key-value stores and temporary files.
func OptWorkDir(s string) Option {
	return func(cnf *Config) {
		cnf.WorkDir = ConvertTilda(s)
	}
}

// OptMaxEditDist sets maximal possible edit distance for fuzzy matching of
// stemmed canonical forms.
func OptMaxEditDist(i int) Option {
	return func(cnf *Config) {
		if i < 1 || i > 2 {
			log.Warn(fmt.Sprintf("MaxEditDist can only be 1 or 2, leaving it at %d.",
				cnf.MaxEditDist))
		} else {
			cnf.MaxEditDist = i
		}
	}
}

// OptPgHost sets the host of gnames database
func OptPgHost(s string) Option {
	return func(cnf *Config) {
		cnf.PgHost = s
	}
}

// OptPgUser sets the user of gnnames database
func OptPgUser(s string) Option {
	return func(cnf *Config) {
		cnf.PgUser = s
	}
}

// OptPgPass sets the password to access gnnames database
func OptPgPass(s string) Option {
	return func(cnf *Config) {
		cnf.PgPass = s
	}
}

// OptPgPort sets the port for gnames database
func OptPgPort(i int) Option {
	return func(cnf *Config) {
		cnf.PgPort = i
	}
}

// OptPgDB sets the name of gnames database
func OptPgDB(s string) Option {
	return func(cnf *Config) {
		cnf.PgDB = s
	}
}

// ConvertTilda expands paths with `~/` to an actual home directory.
func ConvertTilda(path string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		path = filepath.Join(home, path[2:])
	}
	return path
}
