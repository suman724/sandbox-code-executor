package sqlite

import (
	"database/sql"
	_ "embed"
	"strings"
)

//go:embed schema.sql
var schema string

func Bootstrap(db *sql.DB) error {
	statements := strings.Split(schema, ";")
	for _, stmt := range statements {
		query := strings.TrimSpace(stmt)
		if query == "" {
			continue
		}
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
