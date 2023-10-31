package main

import (
	"database/sql"
	"embed"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"os"
)

//
//go:embed migrations/*.sql
var embedMigrations embed.FS

const driver = "postgres"

func main() {
	db, err := sql.Open(driver, os.Getenv("CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(driver); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
