package gnmatcher

import (
  "fmt"
  "github.com/gnames/gnmatcher/dbase"
	log "github.com/sirupsen/logrus"
)

// Config collects and stores external configuration data.
type Config struct {
	WorkDir     string
	NatsURI     string
	JobsNum     int
	MaxEditDist int
	GNamesDB    dbase.Dbase
}

// NewConfig is a Config constructor that takes external options to
// update default values to external ones.
func NewConfig(opts ...Option) Config {
	cnf := Config{
		WorkDir:     "/tmp/gnmatcher",
		NatsURI:     "nats:localhost:4222",
		JobsNum:     8,
		MaxEditDist: 1,
		GNamesDB:    dbase.NewDbase(),
	}
	for _, opt := range opts {
		opt(&cnf)
	}
	return cnf
}

// Option is a type of all options for Config.
type Option func(cnf *Config)

// OptWorkDir sets a directory for key-value stores and temporary files.
func OptWorkDir(s string) Option {
	return func(cnf *Config) {
		cnf.WorkDir = s
	}
}

// OptNatsURI defines a URI to connect to NATS messaging service server.
func OptNatsURI(s string) Option {
	return func(cnf *Config) {
		cnf.NatsURI = s
	}
}

// OptJobsNum sets number of concurrent jobs to run for parallel tasks.
func OptJobsNum(i int) Option {
	return func(cnf *Config) {
		cnf.JobsNum = i
	}
}

// OptMaxEditDist sets maximal possible edit distance for fuzzy matching of
// stemmed canonical forms.
func OptMaxEditDist(i int) Option {
	return func(cnf *Config) {
  if i < 1 || i > 2 {
    log.Warn(fmt.Sprintf("MaxEditDist can only be 1 or 2, leaving it at %d.",
    cnf.MaxEditDist))
  }
		cnf.MaxEditDist = i
	}
}

// OptPgHost sets the host of gnames database
func OptPgHost(s string) Option {
	return func(cnf *Config) {
		cnf.GNamesDB.PgHost = s
	}
}

// OptPgUser sets the user of gnnames database
func OptPgUser(s string) Option {
	return func(cnf *Config) {
		cnf.GNamesDB.PgUser = s
	}
}

// OptPgPass sets the password to access gnnames database
func OptPgPass(s string) Option {
	return func(cnf *Config) {
		cnf.GNamesDB.PgPass = s
	}
}

// OptPgPort sets the port for gnames database
func OptPgPort(i int) Option {
	return func(cnf *Config) {
		cnf.GNamesDB.PgPort = i
	}
}

// OptPgDB sets the name of gnames database
func OptPgDB(s string) Option {
	return func(cnf *Config) {
		cnf.GNamesDB.PgDB = s
	}
}
