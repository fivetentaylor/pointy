package testutils

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

var dynamoRunner = &OnceRunner{
	once: sync.Once{},
	wg:   sync.WaitGroup{},
}

func EnsureDynamoTable() {
	dynamoRunner.mu.Lock()
	dynamoRunner.wg.Add(1)
	go func() {
		dynamoRunner.once.Do(createDynamoTable)
		dynamoRunner.wg.Done()
	}()
	dynamoRunner.wg.Wait()
	dynamoRunner.mu.Unlock()
}

func createDynamoTable() {
	log.Info("DYNAMO: Creating dynamo table")
	dydb, err := dynamo.NewDB()
	if err != nil {
		log.Fatal(fmt.Sprintf("DYNAMO ERROR: %v", err))
	}

	err = dydb.CreateTestTable()
	if err != nil {
		log.Fatal(fmt.Sprintf("DYNAMO CREATE TABLE ERROR: %v", err))
	}
	log.Info("DYNAMO: Created dynamo table")
}
