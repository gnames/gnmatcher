package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mlib "github.com/gnames/gnlib/ent/matcher"
	"github.com/gnames/gnmatcher/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var apiPath = "/api/v1/"

// Run creates and runs a RESTful API service of gnmatcher.
// this API is described by OpenAPI schema at
// https://apidoc.gnames.org/gnmatcher
func Run(m MatcherService) {
	slog.Info("Starting HTTP API server", "port", m.Port())
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.GET("/", root)
	e.GET("/api/v1", root)
	e.GET("/api/v1/", root)
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
https://apidoc.globalnames.org/gnmatcher

API path: /api/v1/
		
Endpoints:
    ping/
    version/
    matches/
`)
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

		fuzzyRelaxed := c.QueryParam("fuzzy_relaxed") == "true"
		if fuzzyRelaxed {
			opts = append(opts, config.OptWithRelaxedFuzzyMatch(true))
		}

		fuzzyUni := c.QueryParam("fuzzy_uninomial") == "true"
		if fuzzyUni {
			opts = append(opts, config.OptWithUninomialFuzzyMatch(true))
		}
		if len(ds) > 0 {
			opts = append(opts, config.OptDataSources(ds))
		}

		result := m.MatchNames(names, opts...)
		if l := len(names); l > 0 {
			slog.Info("Names match",
				"namesNum", l,
				"example", names[0],
				"method", "GET")
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
		if inp.WithRelaxedFuzzyMatch {
			opts = append(opts, config.OptWithRelaxedFuzzyMatch(true))
		}
		if inp.WithUninomialFuzzyMatch {
			opts = append(opts, config.OptWithUninomialFuzzyMatch(true))
		}
		if len(inp.DataSources) > 0 {
			opts = append(opts, config.OptDataSources(inp.DataSources))
		}

		result := m.MatchNames(inp.Names, opts...)
		if l := len(inp.Names); l > 0 {
			slog.Info("Names match",
				"namesNum", l,
				"example", inp.Names[0],
				"method", "POST")
		}
		return c.JSON(http.StatusOK, result)
	}
}
