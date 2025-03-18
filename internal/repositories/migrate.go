package repositories

import (

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/nongrata2/musiclib/migrations"
)

func (db *DB) Migrate() error {
	db.Log.Debug("running migration")
	files, err := iofs.New(migrations.MigrationFiles, ".")
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
