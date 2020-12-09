package config_test

import (
	"testing"

	"github.com/gnames/gnlib/sys"
	"github.com/gnames/gnmatcher/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// NewConfig constructor
func TestNew(t *testing.T) {
	cfg := config.NewConfig()
	deflt := config.Config{
		WorkDir:     sys.ConvertTilda("~/.local/share/gnmatcher"),
		MaxEditDist: 1,
		JobsNum:     1,
		PgHost:      "localhost",
		PgPort:      5432,
		PgUser:      "postgres",
		PgPass:      "",
		PgDB:        "gnames",
	}
	assert.Equal(t, cfg, deflt)
}

// NewConfig with opts
func TestNewOpts(t *testing.T) {
	opts := opts()
	cfg := config.NewConfig(opts...)
	withOpts := config.Config{
		WorkDir:     "/var/opt/gnmatcher",
		MaxEditDist: 2,
		JobsNum:     16,
		PgHost:      "mypg",
		PgPort:      1234,
		PgUser:      "gnm",
		PgPass:      "secret",
		PgDB:        "gnm",
	}
	assert.Equal(t, cfg, withOpts)
}

// 	MaxEditDist is limited to 1 or 2
func TestMaxED(t *testing.T) {
	log.SetLevel(log.PanicLevel)
	cfg := config.NewConfig(config.OptMaxEditDist(5))
	assert.Equal(t, cfg.MaxEditDist, 1)
	cfg = config.NewConfig(config.OptMaxEditDist(0))
	assert.Equal(t, cfg.MaxEditDist, 1)
	cfg = config.NewConfig(config.OptMaxEditDist(1))
	assert.Equal(t, cfg.MaxEditDist, 1)
	cfg = config.NewConfig(config.OptMaxEditDist(2))
	assert.Equal(t, cfg.MaxEditDist, 2)
}

func TestHelpers(t *testing.T) {
	cfg := config.NewConfig()
	assert.Contains(t, cfg.TrieDir(), "/.local/share/gnmatcher/trie")
	assert.Contains(t, cfg.FiltersDir(), "/.local/share/gnmatcher/bloom")
	assert.Contains(t, cfg.StemsDir(), "/.local/share/gnmatcher/stems-kv")
}

func opts() []config.Option {
	return []config.Option{
		config.OptWorkDir("/var/opt/gnmatcher"),
		config.OptMaxEditDist(2),
		config.OptJobsNum(16),
		config.OptPgHost("mypg"),
		config.OptPgUser("gnm"),
		config.OptPgPass("secret"),
		config.OptPgPort(1234),
		config.OptPgDB("gnm"),
	}
}
