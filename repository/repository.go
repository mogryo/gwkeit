package repository

import (
	"context"
	"database/sql"

	"github.com/gwkeit/gwkeitdb"
)

type Repository struct {
	db            *sql.DB
	queries       *gwkeitdb.Queries
	activeVersion int64
}

func New(
	db *sql.DB,
	queries *gwkeitdb.Queries,
) *Repository {
	activeVersion := getLatestAppliedVersion(context.Background(), db)
	return &Repository{db: db, queries: queries, activeVersion: activeVersion}
}

func (r *Repository) checkVersion(ctx context.Context) {
	version := getLatestAppliedVersion(ctx, r.db)
	if version != r.activeVersion {
		panic("Database version mismatch")
	}
}

func getLatestAppliedVersion(ctx context.Context, db *sql.DB) int64 {
	rows, err := db.QueryContext(
		ctx,
		"SELECT max(version_id) FROM goose_db_version WHERE is_applied = 1;",
	)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	if err != nil {
		panic(err)
	}

	if rows == nil || !rows.Next() {
		panic("No applied version found")
	}

	var versionId int64
	err = rows.Scan(&versionId)

	if err != nil {
		panic(err)
	}
	return versionId
}
