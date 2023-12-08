package main

import (
	"database/sql"
	"embed"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/pressly/goose/v3"
)

//
//go:embed migrations/*.sql
var embedMigrations embed.FS

const driver = "pgx"

func main() {
	cfg := config.New()
	db, err := sql.Open(driver, cfg.ConnectionString)
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
