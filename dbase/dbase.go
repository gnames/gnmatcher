// Package dbase is an interface to PostgreSQL database that contains Global
// Names index data
package dbase

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Dbase struct {
	PgHost string
	PgPort int
	PgUser string
	PgPass string
	PgDB   string
}

func NewDbase() Dbase {
	dbase := Dbase{
		PgHost: "0.0.0.0",
		PgPort: 5432,
		PgUser: "postgres",
		PgPass: "",
		PgDB:   "gnames",
	}
	return dbase
}

func (d Dbase) NewDB() *sql.DB {
	db, err := sql.Open("postgres", d.opts())
	if err != nil {
		log.Fatalf("Cannot create PostgreSQL connection: %s.", err)
	}
	return db
}

func (d Dbase) opts() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		d.PgHost, d.PgUser, d.PgPass, d.PgDB)
}
