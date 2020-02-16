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
	"log"
	"os"

	"github.com/gnames/gnmatcher"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	opts    []gnmatcher.Option
)

// config purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type config struct {
	WorkDir string
	PgHost  string
	PgPort  int
	PgUser  string
	PgPass  string
	PgDB    string
	JobsNum int
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnmatcher",
	Short: "Contains tools and algorithms to verify scientific names",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if version {
			fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnmatcher.Version, gnmatcher.Build)
			os.Exit(0)
		}

		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gnmatcher.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "V", false, "Return version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gnmatcher" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gnmatcher")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Config file $HOME/.gnidump.yaml not found")
	}
	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() []gnmatcher.Option {
	cfg := &config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.WorkDir != "" {
		opts = append(opts, gnmatcher.OptWorkDir(cfg.WorkDir))
	}
	if cfg.JobsNum != 0 {
		opts = append(opts, gnmatcher.OptJobsNum(cfg.JobsNum))
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
