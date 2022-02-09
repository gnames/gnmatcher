package rest

import (

	// "github.com/gorilla/mux"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	nsqcfg "github.com/sfgrp/lognsq/config"
	"github.com/sfgrp/lognsq/ent/nsq"
	"github.com/sfgrp/lognsq/io/nsqio"
)

// Run creates and runs a RESTful API service of gnmatcher.
// this API is described by OpenAPI schema at
// https://app.swaggerhub.com/apis/dimus/gnmatcher/1.0.0
func Run(m MatcherService) {
	log.Printf("Starting the HTTP API server on port %d.", m.Port())
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	loggerNSQ := setLogger(e, m)
	if loggerNSQ != nil {
		defer loggerNSQ.Stop()
	}

	e.GET("/", root)
	e.GET("/api/v1/ping", ping(m))
	e.GET("/api/v1/version", ver(m))
	e.POST("/api/v1/matches", matchPOST(m))

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

func matchPOST(m MatcherService) func(echo.Context) error {
	return func(c echo.Context) error {
		var names []string
		if err := c.Bind(&names); err != nil {
			return err
		}
		result := m.MatchNames(names)
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

func setLogger(e *echo.Echo, m MatcherService) nsq.NSQ {
	nsqAddr := m.WebLogsNsqdTCP()
	withLogs := m.WithWebLogs()

	if nsqAddr != "" {
		cfg := nsqcfg.Config{
			StderrLogs: withLogs,
			Topic:      "gnmatcher",
			Address:    nsqAddr,
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
