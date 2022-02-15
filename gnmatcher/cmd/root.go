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
	WorkDir            string
	PgHost             string
	PgPort             int
	PgUser             string
	PgPass             string
	PgDB               string
	MaxEditDist        int
	JobsNum            int
	NsqdTCPAddress     string
	NsqdContainsFilter string
	NsqdRegexFilter    string
	WithWebLogs        bool
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnmatcher",
	Short: "Contains tools and algorithms to verify scientific names",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersionFlag(cmd) {
			os.Exit(0)
		}
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.  This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
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
	var home string
	var err error
	configFile := "gnmatcher"

	// Find home directory.
	home, err = homedir.Dir()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot find home directory")
	}
	home = filepath.Join(home, ".config")

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(configFile)

	// Set environment variables to override
	// config file settings
	_ = viper.BindEnv("WorkDir", "GNM_WORK_DIR")
	_ = viper.BindEnv("PgHost", "GNM_PG_HOST")
	_ = viper.BindEnv("PgPort", "GNM_PG_PORT")
	_ = viper.BindEnv("PgUser", "GNM_PG_USER")
	_ = viper.BindEnv("PgPass", "GNM_PG_PASS")
	_ = viper.BindEnv("PgDB", "GNM_PG_DB")
	_ = viper.BindEnv("MaxEditDist", "GNM_MAX_EDIT_DIST")
	_ = viper.BindEnv("JobsNum", "GNM_JOBS_NUM")
	_ = viper.BindEnv("NsqdTCPAddress", "GNM_NSQD_TCP_ADDRESS")
	_ = viper.BindEnv("NsqdContainsFilter", "GNM_NSQD_CONTAINS_FILTER")
	_ = viper.BindEnv("NsqdRegexFilter", "GNM_NSQD_REGEX_FILTER")
	_ = viper.BindEnv("WithWebLogs", "GNM_WITH_WEB_LOGS")

	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(home, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath, configFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s.", viper.ConfigFileUsed())
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

	if cfg.WorkDir != "" {
		opts = append(opts, config.OptWorkDir(cfg.WorkDir))
	}
	if cfg.MaxEditDist != 0 {
		opts = append(opts, config.OptMaxEditDist(cfg.MaxEditDist))
	}
	if cfg.PgHost != "" {
		opts = append(opts, config.OptPgHost(cfg.PgHost))
	}
	if cfg.PgPort != 0 {
		opts = append(opts, config.OptPgPort(cfg.PgPort))
	}
	if cfg.PgUser != "" {
		opts = append(opts, config.OptPgUser(cfg.PgUser))
	}
	if cfg.PgPass != "" {
		opts = append(opts, config.OptPgPass(cfg.PgPass))
	}
	if cfg.PgDB != "" {
		opts = append(opts, config.OptPgDB(cfg.PgDB))
	}
	if cfg.JobsNum > 0 {
		opts = append(opts, config.OptJobsNum(cfg.JobsNum))
	}
	if cfg.NsqdContainsFilter != "" {
		opts = append(opts, config.OptNsqdContainsFilter(cfg.NsqdContainsFilter))
	}
	if cfg.NsqdRegexFilter != "" {
		opts = append(opts, config.OptNsqdRegexFilter(cfg.NsqdRegexFilter))
	}
	if cfg.NsqdTCPAddress != "" {
		opts = append(opts, config.OptNsqdTCPAddress(cfg.NsqdTCPAddress))
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
func touchConfigFile(configPath string, configFile string) {
	if ok, err := gnsys.FileExists(configPath); ok && err == nil {
		return
	}

	log.Info().Msgf("Creating config file: %s.", configPath)
	createConfig(configPath, configFile)
}

// createConfig creates config file.
func createConfig(path string, file string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create dir %s", path)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot write to file %s", path)
	}
}
