package repositories

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (db *DB) Migrate() error {
	db.Log.Debug("running migration")
	files, err := iofs.New(migrationFiles, "migrations") // get migrations from
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db.Conn.DB, &pgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", files, "pgx", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {
		if err != migrate.ErrNoChange {
			db.Log.Error("migration failed", "error", err)
			return err
		}
		db.Log.Debug("migration did not change anything")
	}

	db.Log.Debug("migration finished")
	return nil
}
