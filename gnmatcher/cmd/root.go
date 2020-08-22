/*
Copyright © 2020 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gnames/gnmatcher"
	"github.com/gnames/gnmatcher/sys"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const configText = `# Path to keep working data and key-value stores
WorkDir: /var/gnmatcher

# URI to NATS, a messaging system server
NatsURI: nats:localhost:4222

# Postgresql host for gnames database
PgHost: localhost

# Postgresql user
PgUser: postgres

# Postgresql password
PgPass:

# Postgresql database
PgDB: gnames

# Number of jobs for parallel tasks
JobsNum: 4

# MaxEditDist is the maximal edit distance for fuzzy matching of
# stemmed canonical forms. Can be 1 or 2, 2 is significantly slower.
MaxEditDist: 1
`

var (
	opts []gnmatcher.Option
)

// config purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type config struct {
	WorkDir string
	NatsURI string
	PgHost  string
	PgPort  int
	PgUser  string
	PgPass  string
	PgDB    string
	JobsNum int
  MaxEditDist int
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
		log.Fatalf("Cannot start gnmatcher: %s.", err)
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
		log.Fatalf("Cannot find home directory: %s.", err)
	}
	home = filepath.Join(home, ".config")

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(configFile)

	// Set environment variables to override
	// config file settings
	viper.BindEnv("WorkDir", "GNM_WORK_DIR")
	viper.BindEnv("NatsURI", "GNM_NATS_URI")
	viper.BindEnv("PgHost", "GNM_PG_HOST")
	viper.BindEnv("PgPort", "GNM_PG_PORT")
	viper.BindEnv("PgUser", "GNM_PG_USER")
	viper.BindEnv("PgPass", "GNM_PG_PASS")
	viper.BindEnv("PgDB", "GNM_PG_DB")
	viper.BindEnv("JobsNum", "GNM_JOBS_NUM")
	viper.BindEnv("MaxEditDist", "GNM_MAX_EDIT_DIST")

	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(home, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath, configFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s.", viper.ConfigFileUsed())
	}
	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() []gnmatcher.Option {
	cfg := &config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Cannot deserialize config data: %s.", err)
	}

	if cfg.WorkDir != "" {
		opts = append(opts, gnmatcher.OptWorkDir(cfg.WorkDir))
	}
	if cfg.NatsURI != "" {
		opts = append(opts, gnmatcher.OptNatsURI(cfg.NatsURI))
	}
	if cfg.JobsNum != 0 {
		opts = append(opts, gnmatcher.OptJobsNum(cfg.JobsNum))
	}
	if cfg.MaxEditDist != 0 {
		opts = append(opts, gnmatcher.OptMaxEditDist(cfg.MaxEditDist))
	}
	if cfg.PgHost != "" {
		opts = append(opts, gnmatcher.OptPgHost(cfg.PgHost))
	}
	if cfg.PgPort != 0 {
		opts = append(opts, gnmatcher.OptPgPort(cfg.PgPort))
	}
	if cfg.PgUser != "" {
		opts = append(opts, gnmatcher.OptPgUser(cfg.PgUser))
	}
	if cfg.PgPass != "" {
		opts = append(opts, gnmatcher.OptPgPass(cfg.PgPass))
	}
	if cfg.PgDB != "" {
		opts = append(opts, gnmatcher.OptPgDB(cfg.PgDB))
	}
	return opts
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatal("Cannot get version flag: %s.", err)
	}

	if hasVersionFlag {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnmatcher.Version, gnmatcher.Build)
	}
	return hasVersionFlag
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string, configFile string) {
	if sys.FileExists(configPath) {
		return
	}

	log.Printf("Creating config file: %s.", configPath)
	createConfig(configPath, configFile)
}

// createConfig creates config file.
func createConfig(path string, file string) {
	err := sys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Fatal("Cannot create dir %s: %s.", path, err)
	}

	err = ioutil.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Fatal("Cannot write to file %s: %s", path, err)
	}
}
