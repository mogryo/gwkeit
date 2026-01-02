package main

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"

	"github.com/gwkeit/appui"
	"github.com/gwkeit/repository"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"github.com/gwkeit/gwkeitdb"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "./gwkeit.db")
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	queries := gwkeitdb.New(db)
	defer queries.Close()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	r := repository.New(db, queries)
	appUI := appui.New(ctx, r)
	err = appUI.Run()
	if err != nil {
		panic(err)
	}
}
