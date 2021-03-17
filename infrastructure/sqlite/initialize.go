package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/chuckha/migrations"
	"github.com/pkg/errors"
)

// Initialize manages migrations for the test db
func (s *Store) Initialize(migrationDir string) error {
	// make sure migrations can persist
	_, err := s.DB.Exec(`CREATE TABLE IF NOT EXISTS migrations (number int)`)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the migration with the highest ID number
	row := s.DB.QueryRow(`SELECT number FROM migrations ORDER BY number DESC LIMIT 1`)
	if err := row.Err(); err != nil {
		return errors.WithStack(err)
	}
	latest := 0

	if err := row.Scan(&latest); err != nil {
		if err != sql.ErrNoRows {
			return errors.WithStack(err)
		}
	}
	// Load up the migrations and run starting from the last known run migration
	migs := migrations.FromDir(migrationDir)
	for _, migration := range migs[:latest] {
		fmt.Printf("Already ran migration: %d\n", migration.Order)
	}
	for _, migration := range migs[latest:] {
		fmt.Printf("Running migration %d\n", migration.Order)
		if _, err := s.DB.Exec(migration.Up); err != nil {
			return err
		}
		if _, err := s.DB.Exec(`INSERT INTO migrations (number) VALUES (?)`, migration.Order); err != nil {
			return err
		}
	}
	return nil
}
