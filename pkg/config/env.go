package config

import (
	"log/slog"
	"os"
	"strconv"
)

// LoadEnv allows to change configuration during runtime without CLI.
func LoadEnv(c *Config) {
	slog.Info("Updating config using environment variables")
	opts := strOpts()
	opts = append(opts, intOpts()...)
	for _, opt := range opts {
		opt(c)
	}
}

func strOpts() []Option {
	var res []Option

	envToOpt := map[string]func(string) Option{
		"GNM_PG_HOST":   OptPgHost,
		"GNM_PG_USER":   OptPgUser,
		"GNM_PG_PASS":   OptPgPass,
		"GNM_PG_DB":     OptPgDB,
		"GNM_CACHE_DIR": OptCacheDir,
	}

	for envVar, optFunc := range envToOpt {
		envVal := os.Getenv(envVar)
		if envVal != "" {
			res = append(res, optFunc(envVal))
		}
	}

	return res
}

func intOpts() []Option {
	var res []Option
	envToOpt := map[string]func(int) Option{
		"GNM_PG_PORT":       OptPgPort,
		"GNM_JOBS_NUM":      OptJobsNum,
		"GNM_MAX_EDIT_DIST": OptMaxEditDist,
	}
	for envVar, optFunc := range envToOpt {
		if envVar == "" {
			continue
		}
		val := os.Getenv(envVar)
		i, err := strconv.Atoi(val)
		if err != nil {
			slog.Warn("Cannot convert to int", "env", envVar, "value", val)
			continue
		}
		res = append(res, optFunc(i))
	}
	return res
}
