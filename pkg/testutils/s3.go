package testutils

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/teamreviso/code/pkg/storage/s3"
)

var s3Runner = &OnceRunner{
	once: sync.Once{},
	wg:   sync.WaitGroup{},
}

func EnsureS3Bucket() {
	s3Runner.mu.Lock()
	s3Runner.wg.Add(1)

	go func() {
		s3Runner.once.Do(createS3Bucket)
		s3Runner.wg.Done()
	}()
	s3Runner.wg.Wait()
	s3Runner.mu.Unlock()
}

func createS3Bucket() {
	log.Info("S3: Creating s3 bucket")
	s3, err := s3.NewS3()
	if err != nil {
		log.Fatal(fmt.Sprintf("S3 ERROR: %v", err))
	}

	err = s3.CreateTestBucket()
	if err != nil {
		log.Fatal(fmt.Sprintf("S3 CREATE BUCKET ERROR: %v", err))
	}
	log.Info("S3: Created s3 bucket")
}
