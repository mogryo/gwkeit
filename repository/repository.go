package repository

import (
	"context"
	"database/sql"

	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/utils"
)

type Repository struct {
	db      *sql.DB
	queries *gwkeitdb.Queries
}

type SnippetInput struct {
	Title string
	Body  string
	Tags  []string
	Urls  []string
}

func New(
	db *sql.DB,
	queries *gwkeitdb.Queries,
) *Repository {
	return &Repository{db: db, queries: queries}
}

func (r *Repository) FindSnippetTags(ctx context.Context, snippetId int64) []gwkeitdb.Tag {
	tags, err := r.queries.FindTagsBySnippetId(ctx, snippetId)

	if err != nil {
		panic(err)
	}

	return tags
}

func (r *Repository) FindSnippetUrls(ctx context.Context, snippetId int64) []gwkeitdb.Url {
	urls, err := r.queries.FindUrlsBySnippetId(ctx, snippetId)

	if err != nil {
		panic(err)
	}

	return urls
}

func (r *Repository) FindSnippets(ctx context.Context, tags []string) []gwkeitdb.Snippet {
	snippets, err := r.queries.FindSnippetsByTags(ctx, tags)

	if err != nil {
		panic(err)
	}

	return snippets
}

func (r *Repository) SaveSnippet(
	ctx context.Context,
	snippetInput *SnippetInput,
) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)

	existingTags, _ := qtx.FindTagsByTag(ctx, snippetInput.Tags)
	existingTagNames := utils.Map(existingTags, func(tag gwkeitdb.Tag) string {
		return tag.Tag
	})
	existingTagIds := utils.Map(existingTags, func(tag gwkeitdb.Tag) int64 {
		return tag.ID
	})

	nonExistingTags := utils.Difference(snippetInput.Tags, existingTagNames)
	if len(nonExistingTags) > 0 {
		for _, tag := range nonExistingTags {
			id, _ := qtx.InsertTag(ctx, tag)
			existingTagIds = append(existingTagIds, id)
		}
	}

	urlIds := make([]int64, 0)
	if len(snippetInput.Urls) > 0 {
		for _, url := range snippetInput.Urls {
			id, _ := qtx.InsertUrl(ctx, url)
			urlIds = append(urlIds, id)
		}
	}

	snippetId, _ := qtx.InsertSnippet(ctx, gwkeitdb.InsertSnippetParams{Title: snippetInput.Title, Body: snippetInput.Body})
	for _, tagId := range existingTagIds {
		_ = qtx.InsertSnippetTag(ctx, gwkeitdb.InsertSnippetTagParams{SnippetID: snippetId, TagID: tagId})
	}
	for _, urlId := range urlIds {
		_ = qtx.InsertUrlSnippet(ctx, gwkeitdb.InsertUrlSnippetParams{SnippetID: snippetId, UrlID: urlId})
	}

	return tx.Commit()
}
