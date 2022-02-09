// package cmd provides command line interface to http server that runs
// gnmatcher functionality.
package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnames/gnmatcher"
	"github.com/gnames/gnsys"

	"github.com/spf13/cobra"

	"github.com/gnames/gnmatcher/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
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
	CacheDir       string
	JobsNum        int
	MaxEditDist    int
	PgHost         string
	PgPort         int
	PgUser         string
	PgPass         string
	PgDB           string
	WebLogsNsqdTCP string
	WithWebLogs    bool
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
		log.Fatal().Err(err).Msg("Cannot start gnmatcher")
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
		log.Fatal().Err(err).Msg("Cannot find home directory")
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
	_ = viper.BindEnv("WebLogsNsqdTCP", "GNM_WEB_LOGS_NSQD_TCP")
	_ = viper.BindEnv("WithWebLogs", "GNM_WITH_WEB_LOGS")

	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(cfgPath, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s.", viper.ConfigFileUsed())
	}
	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() []config.Option {
	cfg := &cfgData{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot deserialize config data")
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
	if cfg.WebLogsNsqdTCP != "" {
		opts = append(opts, config.OptWebLogsNsqdTCP(cfg.WebLogsNsqdTCP))
	}
	if cfg.WithWebLogs {
		opts = append(opts, config.OptWithWebLogs(true))
	}
	return opts
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get version flag")
	}

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

	log.Printf("Creating config file: %s.", configPath)
	createConfig(configPath)
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create dir %s", path)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot write to file %s", path)
	}
}
