package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/gwkeit/dto"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/slicelib"
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

func (r *Repository) FindSnippetsByTags(ctx context.Context, tags []string) []gwkeitdb.Snippet {
	snippets, err := r.queries.FindSnippetsByTags(ctx, tags)

	if err != nil {
		panic(err)
	}

	return snippets
}

func (r *Repository) FindSnippetsByLikeTags(ctx context.Context, tags []string) []gwkeitdb.Snippet {
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

func (r *Repository) SaveSnippet(
	ctx context.Context,
	snippetInput *dto.Snippet,
) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)

	snippetId, _ := qtx.InsertSnippet(
		ctx, gwkeitdb.InsertSnippetParams{
			Title:       snippetInput.Title,
			Body:        snippetInput.Body,
			Description: snippetInput.Description,
			Url:         snippetInput.UrlText,
		},
	)

	existingTags, _ := qtx.FindTagsByTag(ctx, snippetInput.Tags)
	existingTagNames := slicelib.Map(existingTags, func(tag gwkeitdb.Tag) string {
		return tag.Tag
	})
	existingTagIds := slicelib.Map(existingTags, func(tag gwkeitdb.Tag) int64 {
		return tag.ID
	})

	nonExistingTags := slicelib.Difference(snippetInput.Tags, existingTagNames)
	if len(nonExistingTags) > 0 {
		for _, tag := range nonExistingTags {
			id, _ := qtx.InsertTag(ctx, tag)
			existingTagIds = append(existingTagIds, id)
		}
	}

	urlIds := make([]int64, 0)
	if len(snippetInput.UrlList) > 0 {
		for _, url := range snippetInput.UrlList {
			id, _ := qtx.InsertUrl(ctx, gwkeitdb.InsertUrlParams{Url: url, SnippetID: snippetId})
			urlIds = append(urlIds, id)
		}
	}

	for _, tagId := range existingTagIds {
		_ = qtx.InsertSnippetTag(ctx, gwkeitdb.InsertSnippetTagParams{SnippetID: snippetId, TagID: tagId})
	}

	return snippetId, tx.Commit()
}

func (r *Repository) UpdateSnippet(
	ctx context.Context,
	snippetId int64,
	newSnippetData *dto.Snippet,
) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)

	err = qtx.UpdateSnippet(ctx, gwkeitdb.UpdateSnippetParams{ID: snippetId, Title: newSnippetData.Title, Body: newSnippetData.Body})
	if err != nil {
		return err
	}

	err = r.updateSnippetTags(ctx, qtx, snippetId, newSnippetData.Tags)
	if err != nil {
		return err
	}

	err = r.updateSnippetUrls(ctx, qtx, snippetId, newSnippetData.UrlList)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) FindSnippet(ctx context.Context, id int64) gwkeitdb.FindSnippetDataByIdRow {
	snippet, err := r.queries.FindSnippetDataById(ctx, id)
	if err != nil {
		panic(err)
	}

	return snippet
}

func (r *Repository) updateSnippetTags(
	ctx context.Context,
	qtx *gwkeitdb.Queries,
	snippetId int64,
	newTagNames []string,
) error {
	existingSnippetTags, err := qtx.FindTagsBySnippetId(ctx, snippetId)
	if err != nil {
		return err
	}

	unusedTags := slicelib.DifferenceGetA(existingSnippetTags, newTagNames, func(tag gwkeitdb.Tag) string { return tag.Tag })
	err = qtx.DeleteSnippetTags(
		ctx,
		gwkeitdb.DeleteSnippetTagsParams{
			SnippetID: snippetId,
			TagIds:    slicelib.Map(unusedTags, func(tag gwkeitdb.Tag) int64 { return tag.ID }),
		},
	)
	if err != nil {
		return err
	}

	newSnippetTagNames := slicelib.Difference(newTagNames, slicelib.Map(existingSnippetTags, func(tag gwkeitdb.Tag) string { return tag.Tag }))
	existingTagNames, _ := qtx.FindTagsByTag(ctx, newSnippetTagNames)
	existingTagIds := slicelib.Map(existingTagNames, func(tag gwkeitdb.Tag) int64 { return tag.ID })

	tagNamesToAdd := slicelib.DifferenceGetB(newSnippetTagNames, existingTagNames, func(tag gwkeitdb.Tag) string { return tag.Tag })
	for _, tagName := range tagNamesToAdd {
		id, err := qtx.InsertTag(ctx, tagName)
		if err != nil {
			return err
		}

		existingTagIds = append(existingTagIds, id)
	}

	for _, tagId := range existingTagIds {
		_ = qtx.InsertSnippetTag(ctx, gwkeitdb.InsertSnippetTagParams{SnippetID: snippetId, TagID: tagId})
	}

	err = r.deleteOrphanTags(ctx, qtx, unusedTags)
	return err
}

func (r *Repository) updateSnippetUrls(
	ctx context.Context,
	qtx *gwkeitdb.Queries,
	snippetId int64,
	newUrls []string,
) error {
	existingSnippetUrls, err := qtx.FindUrlsBySnippetId(ctx, snippetId)

	unusedUrls := slicelib.DifferenceGetA(existingSnippetUrls, newUrls, func(url gwkeitdb.Url) string { return url.Url })
	err = qtx.DeleteSnippetUrls(ctx, slicelib.Map(unusedUrls, func(url gwkeitdb.Url) int64 { return url.ID }))
	if err != nil {
		return err
	}

	newSnippetUrls := slicelib.Difference(newUrls, slicelib.Map(existingSnippetUrls, func(url gwkeitdb.Url) string { return url.Url }))
	for _, url := range newSnippetUrls {
		_, err = qtx.InsertUrl(ctx, gwkeitdb.InsertUrlParams{Url: url, SnippetID: snippetId})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) deleteOrphanTags(ctx context.Context, qtx *gwkeitdb.Queries, tags []gwkeitdb.Tag) error {
	for _, tag := range tags {
		exists, err := qtx.SnippetTagExists(ctx, tag.Tag)
		if err != nil {
			return err
		}
		if exists == 0 {
			err = qtx.DeleteTagByTag(ctx, tag.Tag)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
