package queryable

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
)

// MigrateAndSeedWithCustomPath Migration + Seed
func MigrateAndSeedWithCustomPath(db *sql.DB, root string) {

	fmt.Println("Start migration & seed process")
	version, cond := doMigrate(db, "file://"+root+"/migrations")
	fmt.Println("Migration (latest,status)", version, cond)
	version, cond = doSeed(sqlx.NewDb(db, "postgres"), "./"+root+"/seeders")
	fmt.Println("Seed (latest,status)", version, cond)
}

// MigrateAndSeed migration + seed
func MigrateAndSeed(db *sql.DB) {
	MigrateAndSeedWithCustomPath(db, "internal/database")
}

func doMigrate(db *sql.DB, migrationDir string) (int, bool) {
	var err error
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("error when creating postgres instance: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationDir,
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal("error when creating database instance: ", err)
	}

	isMigrated := true

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error when migrate up: ", err)
		err := m.Down()
		if err != nil {
			log.Fatal("error when migrate down: ", err)
		}
		isMigrated = false
	}

	return getTableStatus(sqlx.NewDb(db, "postgres"), "schema_migrations").Version, isMigrated
}

func doSeed(db *sqlx.DB, seedDir string) (int, bool) {

	//Seed Log Storage

	rows, err := db.Queryx(`select * from information_schema.tables where "table_name" = 'seed_log'`)

	if err != nil {
		panic(err)
	}

	if rows.Next() == false {
		schema := `CREATE TABLE seed_log(
			version integer,
			dirty boolean
		);`
		_, err := db.Exec(schema)
		if err != nil {
			panic(err)
		}
	}

	var buffer bytes.Buffer
	var isSeeded bool
	var isSeededBefore bool

	status := getTableStatus(db, "seed_log")
	seeds := getVersions(seedDir)

	if status.Version > 0 {
		isSeededBefore = true
	}

	for _, seed := range seeds {
		isSeeded = false
		if seed.Version < status.Version || seed.Version == status.Version && status.Dirty == false {
			continue
		}

		fmt.Printf("Processing %s", seed.Path)
		b, err := ioutil.ReadFile(seed.Path)
		if err != nil {
			panic(err)
		}
		buffer.Reset()
		buffer.Write(b)

		tx, err := db.Beginx()

		if err != nil {
			fmt.Println(err)
		}

		_, err = tx.Exec(buffer.String())
		if err != nil {
			status.Dirty = true
			_, _ = tx.NamedExec(`UPDATE seed_log SET "version"=:version, "dirty"=:dirty`, status)
			tx.Rollback()
			panic(err)
		}

		status.Dirty = false
		status.Version = status.Version + 1

		if !isSeededBefore {
			_, _ = tx.NamedExec(`INSERT INTO seed_log ("version","dirty") VALUES(:version,:dirty);`, status)
			isSeededBefore = true
		} else {
			_, _ = tx.NamedExec(`UPDATE seed_log SET "version"=:version, "dirty"=:dirty`, status)
		}
		tx.Commit()

		isSeeded = true
	}

	return status.Version, isSeeded
}

// Status of seed log
type Status struct {
	Version int  `db:"version" json:"version"`
	Dirty   bool `db:"dirty" json:"dirty"`
}

func getTableStatus(db *sqlx.DB, table string) Status {
	var status Status
	query := fmt.Sprintf(`SELECT "version", "dirty" FROM "%s";`, table)
	err := db.Get(&status, query)
	if err != nil {
		//do nothing
	}
	return status
}

func getVersions(path string) []SeedModel {
	// seed logic

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var seeds []SeedModel

	for _, file := range files {
		fileName := file.Name()
		fileVersion := strings.Split(fileName, "_")[0]
		fileVersionConverted, err := strconv.Atoi(fileVersion)
		if err != nil {
			log.Fatal("seeders did not use correct naming conventions")
			panic(err)
		}
		fileAbsPath, err := filepath.Abs(path + "/" + fileName)
		if err != nil {
			log.Fatal("seeders did not use correct naming conventions")
			panic(err)
		}
		seeds = append(seeds, SeedModel{
			Version: fileVersionConverted,
			Path:    fileAbsPath,
		})
	}
	sort.Sort(ByVersion(seeds))
	return seeds
}

// SeedModel model representation
type SeedModel struct {
	Version int
	Path    string
}

// ByVersion sort
type ByVersion []SeedModel

func (a ByVersion) Len() int           { return len(a) }
func (a ByVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVersion) Less(i, j int) bool { return a[i].Version < a[j].Version }
