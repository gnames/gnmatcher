package config_test

import (
	"testing"

	"github.com/gnames/gnmatcher/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// NewConfig constructor
func TestNew(t *testing.T) {
	cnf := config.NewConfig()
	deflt := config.Config{
		WorkDir:     config.ConvertTilda("~/.local/share/gnmatcher"),
		JobsNum:     8,
		MaxEditDist: 1,
		PgHost:      "localhost",
		PgPort:      5432,
		PgUser:      "postgres",
		PgPass:      "",
		PgDB:        "gnames",
	}
	assert.Equal(t, cnf, deflt)
}

// NewConfig with opts
func TestNewOpts(t *testing.T) {
	opts := opts()
	cnf := config.NewConfig(opts...)
	withOpts := config.Config{
		WorkDir:     "/var/opt/gnmatcher",
		JobsNum:     16,
		MaxEditDist: 2,
		PgHost:      "mypg",
		PgPort:      1234,
		PgUser:      "gnm",
		PgPass:      "secret",
		PgDB:        "gnm",
	}
	assert.Equal(t, cnf, withOpts)
}

// 	MaxEditDist is limited to 1 or 2
func TestMaxED(t *testing.T) {
	log.SetLevel(log.PanicLevel)
	cnf := config.NewConfig(config.OptMaxEditDist(5))
	assert.Equal(t, cnf.MaxEditDist, 1)
	cnf = config.NewConfig(config.OptMaxEditDist(0))
	assert.Equal(t, cnf.MaxEditDist, 1)
	cnf = config.NewConfig(config.OptMaxEditDist(1))
	assert.Equal(t, cnf.MaxEditDist, 1)
	cnf = config.NewConfig(config.OptMaxEditDist(2))
	assert.Equal(t, cnf.MaxEditDist, 2)
}

func TestHelpers(t *testing.T) {
	cnf := config.NewConfig()
	assert.Contains(t, cnf.TrieDir(), "/.local/share/gnmatcher/levenshein")
	assert.Contains(t, cnf.FiltersDir(), "/.local/share/gnmatcher/bloom")
	assert.Contains(t, cnf.StemsDir(), "/.local/share/gnmatcher/stems-kv")
}

func opts() []config.Option {
	return []config.Option{
		config.OptWorkDir("/var/opt/gnmatcher"),
		config.OptJobsNum(16),
		config.OptMaxEditDist(2),
		config.OptPgHost("mypg"),
		config.OptPgUser("gnm"),
		config.OptPgPass("secret"),
		config.OptPgPort(1234),
		config.OptPgDB("gnm"),
	}
}
