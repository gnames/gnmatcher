// package config contains information needed to run gnmatcher project.
package config

import (
	"path/filepath"
	"regexp"

	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
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
	// JobsNum is the number of jobs to run in parallel
	JobsNum int
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

	// NsqdTCPAddress provides an address to the NSQ messenger TCP service. If
	// this value is set and valid, the web logs will be published to the NSQ.
	// The option is ignored if `Port` is not set.
	//
	// If WithWebLogs option is set to `false`, but `NsqdTCPAddress` is set to a
	// valid URL, the logs will be sent to the NSQ messanging service, but they
	// wil not appear as STRERR output.
	// Example: `127.0.0.1:4150`
	NsqdTCPAddress string

	// NsqdContainsFilter logs should match the filter to be sent to NSQ
	// service.
	// Examples:
	// "api" - logs should contain "api"
	// "!api" - logs should not contain "api"
	NsqdContainsFilter string

	// NsqdRegexFilter logs should match the regular expression to be sent to
	// NSQ service.
	// Example: `api\/v(0|1)`
	NsqdRegexFilter *regexp.Regexp

	// WithWebLogs flag enables logs when running web-service. This flag is
	// ignored if `Port` value is not set.
	WithWebLogs bool
}

// NewConfig is a Config constructor that takes external options to
// update default values to external ones.
func NewConfig(opts ...Option) Config {
	workDir := "~/.local/share/gnmatcher"
	workDir, _ = gnsys.ConvertTilda(workDir)
	cfg := Config{
		WorkDir:     workDir,
		MaxEditDist: 1,
		JobsNum:     1,
		PgHost:      "localhost",
		PgPort:      5432,
		PgUser:      "postgres",
		PgPass:      "",
		PgDB:        "gnames",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// TrieDir returns path where to dump/restore
// serialized trie.
func (cfg Config) TrieDir() string {
	return filepath.Join(cfg.WorkDir, "trie")
}

// FiltersDir returns path where to dump/restore
// serialized bloom filters.
func (cfg Config) FiltersDir() string {
	return filepath.Join(cfg.WorkDir, "bloom")
}

// StemsDir returns path where stems key-value store
// is located
func (cfg Config) StemsDir() string {
	return filepath.Join(cfg.WorkDir, "stems-kv")
}

// Option is a type of all options for Config.
type Option func(cfg *Config)

// OptWorkDir sets a directory for key-value stores and temporary files.
func OptWorkDir(s string) Option {
	return func(cfg *Config) {
		workDir, err := gnsys.ConvertTilda(s)
		if err != nil {
			log.Warn().Err(err).Msgf("Cannot expand '%s'", s)
		}
		cfg.WorkDir = workDir
	}
}

// OptMaxEditDist sets maximal possible edit distance for fuzzy matching of
// stemmed canonical forms.
func OptMaxEditDist(i int) Option {
	return func(cfg *Config) {
		if i < 1 || i > 2 {
			log.Warn().Msgf("MaxEditDist can only be 1 or 2, leaving it at %d.",
				cfg.MaxEditDist)
		} else {
			cfg.MaxEditDist = i
		}
	}
}

// OptPgHost sets the host of gnames database
func OptPgHost(s string) Option {
	return func(cfg *Config) {
		cfg.PgHost = s
	}
}

// OptPgUser sets the user of gnnames database
func OptPgUser(s string) Option {
	return func(cfg *Config) {
		cfg.PgUser = s
	}
}

// OptPgPass sets the password to access gnnames database
func OptPgPass(s string) Option {
	return func(cfg *Config) {
		cfg.PgPass = s
	}
}

// OptPgPort sets the port for gnames database
func OptPgPort(i int) Option {
	return func(cfg *Config) {
		cfg.PgPort = i
	}
}

// OptPgDB sets the name of gnames database
func OptPgDB(s string) Option {
	return func(cfg *Config) {
		cfg.PgDB = s
	}
}

// OptJobsNum sets the number of jobs to run in parallel
func OptJobsNum(i int) Option {
	return func(cfg *Config) {
		cfg.JobsNum = i
	}
}

// OptNsqdTCPAddress provides an address of NSQ messanging service.
func OptNsqdTCPAddress(s string) Option {
	return func(cfg *Config) {
		cfg.NsqdTCPAddress = s
	}
}

// OptNsqdContainsFilter provides a filter for logs sent to NSQ service.
func OptNsqdContainsFilter(s string) Option {
	return func(cfg *Config) {
		cfg.NsqdContainsFilter = s
	}
}

// OptNsqdRegexFilter provides a regular expression filter for
// logs sent to NSQ service.
func OptNsqdRegexFilter(s string) Option {
	return func(cfg *Config) {
		r := regexp.MustCompile(s)
		cfg.NsqdRegexFilter = r
	}
}

// OptWithWebLogs sets the WithWebLogs field.
func OptWithWebLogs(b bool) Option {
	return func(cfg *Config) {
		cfg.WithWebLogs = b
	}
}
