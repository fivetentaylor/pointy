package testutils

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrationRunner = &OnceRunner{
	once: sync.Once{},
	wg:   sync.WaitGroup{},
}

func EnsureDatabaseUp() {
	migrationRunner.mu.Lock()
	migrationRunner.wg.Add(1) // Increment the WaitGroup counter

	go func() {
		migrationRunner.once.Do(migrateDb) // Ensure migrateDb() is only executed once
		migrationRunner.wg.Done()          // Decrement the counter when f() is done
	}()

	migrationRunner.wg.Wait()
	migrationRunner.mu.Unlock()
}

func migrateDb() {
	_, fullPath, _, _ := runtime.Caller(0)
	basePath := getPathUpToTarget(fullPath, "pkg")

	migrationsDir := fmt.Sprintf("file://%s/db/migrations", basePath)
	log.Info("Running migrations", "dir", migrationsDir)

	m, err := migrate.New(
		migrationsDir,
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("MIGRATION ERROR: %v", err))
	}

	// Clear all tables before running migrations
	err = truncateAllTables(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Sprintf("TRUNCATE TABLES ERROR: %v", err))
	}

	m.Log = &Log{}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(fmt.Sprintf("MIGRATION UP ERROR: %v", err))
	}
	log.Info("Migration complete")
}

// truncateAllTables truncates all tables in the database other than schema_migrations
func truncateAllTables(databaseURL string) error {
	fmt.Println("TRUNCATE ALL TABLES")
	// Connect to the database
	db, err := sql.Open("postgres", databaseURL) // assuming a PostgreSQL database
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Fetch all table names
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name != 'schema_migrations'
	`)
	if err != nil {
		return fmt.Errorf("failed to fetch table names: %w", err)
	}
	defer rows.Close()

	// Truncate each table except schema_migrations
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		_, err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
		}
		log.Info(fmt.Sprintf("Table %s truncated", tableName))
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error occurred during table truncation: %w", err)
	}

	return nil
}

func getPathUpToTarget(fullPath, targetDir string) string {
	index := strings.Index(fullPath, targetDir)
	if index == -1 {
		return ""
	}
	return fullPath[:index+len(targetDir)]
}

type Log struct{}

func (l *Log) Printf(format string, v ...interface{}) {
	log.Print(fmt.Sprintf(format, v...))
}

func (l *Log) Println(args ...interface{}) {
	log.Info(args)
}

func (l *Log) Verbose() bool {
	return true
}
