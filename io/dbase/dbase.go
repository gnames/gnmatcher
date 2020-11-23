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
func NewDB(cnf config.Config) *sql.DB {
	db, err := sql.Open("postgres", dbUrl(cnf))
	if err != nil {
		log.Fatalf("Cannot create PostgreSQL connection: %s.", err)
	}
	return db
}

func dbUrl(cnf config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cnf.PgUser, cnf.PgPass, cnf.PgHost, cnf.PgPort, cnf.PgDB)
}
