package migration

import (
	"database/sql"
	"log"

	"phonebook/config"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

var db *sql.DB

func currentVersion() int {
	var _version = -1
	_rows, _err := db.Query(`SELECT max("version") as "version" FROM "schema_migrations";`)
	if _err == nil {
		_rows.Next()
		_rows.Scan(&_version)
	}
	return _version
}

//RunMigrations ...
func RunMigrations() {
	var err error
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("error when open configuration %v", err)
	}
	db, err = sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		log.Fatal("error when open postgres connection: ", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("error when creating postgres instance: ", err)
	}

	var m *migrate.Migrate

	m, err = migrate.NewWithDatabaseInstance(
		"file://database/migrations/",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal("error when creating database instance: ", err)
	}

	_prevVersion := currentVersion()
	err = m.Up()
	_currVersion := currentVersion()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("error when migrate up: ", err)
		err := m.Down()
		if err != nil {
			log.Fatal("error when migrate down: ", err)
		}
	} else if _prevVersion < _currVersion && _currVersion == 1 {
		// something todo here
	}

}
