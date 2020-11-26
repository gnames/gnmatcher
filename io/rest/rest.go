package rest

import (

	// "github.com/gorilla/mux"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

// Run creates and runs a RESTful API service of gnmatcher.
// this API is described by OpenAPI schema at
// https://app.swaggerhub.com/apis/dimus/gnmatcher/1.0.0
func Run(m MatcherService) {
	log.Printf("Starting the HTTP API server on port %d.", m.Port())
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/", root)
	e.GET("/api/v1/ping", ping(m))
	e.GET("/api/v1/version", ver(m))
	e.POST("/api/v1/matches", match(m))

	addr := fmt.Sprintf(":%d", m.Port())
	e.Logger.Fatal(e.Start(addr))
}

func root(c echo.Context) error {
	return c.String(http.StatusOK,
		`The OpenAPI is described at
https://app.swaggerhub.com/apis/dimus/gnmatcher/1.0.0`)
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

func match(m MatcherService) func(echo.Context) error {
	return func(c echo.Context) error {
		var names []string
		if err := c.Bind(&names); err != nil {
			return err
		}
		result := m.MatchNames(names)
		return c.JSON(http.StatusOK, result)
	}
}
