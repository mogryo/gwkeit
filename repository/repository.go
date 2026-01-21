package repository

import (
	"database/sql"

	"github.com/gwkeit/gwkeitdb"
)

type Repository struct {
	db      *sql.DB
	queries *gwkeitdb.Queries
}

func New(
	db *sql.DB,
	queries *gwkeitdb.Queries,
) *Repository {
	return &Repository{db: db, queries: queries}
}
