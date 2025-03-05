package testutils

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// EnsureStorage ensures that the storage is up and running and migrated
// - postgres
// - dynamo
// - s3
func EnsureStorage() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		EnsureDatabaseUp()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		EnsureDynamoTable()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		EnsureS3Bucket()
		wg.Done()
	}()

	wg.Wait()
}

func TestGormDb(t *testing.T) *gorm.DB {
	// Gorm
	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	gormdb, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("error opening db: %w", err))
	}
	sqlDB, err := gormdb.DB()
	if err != nil {
		t.Fatal(fmt.Errorf("error getting db: %w", err))
	}
	err = sqlDB.Ping()
	if err != nil {
		t.Fatal(fmt.Errorf("error pinging db: %w", err))
	}

	return gormdb
}
