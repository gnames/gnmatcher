// package cmd provides command line interface to http server that runs
// gnmatcher functionality.
package cmd

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/gnsys"

	"github.com/spf13/cobra"

	gnmatcher "github.com/gnames/gnmatcher/pkg"
	"github.com/gnames/gnmatcher/pkg/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

//go:embed gnmatcher.yaml
var configText string

var (
	opts []config.Option
)

// cfgData purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type cfgData struct {
	CacheDir    string
	JobsNum     int
	MaxEditDist int
	PgHost      string
	PgPort      int
	PgUser      string
	PgPass      string
	PgDB        string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnmatcher",
	Short: "Contains tools and algorithms to verify scientific names",
	Run: func(cmd *cobra.Command, _ []string) {
		if showVersionFlag(cmd) {
			os.Exit(0)
		}
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Cannot start gnmatcher", "error", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "V", false, "Return version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var homePath, cfgPath string
	var err error
	configFile := "gnmatcher"

	// Find home directory.
	homePath, err = homedir.Dir()
	if err != nil {
		slog.Error("Cannot find home directory", "error", err)
		os.Exit(1)
	}
	cfgPath = filepath.Join(homePath, ".config")

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(cfgPath)
	viper.SetConfigName(configFile)

	// Set environment variables to override
	// config file settings
	_ = viper.BindEnv("CacheDir", "GNM_CACHE_DIR")
	_ = viper.BindEnv("JobsNum", "GNM_JOBS_NUM")
	_ = viper.BindEnv("MaxEditDist", "GNM_MAX_EDIT_DIST")
	_ = viper.BindEnv("PgDB", "GNM_PG_DB")
	_ = viper.BindEnv("PgHost", "GNM_PG_HOST")
	_ = viper.BindEnv("PgPass", "GNM_PG_PASS")
	_ = viper.BindEnv("PgPort", "GNM_PG_PORT")
	_ = viper.BindEnv("PgUser", "GNM_PG_USER")

	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(cfgPath, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file", "file", viper.ConfigFileUsed())
	}
	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() []config.Option {
	cfg := &cfgData{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		slog.Error("Cannot deserialize config data", "error", err)
		os.Exit(1)
	}

	if cfg.CacheDir != "" {
		opts = append(opts, config.OptCacheDir(cfg.CacheDir))
	}
	if cfg.JobsNum > 0 {
		opts = append(opts, config.OptJobsNum(cfg.JobsNum))
	}
	if cfg.MaxEditDist != 0 {
		opts = append(opts, config.OptMaxEditDist(cfg.MaxEditDist))
	}
	if cfg.PgDB != "" {
		opts = append(opts, config.OptPgDB(cfg.PgDB))
	}
	if cfg.PgHost != "" {
		opts = append(opts, config.OptPgHost(cfg.PgHost))
	}
	if cfg.PgPass != "" {
		opts = append(opts, config.OptPgPass(cfg.PgPass))
	}
	if cfg.PgPort != 0 {
		opts = append(opts, config.OptPgPort(cfg.PgPort))
	}
	if cfg.PgUser != "" {
		opts = append(opts, config.OptPgUser(cfg.PgUser))
	}
	return opts
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, _ := cmd.Flags().GetBool("version")

	if hasVersionFlag {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnmatcher.Version, gnmatcher.Build)
	}
	return hasVersionFlag
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string) {
	if ok, err := gnsys.FileExists(configPath); ok && err == nil {
		return
	}

	slog.Info("Creating config file", "file", configPath)
	createConfig(configPath)
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		slog.Error("Cannot create dir", "path", path, "error", err)
		os.Exit(1)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		slog.Error("Cannot write to file", "path", path, "error", err)
		os.Exit(1)
	}
}
