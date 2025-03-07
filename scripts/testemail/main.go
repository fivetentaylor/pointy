package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/service/email"
	"github.com/fivetentaylor/pointy/pkg/testutils"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	ctx := testutils.TestContext()

	userTbl := env.Query(ctx).User

	user, err := userTbl.Where(userTbl.Email.Eq("jpozdena@gmail.com")).First()
	if err != nil {
		log.Fatal(err)
	}

	docTbl := env.Query(ctx).Document

	doc, err := docTbl.Where(docTbl.ID.Eq("12c2e68c-0e8d-4fe9-98d9-eedf9982a7cc")).First()
	if err != nil {
		log.Fatal(err)
	}

	err = email.SendFirstOpen(ctx, "james@revi.so", user, doc)
	if err != nil {
		log.Fatal(err)
	}
}
