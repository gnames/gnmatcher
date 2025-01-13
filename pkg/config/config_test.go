package config_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/stretchr/testify/assert"
)

// NewConfig constructor
func TestNew(t *testing.T) {
	cfg := config.New()
	cacheDir, _ := os.UserCacheDir()
	deflt := config.Config{
		CacheDir:    filepath.Join(cacheDir, "gnmatcher"),
		MaxEditDist: 1,
		JobsNum:     4,
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
	oldLevel := slog.SetLogLoggerLevel(10)
	defer slog.SetLogLoggerLevel(oldLevel)

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
	assert.Contains(t, cfg.TrieDir(), "/gnmatcher/trie")
	assert.Contains(t, cfg.FiltersDir(), "/gnmatcher/bloom")
	assert.Contains(t, cfg.StemsDir(), "/gnmatcher/stems-kv")
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
