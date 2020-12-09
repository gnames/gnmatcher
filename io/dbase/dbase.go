// package dbase provides convenience methods for accessing PostgreSQL
// database.
package dbase

import (
	"database/sql"
	"fmt"

	"github.com/gnames/gnmatcher/config"
	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewDB creates a new instance of sql.DB using configuration data.
func NewDB(cfg config.Config) *sql.DB {
	db, err := sql.Open("postgres", dbUrl(cfg))
	if err != nil {
		log.Fatalf("Cannot create PostgreSQL connection: %s.", err)
	}
	return db
}

func dbUrl(cfg config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDB)
}
