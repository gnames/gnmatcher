// package dbase provides convenience methods for accessing PostgreSQL
// database.
package dbase

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/gnames/gnmatcher/pkg/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewDB creates a new instance of sql.DB using configuration data.
func NewDB(cfg config.Config) *sql.DB {
	db, err := sql.Open("postgres", dbUrl(cfg))
	if err != nil {
		slog.Error("Cannot create PostgreSQL connection", "error", err)
		os.Exit(1)
	}
	return db
}

func dbUrl(cfg config.Config) string {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDB)
	return dbURL
}
