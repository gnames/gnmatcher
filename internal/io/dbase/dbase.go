// package dbase provides convenience methods for accessing PostgreSQL
// database.
package dbase

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gnames/gnmatcher/pkg/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// StemIDsSQL selects IDs of stems that have at least one canonical form
// without the maybe_misspelling flag. Other stems are hidden from matching,
// otherwise a misspelling matches itself exactly, and fuzzy matching to the
// correct name never runs. gndb computes the flag, so the database has to be
// created by a version of gndb that has this column.
const StemIDsSQL = `
SELECT DISTINCT ns.canonical_stem_id
FROM name_strings ns
JOIN canonicals c ON c.id = ns.canonical_id
WHERE ns.canonical_stem_id IS NOT NULL AND NOT c.maybe_misspelling`

// NewDB creates a new instance of sql.DB using configuration data.
func NewDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl(cfg))
	if err != nil {
		slog.Error("Cannot create PostgreSQL connection", "error", err)
		return nil, err
	}
	return db, nil
}

func dbUrl(cfg config.Config) string {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDB)
	return dbURL
}
