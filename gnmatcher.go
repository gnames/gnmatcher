package gnmatcher

import (
	"path/filepath"

	"github.com/gnames/gnmatcher/dbase"
	"github.com/gnames/gnmatcher/sys"
)

// GNmatcher keeps most general configuration settings and high level
// methods for scientific name matching.
type GNmatcher struct {
	WorkDir string
	JobsNum int
	dbase.Dbase
}

// NewGNmatcher is a constructor for GNmatcher instance
func NewGNmatcher(opts ...Option) GNmatcher {
	gnm := GNmatcher{
		WorkDir: "/tmp/gnmatcher",
		JobsNum: 4,
		Dbase:   dbase.NewDbase(),
	}
	for _, opt := range opts {
		opt(&gnm)
	}
	return gnm
}

func (gnm GNmatcher) FiltersDir() string {
	return filepath.Join(gnm.WorkDir, "bloom")
}

func (gnm GNmatcher) CreateWorkDir() error {
	return sys.MakeDir(gnm.FiltersDir())
}

// Option is a type of all options for GNmatcher.
type Option func(gnm *GNmatcher)

// OptWorkDir sets a directory for key-value stores and temporary files.
func OptWorkDir(d string) Option {
	return func(gnm *GNmatcher) {
		gnm.WorkDir = d
	}
}

// OptJobsNum sets number of concurrent jobs to run for parallel tasks.
func OptJobsNum(i int) Option {
	return func(gnm *GNmatcher) {
		gnm.JobsNum = i
	}
}

// OptPgHost sets the host of gnindex database
func OptPgHost(h string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgHost = h
	}
}

// OptPgUser sets the user of gnindex database
func OptPgUser(u string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgUser = u
	}
}

// OptPgPass sets the password to access gnindex database
func OptPgPass(p string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgPass = p
	}
}

// OptPgPort sets the port for gnindex database
func OptPgPort(p int) Option {
	return func(gnm *GNmatcher) {
		gnm.PgPort = p
	}
}

// OptPgDB sets the name of gnindex database
func OptPgDB(n string) Option {
	return func(gnm *GNmatcher) {
		gnm.PgDB = n
	}
}
