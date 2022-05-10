package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	nsqcfg "github.com/sfgrp/lognsq/config"
	"github.com/sfgrp/lognsq/ent/nsq"
	"github.com/sfgrp/lognsq/io/nsqio"
)

var apiPath = "/api/v0/"

// Run creates and runs a RESTful API service of gnmatcher.
// this API is described by OpenAPI schema at
// https://app.swaggerhub.com/apis/dimus/gnmatcher/1.0.0
func Run(m MatcherService) {
	log.Info().Int("port", m.Port()).Msg("Starting HTTP API server")
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	loggerNSQ := setLogger(e, m)
	if loggerNSQ != nil {
		defer loggerNSQ.Stop()
	}

	e.GET("/", root)
	e.GET(apiPath+"ping", ping(m))
	e.GET(apiPath+"version", ver(m))
	e.POST(apiPath+"matches", matchPOST(m))
	e.GET(apiPath+"matches/:names", matchGET(m))

	addr := fmt.Sprintf(":%d", m.Port())
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

func root(c echo.Context) error {
	return c.String(http.StatusOK,
		`The OpenAPI is described at
https://apidoc.globalnames.org/gnmatcher`)
}

func ping(m MatcherService) func(echo.Context) error {
	return func(c echo.Context) error {
		result := m.Ping()
		return c.String(http.StatusOK, result)
	}
}

func ver(m MatcherService) func(echo.Context) error {
	return func(c echo.Context) error {
		result := m.GetVersion()
		return c.JSON(http.StatusOK, result)
	}
}

func matchGET(m MatcherService) func(echo.Context) error {
	fmt.Println("HERE")
	return func(c echo.Context) error {
		nameStr, _ := url.QueryUnescape(c.Param("names"))
		names := strings.Split(nameStr, "|")
		dsStr, _ := url.QueryUnescape(c.QueryParam("data_sources"))
		var ds []int
		for _, v := range strings.Split(dsStr, "|") {
			if id, err := strconv.Atoi(v); err == nil {
				ds = append(ds, id)
			}
		}
		var opts []config.Option
		spGrp := c.QueryParam("species_group") == "true"
		if spGrp {
			opts = append(opts, config.OptWithSpeciesGroup(true))
		}
		if len(ds) > 0 {
			opts = append(opts, config.OptDataSources(ds))
		}

		result := m.MatchNames(names, opts...)
		if l := len(names); l > 0 {
			log.Info().
				Int("namesNum", l).
				Str("example", names[0]).
				Str("method", "POST").
				Msg("Name Match")
		}
		return c.JSON(http.StatusOK, result)
	}
}

func matchPOST(m MatcherService) func(echo.Context) error {
	return func(c echo.Context) error {
		var inp mlib.Input
		var opts []config.Option

		if err := c.Bind(&inp); err != nil {
			return err
		}
		if inp.WithSpeciesGroup {
			opts = append(opts, config.OptWithSpeciesGroup(true))
		}
		if len(inp.DataSources) > 0 {
			opts = append(opts, config.OptDataSources(inp.DataSources))
		}

		result := m.MatchNames(inp.Names, opts...)
		if l := len(inp.Names); l > 0 {
			log.Info().
				Int("namesNum", l).
				Str("example", inp.Names[0]).
				Str("method", "POST").
				Msg("Name Match")
		}
		return c.JSON(http.StatusOK, result)
	}
}

func setLogger(e *echo.Echo, m MatcherService) nsq.NSQ {
	cfg := m.GetConfig()
	nsqAddr := cfg.NsqdTCPAddress
	withLogs := cfg.WithWebLogs
	contains := cfg.NsqdContainsFilter
	regex := cfg.NsqdRegexFilter

	if nsqAddr != "" {
		cfg := nsqcfg.Config{
			StderrLogs: withLogs,
			Topic:      "gnmatcher",
			Address:    nsqAddr,
			Contains:   contains,
			Regex:      regex,
		}
		remote, err := nsqio.New(cfg)
		logCfg := middleware.DefaultLoggerConfig
		if err == nil {
			logCfg.Output = remote
			// set app logger too
			log.Logger = log.Output(remote)
		}
		e.Use(middleware.LoggerWithConfig(logCfg))
		if err != nil {
			log.Warn().Err(err)
		}
		return remote
	} else if withLogs {
		e.Use(middleware.Logger())
		return nil
	}
	return nil
}
