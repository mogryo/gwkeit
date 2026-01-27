package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/slicelib"
)

func (r *Repository) FindSnippetTags(ctx context.Context, snippetId int64) []gwkeitdb.Tag {
	r.checkVersion(ctx)
	tags, err := r.queries.FindTagsBySnippetId(ctx, snippetId)

	if err != nil {
		panic(err)
	}

	return tags
}

func (r *Repository) FindSnippetUrls(ctx context.Context, snippetId int64) []gwkeitdb.Url {
	r.checkVersion(ctx)
	urls, err := r.queries.FindUrlsBySnippetId(ctx, snippetId)

	if err != nil {
		panic(err)
	}

	return urls
}

func (r *Repository) FindSnippetsByTags(ctx context.Context, tags []string) []gwkeitdb.Snippet {
	r.checkVersion(ctx)
	snippets, err := r.queries.FindSnippetsByTags(ctx, tags)

	if err != nil {
		panic(err)
	}

	return snippets
}

func (r *Repository) FindSnippetsByLikeTags(ctx context.Context, tags []string) []gwkeitdb.Snippet {
	r.checkVersion(ctx)
	allSnippets := make([]gwkeitdb.Snippet, 0)
	for _, tag := range tags {
		snippets, err := r.queries.FindSnippetsByLikeTags(ctx, "%"+tag+"%")

		if err != nil {
			panic(err)
		}
		allSnippets = append(allSnippets, snippets...)
	}

	return slicelib.UniqueGet(allSnippets, func(snippet gwkeitdb.Snippet) int64 { return snippet.ID })
}

func (r *Repository) FindSnippetsByFts(ctx context.Context, tags []string) []gwkeitdb.Snippet {
	r.checkVersion(ctx)
	prefixTags := slicelib.Map(tags, func(tag string) string { return tag + "*" })
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT rowid, title, body, description, url FROM snippets_fts WHERE snippets_fts MATCH $1;",
		strings.Join(prefixTags, " OR "),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	resultRows := make([]gwkeitdb.Snippet, 0)
	for rows.Next() {
		var (
			id          int64
			title       string
			body        string
			description string
			url         string
		)
		if err := rows.Scan(&id, &title, &body, &description, &url); err != nil {
			log.Fatal(err)
		}
		newSnippet := gwkeitdb.Snippet{
			ID:          id,
			Title:       title,
			Body:        body,
			Description: description,
			Url:         url,
		}
		resultRows = append(resultRows, newSnippet)
	}

	return resultRows
}

func (r *Repository) FindSnippetsByPage(ctx context.Context, page int64, size int64) ([]gwkeitdb.Snippet, error) {
	r.checkVersion(ctx)

	if page <= 0 {
		return nil, errors.New("page must be greater than 0")
	}

	return r.queries.FindSnippetsPaginated(
		ctx,
		gwkeitdb.FindSnippetsPaginatedParams{Limit: size, Offset: (page - 1) * size},
	)
}

func (r *Repository) FindSnippet(ctx context.Context, id int64) gwkeitdb.FindSnippetDataByIdRow {
	r.checkVersion(ctx)
	snippet, err := r.queries.FindSnippetDataById(ctx, id)
	if err != nil {
		panic(err)
	}

	return snippet
}

func (r *Repository) GetSnippetCount(ctx context.Context) int64 {
	r.checkVersion(ctx)
	count, err := r.queries.GetSnippetCount(ctx)
	if err != nil {
		panic(err)
	}

	return count
}
