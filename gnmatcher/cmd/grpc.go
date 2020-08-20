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
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gnames/gnmatcher"
	"github.com/gnames/gnmatcher/rpc"
	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "A gRPC interface to name matching functionality.",
	Long: `Runs a gRPC server that listens for packages of scientific names. It
tries to match the names using exact and fuzzy matching algorithms and returns
UUIDs of canonical forms that did match together with the edit distances to
estimate differences between input and output names.`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			log.SetLevel(log.DebugLevel)
			log.Printf("Log level is set to '%s'.", log.Level.String(log.GetLevel()))
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatalf("Cannot get port flag: %s", err)
		}
		cnf := gnmatcher.NewConfig(opts...)
		gnm, err := gnmatcher.NewGNMatcher(cnf)
		if err != nil {
			log.Printf("Cannot create an instance of GNMatcher: %s.", err)
			os.Exit(1)
		}
		rpc.Run(port, &gnm)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	grpcCmd.Flags().IntP("port", "p", 8778, "grpc's port")
	grpcCmd.Flags().BoolP("debug", "d", false, "set logs level to DEBUG")
}
