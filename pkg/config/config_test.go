package config_test

import (
	"testing"

	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// NewConfig constructor
func TestNew(t *testing.T) {
	cfg := config.New()
	cacheDir, _ := gnsys.ConvertTilda("~/.cache/gnmatcher")
	deflt := config.Config{
		CacheDir:    cacheDir,
		MaxEditDist: 1,
		JobsNum:     1,
		PgHost:      "0.0.0.0",
		PgPort:      5432,
		PgUser:      "postgres",
		PgPass:      "postgres",
		PgDB:        "gnames",
	}
	assert.Equal(t, deflt, cfg)
}

// NewConfig with opts
func TestNewOpts(t *testing.T) {
	opts := opts()
	cfg := config.New(opts...)
	withOpts := config.Config{
		CacheDir:    "/var/opt/gnmatcher",
		MaxEditDist: 2,
		JobsNum:     16,
		PgHost:      "mypg",
		PgPort:      1234,
		PgUser:      "gnm",
		PgPass:      "secret",
		PgDB:        "gnm",
	}
	assert.Equal(t, withOpts, cfg)
}

// MaxEditDist is limited to 1 or 2
func TestMaxED(t *testing.T) {
	logLevel := log.Logger.GetLevel()
	zerolog.SetGlobalLevel(zerolog.Disabled)
	defer zerolog.SetGlobalLevel(logLevel)

	cfg := config.New(config.OptMaxEditDist(5))
	assert.Equal(t, 1, cfg.MaxEditDist)
	cfg = config.New(config.OptMaxEditDist(0))
	assert.Equal(t, 1, cfg.MaxEditDist)
	cfg = config.New(config.OptMaxEditDist(1))
	assert.Equal(t, 1, cfg.MaxEditDist)
	cfg = config.New(config.OptMaxEditDist(2))
	assert.Equal(t, 2, cfg.MaxEditDist)
}

func TestHelpers(t *testing.T) {
	cfg := config.New()
	assert.Contains(t, cfg.TrieDir(), "/.cache/gnmatcher/trie")
	assert.Contains(t, cfg.FiltersDir(), "/.cache/gnmatcher/bloom")
	assert.Contains(t, cfg.StemsDir(), "/.cache/gnmatcher/stems-kv")
}

func opts() []config.Option {
	return []config.Option{
		config.OptCacheDir("/var/opt/gnmatcher"),
		config.OptMaxEditDist(2),
		config.OptJobsNum(16),
		config.OptPgHost("mypg"),
		config.OptPgUser("gnm"),
		config.OptPgPass("secret"),
		config.OptPgPort(1234),
		config.OptPgDB("gnm"),
	}
}
