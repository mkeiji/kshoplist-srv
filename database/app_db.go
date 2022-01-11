package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

var Db *sql.DB

type AppDb struct{}

func (this AppDb) Init() {
	var (
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASS")
		dbName = os.Getenv("DB_NAME")
		dbHost = os.Getenv("DB_HOST")
		dbAddr = fmt.Sprintf(
			"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
			dbUser,
			dbPass,
			dbHost,
			dbName,
		)
		migrationDir = flag.String("migration.files", "./database/migrations", "Directory where the migration files are located ?")
		mysqlDSN     = flag.String("mysql.dsn", dbAddr, "Mysql DSN")
	)
	// dbAddr: "root:secret@tcp(localhost)/testdb"

	flag.Parse()

	var dbErr error
	Db, dbErr = sql.Open("mysql", *mysqlDSN)
	if dbErr != nil {
		log.Fatalf("could not connect to the MySQL database... %v", dbErr)
	}

	if err := Db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	// Run migrations
	driver, err := mysql.WithInstance(Db, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationDir), // file://path/to/directory
		"mysql",
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
