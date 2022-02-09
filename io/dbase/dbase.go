// package dbase provides convenience methods for accessing PostgreSQL
// database.
package dbase

import (
	"database/sql"
	"fmt"

	"github.com/gnames/gnmatcher/config"
	"github.com/rs/zerolog/log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewDB creates a new instance of sql.DB using configuration data.
func NewDB(cfg config.Config) *sql.DB {
	db, err := sql.Open("postgres", dbUrl(cfg))
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create PostgreSQL connection")
	}
	return db
}

func dbUrl(cfg config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDB)
}
