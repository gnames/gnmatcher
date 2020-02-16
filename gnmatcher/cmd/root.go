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
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

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
		versionFlag(cmd)

		if len(args) == 0 {
			processStdin(cmd)
			os.Exit(0)
		}
		data := getInput(cmd, args)
		match(data)
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.  This is called by main.main(). It only needs to happen once
// to the rootCmd.
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.gnmatcher.yaml)")

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

func versionFlag(cmd *cobra.Command) {
	version, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if version {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnmatcher.Version, gnmatcher.Build)
		os.Exit(0)
	}
}

func getInput(cmd *cobra.Command, args []string) string {
	var data string
	switch len(args) {
	case 1:
		data = args[0]
	default:
		_ = cmd.Help()
		os.Exit(0)
	}
	return data
}

func match(data string) {
	gnm := gnmatcher.NewGNmatcher(opts...)
	err := gnm.CreateWorkDir()
	if err != nil {
		log.Fatal(err)
	}

	path := string(data)
	if fileExists(path) {
		f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		matchFile(gnm, f)
		f.Close()
	} else {
		matchString(gnm, data)
	}
}

func processStdin(cmd *cobra.Command) {
	if !checkStdin() {
		_ = cmd.Help()
		return
	}
	gnm := gnmatcher.NewGNmatcher(opts...)
	matchFile(gnm, os.Stdin)
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		log.Panic(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func fileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

func matchFile(gnm gnmatcher.GNmatcher, f io.Reader) {
	in := make(chan string)
	out := make(chan gnmatcher.MatchResult)
	var wg sync.WaitGroup
	wg.Add(1)

	go gnm.MatchStream(in, out)
	go processResults(gnm, out, &wg)
	sc := bufio.NewScanner(f)
	count := 0
	for sc.Scan() {
		count++
		if count%50000 == 0 {
			log.Printf("Matching %d-th line\n", count)
		}
		name := sc.Text()
		in <- name
	}
	close(in)
	wg.Wait()
}

func processResults(gnm gnmatcher.GNmatcher,
	out <-chan gnmatcher.MatchResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for r := range out {
		if r.Error != nil {
			log.Println(r.Error)
		}
		fmt.Println(r.Output)
	}
}

func matchString(gnm gnmatcher.GNmatcher, data string) {
	res, err := gnm.MatchAndFormat(data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(res)
}
