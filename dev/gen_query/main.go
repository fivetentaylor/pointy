package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var dataTypeMap = map[string]func(columnType gorm.ColumnType) (dataType string){
	"uuid": func(columnType gorm.ColumnType) (dataType string) {
		nullable, ok := columnType.Nullable()
		if !ok {
			panic("idk why this would happen")
		}

		if nullable {
			return "*string"
		}

		return "string"
	},
	"text": func(columnType gorm.ColumnType) (dataType string) {
		nullable, ok := columnType.Nullable()
		if !ok {
			panic("idk why this would happen")
		}

		if nullable {
			return "*string"
		}

		return "string"
	},
}

// Put your custom queries here:

// Addyour Query interfaces here:
var interfaces = map[string][]interface{}{}

func main() {
	gormdb, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(fmt.Errorf("error opening db: %w", err))
		return
	}

	cfg := &gen.Config{
		OutPath:      "./pkg/query",
		ModelPkgPath: "./pkg/models",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	}
	cfg.WithDataTypeMap(dataTypeMap)

	g := gen.NewGenerator(*cfg)

	g.UseDB(gormdb)

	tableList, err := gormdb.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get all tables fail: %w", err))
	}

	for _, tableName := range tableList {
		log.Println(tableName)
		if inter, ok := interfaces[tableName]; ok {
			g.ApplyInterface(inter[0], inter[1])
			continue
		}
		model := g.GenerateModel(tableName)

		g.ApplyBasic(model)
	}

	// Generate the code
	g.Execute()
}
