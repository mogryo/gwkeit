package main

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"os"
	"path"

	"github.com/gwkeit/configuration"
	"github.com/gwkeit/pages"
	"github.com/gwkeit/repository"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"github.com/gwkeit/gwkeitdb"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	ctx := context.Background()
	defer ctx.Done()

	homeDir, _ := os.UserHomeDir()
	dataSourceName := path.Join(homeDir, configuration.AppDirectory, configuration.DbName)

	err := os.Mkdir(path.Dir(dataSourceName), 0755)

	db, err := sql.Open("sqlite", dataSourceName)
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
	defer func(queries *gwkeitdb.Queries) {
		err := queries.Close()
		if err != nil {
			panic(err)
		}
	}(queries)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	r := repository.New(db, queries)

	initialState := configuration.ReadConfiguration()

	err = pages.Run(ctx, r, initialState)
	if err != nil {
		panic(err)
	}
}
