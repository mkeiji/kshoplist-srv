package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

var Db *sql.DB

type AppDb struct{}

func (this AppDb) Init() {
	var (
		dbUser       = os.Getenv("DB_USER")
		dbPass       = os.Getenv("DB_PASS")
		dbName       = os.Getenv("DB_NAME")
		dbHost       = os.Getenv("DB_HOST")
		dbPort       = os.Getenv("DB_PORT")
		dbAddr       = fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v?sslmode=disable",
			dbUser,
			dbPass,
			dbHost,
			dbPort,
			dbName,
		)
		migrationDir = flag.String("migration.files", "./database/migrations", "Directory where the migration files are located")
		pgDSN        = flag.String("postgres.dsn", dbAddr, "PostgreSQL DSN")
	)

	flag.Parse()

	var dbErr error
	Db, dbErr = sql.Open("postgres", *pgDSN)
	if dbErr != nil {
		log.Fatalf("could not connect to the PostgreSQL database... %v", dbErr)
	}

	if err := Db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	// Run migrations
	driver, err := postgres.WithInstance(Db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationDir),
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	// m.Down()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	} else {
		log.Println("Database migrated")
	}
}

