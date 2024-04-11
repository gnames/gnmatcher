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
	"log/slog"
	"os"

	"github.com/gnames/gnfmt"
	gnmatcher "github.com/gnames/gnmatcher/pkg"
	gnmcnf "github.com/gnames/gnmatcher/pkg/config"
	"github.com/gnames/gnmatcher/pkg/io/bloom"
	"github.com/gnames/gnmatcher/pkg/io/rest"
	"github.com/gnames/gnmatcher/pkg/io/trie"
	"github.com/gnames/gnmatcher/pkg/io/virusio"

	"github.com/spf13/cobra"
)

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "RESTful interface to scientific names matching.",
	Long: `Runs a RESTful HTTP/1 server that takes a list of scientific names
in binary protobuf-based format and returns output in protobuf format
as well.`,
	Run: func(cmd *cobra.Command, _ []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
		port, _ := cmd.Flags().GetInt("port")

		cfg := gnmcnf.New(opts...)
		em := bloom.New(cfg)
		fm := trie.New(cfg)
		vm := virusio.New(cfg)
		gnm, err := gnmatcher.New(em, fm, vm, cfg)
		if err != nil {
			slog.Error("Error initializing matcher", "error", err)
			os.Exit(1)
		}

		var enc gnfmt.Encoder = gnfmt.GNjson{}

		service := rest.NewMatcherService(gnm, port, enc)
		rest.Run(service)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(restCmd)

	restCmd.Flags().IntP("port", "p", 8080, "REST port")
	restCmd.Flags().BoolP("debug", "d", false, "set logs level to DEBUG")
}
