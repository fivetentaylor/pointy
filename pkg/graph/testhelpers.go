package graph

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	pgm "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	fixtures *testfixtures.Loader
)

func setupTestDB() (*gorm.DB, error) {
	connStr, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		return nil, fmt.Errorf("DATABASE_URL not set in environment variables")
	}

	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func migrator(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := pgm.WithInstance(db, &pgm.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../db/migrations",
		"test", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func loadFixtures(db *sql.DB) error {
	var err error

	// Configure the loader to match your database setup and fixtures directory
	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../db/fixtures"),
	)
	if err != nil {
		return err
	}

	return fixtures.Load()
}

func prettyPrintTable(db *sql.DB, tableName string) error {
	// Fetching column names
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 0", tableName))
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Preparing query
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err = db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Preparing to scan rows based on the number of columns
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Printing column names
	fmt.Println(strings.Join(columns, "\t"))

	// Printing rows
	for rows.Next() {
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return err
		}
		for i, col := range values {
			if col != nil {
				fmt.Printf("%v", col)
			}
			// Print tab unless it's the last column
			if i < len(columns)-1 {
				fmt.Print("\t")
			}
		}
		fmt.Print("\n")
	}
	return rows.Err()
}
